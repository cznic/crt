// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"unsafe"

	"github.com/cznic/ccir/libc/stdio"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
)

var (
	_ vaReader = (*argsReader)(nil)

	stdin, stdout, stderr uintptr
)

var (
	files = &fmap{
		m: map[uintptr]*os.File{},
	}
	nullReader = bytes.NewBuffer(nil)
)

type fmap struct {
	m  map[uintptr]*os.File
	mu sync.Mutex
}

func (m *fmap) add(f *os.File, u uintptr) {
	m.mu.Lock()
	m.m[u] = f
	m.mu.Unlock()
}

func (m *fmap) reader(u uintptr) io.Reader {
	switch u {
	case stdin:
		return os.Stdin
	case stdout, stderr:
		return nullReader
	}

	m.mu.Lock()
	f := m.m[u]
	m.mu.Unlock()
	return f
}

func (m *fmap) writer(u uintptr) io.Writer {
	switch u {
	case stdin:
		return ioutil.Discard
	case stdout:
		return os.Stdout
	case stderr:
		return os.Stderr
	}

	m.mu.Lock()
	f := m.m[u]
	m.mu.Unlock()
	return f
}

func (m *fmap) extract(u uintptr) *os.File {
	m.mu.Lock()
	f := m.m[u]
	delete(m.m, u)
	m.mu.Unlock()
	return f
}

type vaReader interface {
	readF64() float64
	readI32() int32
	readI64() int64
	readLong() int64
	readPtr() unsafe.Pointer
	readU32() uint32
	readU64() uint64
	readULong() uint64
}

type argsReader []interface{}

func (r *argsReader) readPtr() unsafe.Pointer {
	s := *r
	v := s[0].(unsafe.Pointer)
	*r = s[1:]
	return v
}

func (r *argsReader) readF64() float64 {
	s := *r
	v := s[0].(float64)
	*r = s[1:]
	return v
}

func (r *argsReader) readI32() (v int32) {
	s := *r
	switch x := s[0].(type) {
	case int32:
		v = x
	case uint32:
		v = int32(x)
	case int64:
		v = int32(x)
	case uint64:
		v = int32(x)
	default:
		panic(fmt.Errorf("%T", x))
	}
	*r = s[1:]
	return v
}

func (r *argsReader) readU32() uint32 {
	s := *r
	v := s[0].(uint32)
	*r = s[1:]
	return v
}

func (r *argsReader) readI64() int64 {
	s := *r
	v := s[0].(int64)
	*r = s[1:]
	return v
}

func (r *argsReader) readU64() uint64 {
	s := *r
	v := s[0].(uint64)
	*r = s[1:]
	return v
}

// void __register_stdfiles(void *, void *, void *);
func X__register_stdfiles(in, out, err uintptr) {
	stdin = in
	stdout = out
	stderr = err
}

// int printf(const char *format, ...);
func Xprintf(format *int8, args ...interface{}) int32 {
	r := argsReader(args)
	return goFprintf(os.Stdout, format, &r)
}

// int sprintf(char *str, const char *format, ...);
func Xsprintf(str, format *int8, args ...interface{}) int32 {
	w := memWriter(uintptr(unsafe.Pointer(str)))
	r := argsReader(args)
	n := goFprintf(&w, format, &r)
	w.WriteByte(0)
	return n
}

func goFprintf(w io.Writer, format *int8, ap vaReader) int32 {
	var b buffer.Bytes
	written := 0
	for {
		ch := *format
		*(*uintptr)(unsafe.Pointer(&format))++
		switch ch {
		case 0:
			_, err := b.WriteTo(w)
			b.Close()
			if err != nil {
				return -1
			}

			return int32(written)
		case '%':
			modifiers := ""
			long := 0
			var w []interface{}
		more:
			ch := *format
			*(*uintptr)(unsafe.Pointer(&format))++
			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				modifiers += string(ch)
				goto more
			case '*':
				w = append(w, ap.readI32())
				modifiers += string(ch)
				goto more
			case 'c':
				arg := ap.readI32()
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sc", modifiers), append(w, arg)...)
				written += n
			case 'd', 'i':
				var arg interface{}
				switch long {
				case 0:
					arg = ap.readI32()
				case 1:
					arg = ap.readLong()
				default:
					arg = ap.readI64()
				}
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sd", modifiers), append(w, arg)...)
				written += n
			case 'u':
				var arg interface{}
				switch long {
				case 0:
					arg = ap.readU32()
				case 1:
					arg = ap.readULong()
				default:
					arg = ap.readU64()
				}
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sd", modifiers), append(w, arg)...)
				written += n
			case 'x':
				var arg interface{}
				switch long {
				case 0:
					arg = ap.readU32()
				case 1:
					arg = ap.readULong()
				default:
					arg = ap.readU64()
				}
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sx", modifiers), append(w, arg)...)
				written += n
			case 'l':
				long++
				goto more
			case 'f':
				arg := ap.readF64()
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sf", modifiers), append(w, arg)...)
				written += n
			case 'p':
				arg := ap.readPtr()
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sp", modifiers), append(w, arg)...)
				written += n
			case 'g':
				arg := ap.readF64()
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sg", modifiers), append(w, arg)...)
				written += n
			case 's':
				arg := (*int8)(ap.readPtr())
				if arg == nil {
					break
				}

				var b2 buffer.Bytes
				for {
					c := *arg
					*(*uintptr)(unsafe.Pointer(&arg))++
					if c == 0 {
						n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%ss", modifiers), append(w, b2.Bytes())...)
						b2.Close()
						written += n
						break
					}

					b2.WriteByte(byte(c))
				}
			default:
				panic(fmt.Errorf("TODO %q", "%"+string(ch)))
			}
		default:
			b.WriteByte(byte(ch))
			written++
			if ch == '\n' {
				if _, err := b.WriteTo(w); err != nil {
					b.Close()
					return -1
				}
				b.Reset()
			}
		}
	}
}

// FILE *fopen64(const char *path, const char *mode);
func Xfopen64(path, mode *int8) file {
	p := GoString(path)
	var u uintptr
	switch p {
	case os.Stderr.Name():
		u = stderr
	case os.Stdin.Name():
		u = stdin
	case os.Stdout.Name():
		u = stdout
	default:
		var f *os.File
		var err error
		switch mode := GoString(mode); mode {
		case "r":
			if f, err = os.OpenFile(p, os.O_RDONLY, 0666); err != nil {
				switch {
				case os.IsNotExist(err):
					TODO("") // c.setErrno(errno.XENOENT)
				case os.IsPermission(err):
					TODO("") // c.setErrno(errno.XEPERM)
				default:
					TODO("") // c.setErrno(errno.XEACCES)
				}
			}
		case "w":
			if f, err = os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
				switch {
				case os.IsPermission(err):
					TODO("") // c.setErrno(errno.XEPERM)
				default:
					TODO("") // c.setErrno(errno.XEACCES)
				}
			}
		default:
			panic(mode)
		}
		if f != nil {
			u = malloc(ptrSize)
			files.add(f, u)
		}
	}
	return (file)(unsafe.Pointer(u))
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func fwrite(ptr uintptr, size, nmemb uint64, stream file) uint64 {
	hi, lo := mathutil.MulUint128_64(size, nmemb)
	if hi != 0 || lo > math.MaxInt32 {
		TODO("") // c.setErrno(errno.XE2BIG)
		return 0
	}

	n, err := files.writer(uintptr(unsafe.Pointer(stream))).Write((*[math.MaxInt32]byte)(unsafe.Pointer(ptr))[:lo])
	if err != nil {
		TODO("") // c.setErrno(errno.XEIO)
	}
	return uint64(n) / size
}

// int fclose(FILE *stream);
func Xfclose(stream file) int32 {
	u := uintptr(unsafe.Pointer(stream))
	switch u {
	case stdin, stdout, stderr:
		TODO("") // c.setErrno(errno.XEIO)
		return stdio.XEOF
	}

	f := files.extract(u)
	if f == nil {
		TODO("") // c.setErrno(errno.XEBADF)
		return stdio.XEOF
	}

	Xfree(u)
	if err := f.Close(); err != nil {
		TODO("") // c.setErrno(errno.XEIO)
		return stdio.XEOF
	}

	return 0
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func fread(ptr uintptr, size, nmemb uint64, stream file) uint64 {
	hi, lo := mathutil.MulUint128_64(size, nmemb)
	if hi != 0 || lo > math.MaxInt32 {
		TODO("") // c.setErrno(errno.XE2BIG)
		return 0
	}

	n, err := files.reader(uintptr(unsafe.Pointer(stream))).Read((*[math.MaxInt32]byte)(unsafe.Pointer(ptr))[:lo])
	if err != nil {
		TODO("") // c.setErrno(errno.XEIO)
	}
	return uint64(n) / size
}

// int fgetc(FILE *stream);
func Xfgetc(stream file) int32 {
	p := buffer.Get(1)
	if _, err := files.reader(uintptr(unsafe.Pointer(stream))).Read(*p); err != nil {
		buffer.Put(p)
		return stdio.XEOF
	}

	r := int32((*p)[0])
	buffer.Put(p)
	return r
}

// char *fgets(char *s, int size, FILE *stream);
func Xfgets(s *int8, size int32, stream file) *int8 {
	f := files.reader(uintptr(unsafe.Pointer(stream)))
	p := buffer.Get(1)
	b := *p
	w := memWriter(uintptr(unsafe.Pointer(s)))
	ok := false
	for i := int(size) - 1; i > 0; i-- {
		_, err := f.Read(b)
		if err != nil {
			if !ok {
				buffer.Put(p)
				return nil
			}

			break
		}

		ok = true
		w.WriteByte(b[0])
		if b[0] == '\n' {
			break
		}
	}
	w.WriteByte(0)
	buffer.Put(p)
	return s

}

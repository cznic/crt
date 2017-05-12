// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64
// +build !windows

package crt

import (
	"math"
	"os"
	"unsafe"

	"github.com/cznic/ccir/libc/stdio"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
)

const (
	longBits = 64
)

func (r *argsReader) readLong() int64 {
	s := *r
	v := s[0].(int64)
	*r = s[1:]
	return v
}

func (r *argsReader) readULong() uint64 {
	s := *r
	v := s[0].(uint64)
	*r = s[1:]
	return v
}

// FILE *fopen64(const char *path, const char *mode);
func Xfopen64(path, mode *int8) *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
} {
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
			TODO("") // u = c.m.malloc(int(unsafe.Sizeof(file{})))
			files.add(f, u)
		}
	}
	return (*struct {
		X0  int32
		X1  *int8
		X2  *int8
		X3  *int8
		X4  *int8
		X5  *int8
		X6  *int8
		X7  *int8
		X8  *int8
		X9  *int8
		X10 *int8
		X11 *int8
		X12 uintptr
		X13 uintptr
		X14 int32
		X15 int32
		X16 int64
		X17 uint16
		X18 int8
		X19 [1]int8
		X20 uintptr
		X21 int64
		X22 uintptr
		X23 uintptr
		X24 uintptr
		X25 uintptr
		X26 uint64
		X27 int32
		X28 [20]int8
	})(unsafe.Pointer(u))
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfwrite(ptr uintptr, size, nmemb uint64, stream *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}) uint64 {
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
func Xfclose(stream *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}) int32 {
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

	TODO("") // c.m.free(u)
	if err := f.Close(); err != nil {
		TODO("") // c.setErrno(errno.XEIO)
		return stdio.XEOF
	}

	return 0
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfread(ptr uintptr, size, nmemb uint64, stream *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}) uint64 {
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
func Xfgetc(stream *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}) int32 {
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
func Xfgets(s *int8, size int32, stream *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}) *int8 {
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

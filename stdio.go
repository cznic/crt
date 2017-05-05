// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"io"
	"os"
	"unsafe"

	"github.com/cznic/internal/buffer"
)

var (
	_ vaReader = (*argsReader)(nil)

	stdin, stdout, stderr uintptr
)

type vaReader interface {
	readF64() float64
	readI32() int32
	readI64() int64
	readLong() int64
	readPtr() uintptr
	readU32() uint32
	readU64() uint64
	readULong() uint64
}

type argsReader []interface{}

func (r *argsReader) readPtr() uintptr {
	s := *r
	v := s[0].(uintptr)
	*r = s[1:]
	return v
}

func (r *argsReader) readF64() float64 {
	s := *r
	v := s[0].(float64)
	*r = s[1:]
	return v
}

func (r *argsReader) readI32() int32 {
	s := *r
	v := s[0].(int32)
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
func Xprintf(format uintptr, args ...interface{}) int32 {
	r := argsReader(args)
	return goFprintf(os.Stdout, format, &r)
}

func goFprintf(w io.Writer, format uintptr, ap vaReader) int32 {
	var b buffer.Bytes
	written := 0
	for {
		ch := readI8(format)
		format++
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
			ch := readI8(format)
			format++
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
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sp", modifiers), append(w, unsafe.Pointer(arg))...)
				written += n
			case 'g':
				arg := ap.readF64()
				n, _ := fmt.Fprintf(&b, fmt.Sprintf("%%%sg", modifiers), append(w, arg)...)
				written += n
			case 's':
				arg := ap.readPtr()
				if arg == 0 {
					break
				}

				var b2 buffer.Bytes
				for {
					c := readI8(arg)
					arg++
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

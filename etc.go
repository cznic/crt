// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crt provides C-runtime services. (Work In Progress)
package crt

import (
	"fmt"
	"math"
	"os"
	"path"
	"runtime"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
)

const (
	ptrSize   = mathutil.UintPtrBits / 8
	heapAlign = 2 * ptrSize
)

var (
	brk           unsafe.Pointer
	heapAvailable int64
	threadID      uintptr
)

func writeU8(p uintptr, v uint8) { *(*uint8)(unsafe.Pointer(p)) = v }

type TLS struct {
	threadID uintptr
	errno    int32
}

func NewTLS() *TLS { return &TLS{threadID: atomic.AddUintptr(&threadID, 1)} }

func (t *TLS) setErrno(err interface{}) {
	switch x := err.(type) {
	case int:
		t.errno = int32(x)
	case *os.PathError:
		t.setErrno(x.Err)
	case syscall.Errno:
		t.errno = int32(x)
	default:
		panic(fmt.Errorf("TODO %T(%#v)", x, x))
	}
}

func (t *TLS) Close() { free(unsafe.Pointer(t)) }

//TODO remove me.
func TODO(msg string, more ...interface{}) string { //TODOOK
	_, fn, fl, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%s:%d: %v\n", path.Base(fn), fl, fmt.Sprintf(msg, more...))
	os.Stderr.Sync()
	panic(fmt.Errorf("%s:%d: TODO %v", path.Base(fn), fl, fmt.Sprintf(msg, more...))) //TODOOK
}

type memWriter uintptr

func (m *memWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	*m += memWriter(movemem(unsafe.Pointer(*m), unsafe.Pointer(&b[0]), len(b)))
	return len(b), nil
}

func (m *memWriter) WriteByte(b byte) error {
	*(*byte)(unsafe.Pointer(*m)) = b
	*m++
	return nil
}

func movemem(dst, src unsafe.Pointer, n int) int {
	return copy((*[math.MaxInt32]byte)(dst)[:n], (*[math.MaxInt32]byte)(src)[:n])
}

// GoString returns a string from a C char* null terminated string s.
func GoString(s *int8) string {
	if s == nil {
		return ""
	}

	var b buffer.Bytes
	for {
		ch := *s
		if ch == 0 {
			r := string(b.Bytes())
			b.Close()
			return r
		}

		b.WriteByte(byte(ch))
		*(*uintptr)(unsafe.Pointer(&s))++
	}
}

// GoStringLen returns a string from a C char* string s having length len bytes.
func GoStringLen(s *int8, len int) string {
	var b buffer.Bytes
	for ; len != 0; len-- {
		b.WriteByte(byte(*s))
		*(*uintptr)(unsafe.Pointer(&s))++
	}
	r := string(b.Bytes())
	b.Close()
	return r
}

func RegisterHeap(h unsafe.Pointer, n int64) {
	brk = h
	heapAvailable = n
}

// if n%m != 0 { n += m-n%m }. m must be a power of 2.
func roundupI64(n, m int64) int64 { return (n + m - 1) &^ (m - 1) }

// CString allocates a C string initialized from s.
func CString(s string) unsafe.Pointer {
	n := len(s)
	p := malloc(n + 1)
	copy((*[math.MaxInt32]byte)(p)[:n], s)
	(*[math.MaxInt32]byte)(p)[n] = 0
	return p
}

func malloc(size int) unsafe.Pointer {
	if size != 0 {
		var tls TLS
		p := sbrk(&tls, int64(size))
		if int64(uintptr(p)) > 0 {
			return p
		}
	}

	return nil
}

func calloc(size int) unsafe.Pointer {
	p := malloc(size)
	if p == nil {
		return nil
	}

	q := (*byte)(p)
	for ; size != 0; size-- {
		*q = 0
		*(*uintptr)(unsafe.Pointer(&q))++
	}
	return p
}

func realloc(tls *TLS, ptr unsafe.Pointer, size int) unsafe.Pointer {
	q := malloc(size)
	if q == nil {
		return nil
	}

	movemem(q, ptr, size)
	return q
}

func free(ptr unsafe.Pointer) { /*TODO*/ }

// CopyString copies src to dest, optionally adding a zero byte at the end.
func CopyString(dst unsafe.Pointer, src string, addNull bool) {
	copy((*[math.MaxInt32]byte)(dst)[:len(src)], src)
	if addNull {
		writeU8(uintptr(dst)+uintptr(len(src)), 0)
	}
}

// CopyBytes copies src to dest, optionally adding a zero byte at the end.
func CopyBytes(dst unsafe.Pointer, src []byte, addNull bool) {
	copy((*[math.MaxInt32]byte)(dst)[:len(src)], src)
	if addNull {
		writeU8(uintptr(dst)+uintptr(len(src)), 0)
	}
}

// GoBytesLen returns a []byte copied from a C char* string s having length len bytes.
func GoBytesLen(s *int8, len int) []byte {
	var b buffer.Bytes
	for ; len != 0; len-- {
		b.WriteByte(byte(*s))
		*(*uintptr)(unsafe.Pointer(&s))++
	}
	return b.Bytes()
}

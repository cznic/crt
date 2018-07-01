// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crt provides C-runtime services. (Work In Progress)
//
// Installation
//
//     $ go get github.com/cznic/crt
//
// Documentation: http://godoc.org/github.com/cznic/crt
//
// This package contains documentation from
// 	http://pubs.opengroup.org/onlinepubs/9699919799/
//
// 	The Open Group Base Specifications Issue 7, 2018 edition
// 	IEEE Std 1003.1™-2017 (Revision of IEEE Std 1003.1-2008)
// 	Copyright © 2001-2018 IEEE and The Open Group
package crt

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/cznic/crt/errno"
	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
	"github.com/cznic/memory"
)

const (
	ptrSize = mathutil.UintPtrBits / 8
)

var (
	_ io.Writer = (*memWriter)(nil)

	// TraceWriter is used for trace output.
	TraceWriter = io.Writer(os.Stderr)

	allocMu   sync.Mutex
	allocator memory.Allocator
	threadID  uintptr
)

// TLS represents the C-thread local storage.
type TLS uintptr

type tls struct {
	threadID uintptr
	errno    int32
}

func (t TLS) setErrno(err interface{}) {
	switch x := err.(type) {
	case int:
		(*tls)(unsafe.Pointer(t)).errno = int32(x)
	case *os.PathError:
		t.setErrno(x.Err)
	case *os.SyscallError:
		switch y := x.Err.(type) {
		case syscall.Errno:
			(*tls)(unsafe.Pointer(t)).errno = int32(y)
		default:
			panic(fmt.Errorf("crt.setErrno %T(%#v)", y, y))
		}
	case syscall.Errno:
		(*tls)(unsafe.Pointer(t)).errno = int32(x)
	default:
		panic(fmt.Errorf("crt.setErrno %T(%#v)", x, x))
	}
}

func (t TLS) err() int32 { return (*tls)(unsafe.Pointer(t)).errno }

func (t TLS) getThreadID() uintptr { return (*tls)(unsafe.Pointer(t)).threadID }

// NewTLS returns a newly created TLS, allocated outside of the Go runtime
// heap.  To free the TLS use Free(uintptr(unsafe.Pointer(tls))).
func NewTLS() TLS {
	t := MustCalloc(int(unsafe.Sizeof(tls{})))
	(*tls)(unsafe.Pointer(t)).threadID = atomic.AddUintptr(&threadID, 1)
	return TLS(t)
}

// void __register_stdfiles(void *stdin, void *stdout, void *stderr, void *environ);
func X__register_stdfiles(tls TLS, in, out, err, env uintptr) {
	stdin = in
	stdout = out
	stderr = err
	penviron = env
}

// void exit(int);
func X__builtin_exit(tls TLS, n int32) { os.Exit(int(n)) }

// BSS allocates the the bss segment of a package/command.
func BSS(init *byte) uintptr {
	r := uintptr(unsafe.Pointer(init))
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// DS allocates the the data segment of a package/command.
func DS(init []byte) uintptr {
	r := (*reflect.SliceHeader)(unsafe.Pointer(&init)).Data
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// TS allocates the R/O text segment of a package/command.
func TS(init string) uintptr { return (*reflect.StringHeader)(unsafe.Pointer(&init)).Data }

// Free frees memory allocated by Calloc, Malloc or Realloc.
func Free(p uintptr) error {
	allocMu.Lock()
	err := allocator.UintptrFree(p)
	allocMu.Unlock()
	return err
}

// Calloc allocates zeroed memory.
func Calloc(size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrCalloc(size)
	allocMu.Unlock()
	return p, err
}

// Malloc allocates memory.
func Malloc(size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrMalloc(size)
	allocMu.Unlock()
	return p, err
}

// MustCalloc is like Calloc but panics if the allocation cannot be made.
func MustCalloc(size int) uintptr {
	p, err := Calloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// MustMalloc is like Malloc but panics if the allocation cannot be made.
func MustMalloc(size int) uintptr {
	p, err := Malloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// Realloc reallocates memory.
func Realloc(p uintptr, size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrRealloc(p, size)
	allocMu.Unlock()
	return p, err
}

type memWriter uintptr

func (m *memWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	Copy(uintptr(*m), uintptr(unsafe.Pointer(&b[0])), len(b))
	*m += memWriter(len(b))
	return len(b), nil
}

func (m *memWriter) WriteByte(b byte) error {
	*(*byte)(unsafe.Pointer(*m)) = b
	*m++
	return nil
}

// Copy copies n bytes form src to dest and returns n.
func Copy(dst, src uintptr, n int) int {
	return copy((*rawmem)(unsafe.Pointer(dst))[:n], (*rawmem)(unsafe.Pointer(src))[:n])
}

// CString allocates a C string initialized from s.
func CString(s string) uintptr {
	n := len(s)
	var t tls
	p := malloc(TLS(unsafe.Pointer(&t)), n+1)
	if p == 0 {
		return 0
	}

	copy((*rawmem)(unsafe.Pointer(p))[:n], s)
	(*rawmem)(unsafe.Pointer(p))[n] = 0
	return p
}

// GoString returns a string from a C char* null terminated string s.
func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}

	var b buffer.Bytes
	for {
		ch := *(*byte)(unsafe.Pointer(s))
		if ch == 0 {
			r := string(b.Bytes())
			b.Close()
			return r
		}

		b.WriteByte(ch)
		s++
	}
}

func calloc(tls TLS, size int) uintptr {
	p, err := Calloc(size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}

func free(tls TLS, p uintptr) { Free(p) }

func malloc(tls TLS, size int) uintptr {
	p, err := Malloc(size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}

func realloc(tls TLS, p uintptr, size int) uintptr {
	p, err := Realloc(p, size)
	if err != nil {
		tls.setErrno(errno.XENOMEM)
		return 0
	}

	return p
}

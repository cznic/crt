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
	"unsafe"

	"github.com/cznic/internal/buffer"
	"github.com/cznic/mathutil"
)

const (
	ptrSize   = mathutil.UintPtrBits / 8
	heapAlign = 2 * ptrSize
)

var (
	heap          unsafe.Pointer
	heapAvailable int64
)

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

func RegisterHeap(h unsafe.Pointer, n int64) {
	heap = h
	heapAvailable = n
}

// if n%m != 0 { n += m-n%m }. m must be a power of 2.
func roundupI64(n, m int64) int64 { return (n + m - 1) &^ (m - 1) }

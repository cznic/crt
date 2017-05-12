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

	*m += memWriter(movemem(uintptr(*m), uintptr((unsafe.Pointer)(&b[0])), len(b)))
	return len(b), nil
}

func (m *memWriter) WriteByte(b byte) error {
	*(*byte)(unsafe.Pointer(m)) = b
	*m++
	return nil
}

func movemem(dst, src uintptr, n int) int {
	return copy((*[math.MaxInt32]byte)(unsafe.Pointer(dst))[:n], (*[math.MaxInt32]byte)(unsafe.Pointer(src))[:n])
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

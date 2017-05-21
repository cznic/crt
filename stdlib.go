// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os"
	"unsafe"

	"github.com/cznic/mathutil"
)

// void exit(int);
func Xexit(n int32) { X__builtin_exit(n) }

// void exit(int);
func X__builtin_exit(n int32) {
	os.Exit(int(n))
}

// void *malloc(size_t size);
func malloc(size int) unsafe.Pointer {
	if size != 0 && size <= mathutil.MaxInt {
		if size != 0 {
			p := sbrk(int64(size))
			if int64(uintptr(p)) > 0 {
				return p
			}
		}

	}

	return nil
}

// void free(void *ptr);
func Xfree(ptr unsafe.Pointer) {
	//TODO
}

// void abort();
func Xabort() { X__builtin_abort() }

// void abort();
func X__builtin_abort() { os.Exit(1) }

// void *realloc(void *ptr, size_t size);
func realloc(ptr unsafe.Pointer, size int) unsafe.Pointer {
	q := malloc(size)
	if q == nil {
		return nil
	}

	movemem(q, ptr, size)
	return q
}

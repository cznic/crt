// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"unsafe"
)

// void *malloc(size_t size);
func Xmalloc(size uint64) unsafe.Pointer { return malloc(int(size)) }

// void *realloc(void *ptr, size_t size);
func Xrealloc(ptr unsafe.Pointer, size uint64) unsafe.Pointer {
	return realloc(ptr, int(size))
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// void *malloc(size_t size);
func Xmalloc(size uint64) uintptr { return malloc(int(size)) }

// void *realloc(void *ptr, size_t size);
func Xrealloc(ptr uintptr, size uint64) uintptr {
	return realloc(ptr, int(size))
}

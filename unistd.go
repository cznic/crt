// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// void *sbrk(intptr_t increment);
func sbrk(increment int64) uintptr {
	if increment > heapAvailable {
		TODO("") // On error, (void *) -1 is returned, and errno is set to ENOMEM.
		return ^uintptr(0)
	}

	increment = roundupI64(increment, heapAlign)
	heapAvailable -= increment
	return heap + uintptr(increment)
}
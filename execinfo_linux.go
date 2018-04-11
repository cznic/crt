// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int backtrace(void **buffer, int size);
func Xbacktrace(tls TLS, buffer uintptr, size int32) int32 {
	return 0 //TODO
}

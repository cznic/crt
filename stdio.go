// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

var stdin, stdout, stderr uintptr

// void __register_stdfiles(void *, void *, void *);
func X__register_stdfiles(in, out, err uintptr) {
	stdin = in
	stdout = out
	stderr = err
}

// int printf(const char *format, ...);
func Xprintf(fomat uintptr, ap ...interface{}) int32 {
	return 0 //TODO
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"syscall"
	"unsafe"
)

// int open64(const char *pathname, int flags, ...);
func Xopen64(pathname *int8, flags int32, args ...interface{}) int32 {
	var mode int32
	if len(args) != 0 {
		mode = args[0].(int32)
	}
	r, _, err := syscall.Syscall(syscall.SYS_OPEN, uintptr(unsafe.Pointer(pathname)), uintptr(flags), uintptr(mode))
	//TODO if strace {
	//TODO 	fmt.Fprintf(os.Stderr, "open(%q, %v, %#o) %v %v\n", GoString(pathname), modeString(flags), mode, r, err)
	//TODO }
	if err != 0 {
		TODO("") //c.thread.setErrno(err)
	}
	return int32(r)
}

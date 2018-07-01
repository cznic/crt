// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// off_t lseek(int fd, off_t offset, int whence);
func Xlseek(tls TLS, fd int32, offset int32, whence int32) int32 {
	return int32(Xlseek64(tls, fd, int64(offset), whence))
}

// int ftruncate(int fildes, off_t length);
func Xftruncate(tls TLS, fildes int32, length int32) int32 {
	panic("TODO")
}

// int ftruncate(int fildes, off_t length);
func Xftruncate64(tls TLS, fildes int32, length int64) int32 {
	panic("TODO")
	r, _, err := syscall.Syscall(syscall.SYS_FTRUNCATE64, uintptr(fildes), uintptr(length), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "ftruncate(%#x, %#x) %v, %v\n", fildes, length, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

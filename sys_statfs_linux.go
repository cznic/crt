// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int fstatfs(int fd, struct statfs *buf);
func Xfstatfs(tls TLS, fd int32, buf uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FSTATFS, uintptr(fd), buf, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fstatfs(%v, %#x) %v %v\n", fd, buf, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

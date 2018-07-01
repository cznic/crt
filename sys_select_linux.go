// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int select(int nfds, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout);
func Xselect(tls TLS, nfds int32, readfds, writefds, exceptfds, timeout uintptr) int32 {
	r, _, err := syscall.Syscall6(syscall.SYS_SELECT, uintptr(nfds), readfds, writefds, exceptfds, timeout, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "select(%v, %#x, %#x, %#x, %#x) %v %v\n", nfds, readfds, writefds, exceptfds, timeout, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

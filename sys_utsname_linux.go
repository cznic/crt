// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// int uname(struct utsname *buf);
func Xuname(tls TLS, buf uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_UNAME, buf, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "uname(%#x) %v %v\n", buf, r, err)
		fmt.Fprintf(os.Stderr, "\tXsysname %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xsysname)))
		fmt.Fprintf(os.Stderr, "\tXnodename %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xnodename)))
		fmt.Fprintf(os.Stderr, "\tXrelease %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xrelease)))
		fmt.Fprintf(os.Stderr, "\tXversion %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xversion)))
		fmt.Fprintf(os.Stderr, "\tXmachine %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xmachine)))
		fmt.Fprintf(os.Stderr, "\tX__domainname %q\n", GoString(buf+unsafe.Offsetof(utsname{}.X__domainname)))
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

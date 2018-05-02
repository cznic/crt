// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"syscall"
	"unsafe"
)

// int uname(struct utsname *buf);
func Xuname(tls TLS, buf uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_UNAME, buf, 0, 0)
	if strace {
		fmt.Fprintf(TraceWriter, "uname(%#x) %v %v\n", buf, r, err)
		fmt.Fprintf(TraceWriter, "\tXsysname %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xsysname)))
		fmt.Fprintf(TraceWriter, "\tXnodename %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xnodename)))
		fmt.Fprintf(TraceWriter, "\tXrelease %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xrelease)))
		fmt.Fprintf(TraceWriter, "\tXversion %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xversion)))
		fmt.Fprintf(TraceWriter, "\tXmachine %q\n", GoString(buf+unsafe.Offsetof(utsname{}.Xmachine)))
		fmt.Fprintf(TraceWriter, "\tX__domainname %q\n", GoString(buf+unsafe.Offsetof(utsname{}.X__domainname)))
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

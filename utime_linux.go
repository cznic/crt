// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// int utime(const char *path, const struct utimbuf *times);
//
// The utime() function shall set the access and modification times of the file
// named by the path argument.
//
// If times is a null pointer, the access and modification times of the file
// shall be set to the current time. The effective user ID of the process shall
// match the owner of the file, or the process has write permission to the file
// or has appropriate privileges, to use utime() in this manner.
//
// If times is not a null pointer, times shall be interpreted as a pointer to a
// utimbuf structure and the access and modification times shall be set to the
// values contained in the designated structure. Only a process with the
// effective user ID equal to the user ID of the file or a process with
// appropriate privileges may use utime() this way.
//
// The utimbuf structure is defined in the <utime.h> header. The times in the
// structure utimbuf are measured in seconds since the Epoch.
//
// Upon successful completion, the utime() function shall mark the last file
// status change timestamp for update; see <sys/stat.h>.
//
// Upon successful completion, 0 shall be returned. Otherwise, -1 shall be
// returned and errno shall be set to indicate the error, and the file times
// shall not be affected.
func Xutime(tls TLS, path, times uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_UTIME, path, times, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "utime(%q, %#x) %v %v\n", GoString(path), times, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

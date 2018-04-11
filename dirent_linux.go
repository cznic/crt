// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"unsafe"
)

// int closedir(DIR *dirp);
//
// The closedir() function shall close the directory stream referred to by the
// argument dirp. Upon return, the value of dirp may no longer point to an
// accessible object of the type DIR. If a file descriptor is used to implement
// type DIR, that file descriptor shall be closed.
//
// Upon successful completion, closedir() shall return 0; otherwise, -1 shall
// be returned and errno set to indicate the error.
func Xclosedir(tls TLS, dirp uintptr) int32 {
	buf := *(*uintptr)(unsafe.Pointer(dirp + unsafe.Offsetof(S__dirstream{}.buf)))
	if err := Free(buf); err != nil {
		tls.setErrno(err)
		Free(dirp)
		return -1
	}

	if err := Free(dirp); err != nil {
		tls.setErrno(err)
		return -1
	}

	return 0
}

// int readdir64_r(DIR *dirp, struct dirent *entry, struct dirent64 **result);
func Xreaddir64_r(tls TLS, dirp, entry, result uintptr) int32 {
	panic("TODO")
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"io/ioutil"
	"syscall"
	"unsafe"
)

// DIR *opendir(const char *name);
//
// The opendir() function shall open a directory stream corresponding to the
// directory named by the dirname argument. The directory stream is positioned
// at the first entry. If the type DIR is implemented using a file descriptor,
// applications shall only be able to open up to a total of {OPEN_MAX} files
// and directories.
//
// If the type DIR is implemented using a file descriptor, the descriptor shall
// be obtained as if the O_DIRECTORY flag was passed to open().
//
// Upon successful completion, these functions shall return a pointer to an
// object of type DIR. Otherwise, these functions shall return a null pointer
// and set errno to indicate the error.
func Xopendir(tls TLS, name uintptr) (r uintptr) {
	s, err := ioutil.ReadDir(GoString(name))
	if err != nil {
		tls.setErrno(err)
		return 0
	}

	sz := unsafe.Sizeof(syscall.Dirent{})
	buf, err := Malloc((len(s) + 1) * int(sz))
	if err != nil {
		tls.setErrno(err)
		return 0
	}

	if r, err = Malloc(int(unsafe.Sizeof(S__dirstream{}))); err != nil {
		Free(buf)
		tls.setErrno(err)
		return 0
	}

	var stream S__dirstream
	stream.buf = buf
	stream.next = buf
	stream.n = len(s)
	*(*S__dirstream)(unsafe.Pointer(r)) = stream

	p := buf
	for _, v := range s {
		var d syscall.Dirent
		nm := v.Name()
		for i := 0; i < len(nm); i++ {
			if i == len(d.Name)-1 {
				break
			}

			d.Name[i] = nm[i]
			d.Ino = v.Sys().(*syscall.Stat_t).Ino
		}
		*(*syscall.Dirent)(unsafe.Pointer(p)) = d
		p += sz
	}
	return r
}

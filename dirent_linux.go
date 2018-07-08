// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"io/ioutil"
	"syscall"
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

			d.Name[i] = int8(nm[i])
			d.Ino = v.Sys().(*syscall.Stat_t).Ino
		}
		*(*syscall.Dirent)(unsafe.Pointer(p)) = d
		p += sz
	}
	return r
}

// struct dirent *readdir(DIR *dirp);
//
// The type DIR, which is defined in the <dirent.h> header, represents a
// directory stream, which is an ordered sequence of all the directory entries
// in a particular directory. Directory entries represent files; files may be
// removed from a directory or added to a directory asynchronously to the
// operation of readdir().
//
// The readdir() function shall return a pointer to a structure representing
// the directory entry at the current position in the directory stream
// specified by the argument dirp, and position the directory stream at the
// next entry. It shall return a null pointer upon reaching the end of the
// directory stream. The structure dirent defined in the <dirent.h> header
// describes a directory entry. The value of the structure's d_ino member shall
// be set to the file serial number of the file named by the d_name member. If
// the d_name member names a symbolic link, the value of the d_ino member shall
// be set to the file serial number of the symbolic link itself.
//
// The readdir() function shall not return directory entries containing empty
// names. If entries for dot or dot-dot exist, one entry shall be returned for
// dot and one entry shall be returned for dot-dot; otherwise, they shall not
// be returned.
//
// The application shall not modify the structure to which the return value of
// readdir() points, nor any storage areas pointed to by pointers within the
// structure. The returned pointer, and pointers within the structure, might be
// invalidated or the structure or the storage areas might be overwritten by a
// subsequent call to readdir() on the same directory stream. They shall not be
// affected by a call to readdir() on a different directory stream. The
// returned pointer, and pointers within the structure, might also be
// invalidated if the calling thread is terminated.
//
// If a file is removed from or added to the directory after the most recent
// call to opendir() or rewinddir(), whether a subsequent call to readdir()
// returns an entry for that file is unspecified.
//
// The readdir() function may buffer several directory entries per actual read
// operation; readdir() shall mark for update the last data access timestamp of
// the directory each time the directory is actually read.
//
// After a call to fork(), either the parent or child (but not both) may
// continue processing the directory stream using readdir(), rewinddir(), [XSI]
// ￼  or seekdir(). ￼ If both the parent and child processes use these
// functions, the result is undefined.
//
// The readdir() function need not be thread-safe.
//
// Applications wishing to check for error situations should set errno to 0
// before calling readdir(). If errno is set to non-zero on return, an error
// occurred.
//
// The storage pointed to by entry shall be large enough for a dirent with an
// array of char d_name members containing at least {NAME_MAX}+1 elements.
//
// Upon successful completion, readdir() shall return a pointer to an object of
// type struct dirent. When an error is encountered, a null pointer shall be
// returned and errno shall be set to indicate the error. When the end of the
// directory is encountered, a null pointer shall be returned and errno is not
// changed.
//
// Note: only the d_name and d_ino fields of struct dirent are valid.
func Xreaddir(tls TLS, dirp uintptr) uintptr {
	if *(*int)(unsafe.Pointer(dirp + unsafe.Offsetof(S__dirstream{}.n))) == 0 {
		return 0
	}

	*(*int)(unsafe.Pointer(dirp + unsafe.Offsetof(S__dirstream{}.n)))--
	r := *(*uintptr)(unsafe.Pointer(dirp + unsafe.Offsetof(S__dirstream{}.next)))
	*(*uintptr)(unsafe.Pointer(dirp + unsafe.Offsetof(S__dirstream{}.next))) += unsafe.Sizeof(syscall.Dirent{})
	return r
}

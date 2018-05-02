// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"syscall"
)

// int fchmod(int fd, mode_t mode);
func Xfchmod(tls TLS, fd int32, mode uint32) int32 {
	panic("TODO fchmod")
}

// int mkdir(const char *pathname, mode_t mode);
func Xmkdir(tls TLS, pathname uintptr, mode uint32) int32 {
	// The mkdir() function shall create a new directory with name path.
	// The file permission bits of the new directory shall be initialized
	// from mode. These file permission bits of the mode argument shall be
	// modified by the process' file creation mask.
	//
	// When bits in mode other than the file permission bits are set, the
	// meaning of these additional bits is implementation-defined.
	//
	// The directory's user ID shall be set to the process' effective user
	// ID. The directory's group ID shall be set to the group ID of the
	// parent directory or to the effective group ID of the process.
	// Implementations shall provide a way to initialize the directory's
	// group ID to the group ID of the parent directory. Implementations
	// may, but need not, provide an implementation-defined way to
	// initialize the directory's group ID to the effective group ID of the
	// calling process.
	//
	// The newly created directory shall be an empty directory.
	//
	// If path names a symbolic link, mkdir() shall fail and set errno to
	// [EEXIST].
	//
	// Upon successful completion, mkdir() shall mark for update the last
	// data access, last data modification, and last file status change
	// timestamps of the directory. Also, the last data modification and
	// last file status change timestamps of the directory that contains
	// the new entry shall be marked for update.
	//
	//
	// Upon successful completion, these functions shall return 0.
	// Otherwise, these functions shall return -1 and set errno to indicate
	// the error. If -1 is returned, no directory shall be created.
	r, _, err := syscall.Syscall(syscall.SYS_MKDIR, pathname, uintptr(mode), 0)
	if strace {
		fmt.Fprintf(TraceWriter, "mkdir(%q, %#o) %v %v\n", GoString(pathname), mode, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int chmod(const char *path, mode_t mode);
//
// The chmod() function shall change S_ISUID, S_ISGID, [XSI] ￼  S_ISVTX, ￼ and
// the file permission bits of the file named by the pathname pointed to by the
// path argument to the corresponding bits in the mode argument. The
// application shall ensure that the effective user ID of the process matches
// the owner of the file or the process has appropriate privileges in order to
// do this.
//
// S_ISUID, S_ISGID, [XSI] ￼  S_ISVTX, ￼ and the file permission bits are
// described in <sys/stat.h>.
//
// If the calling process does not have appropriate privileges, and if the
// group ID of the file does not match the effective group ID or one of the
// supplementary group IDs and if the file is a regular file, bit S_ISGID
// (set-group-ID on execution) in the file's mode shall be cleared upon
// successful return from chmod().
//
// Additional implementation-defined restrictions may cause the S_ISUID and
// S_ISGID bits in mode to be ignored.
//
// Upon successful completion, chmod() shall mark for update the last file
// status change timestamp of the file.
//
// Upon successful completion, these functions shall return 0. Otherwise, these
// functions shall return -1 and set errno to indicate the error. If -1 is
// returned, no change to the file mode occurs.
func Xchmod(tls TLS, path uintptr, mode uint32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_CHMOD, path, uintptr(mode), 0)
	if strace {
		fmt.Fprintf(TraceWriter, "chmod(%q, %#o) %v %v\n", GoString(path), mode, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int mkfifo(const char *pathname, mode_t mode);
func Xmkfifo(tls TLS, pathname uintptr, mode uint32) int32 {
	panic("TODO")
}

// mode_t umask(mode_t cmask);
func Xumask(tls TLS, cmask uint32) uint32 {
	// The umask() function shall set the file mode creation mask of the
	// process to cmask and return the previous value of the mask. Only the
	// file permission bits of cmask (see <sys/stat.h>) are used; the
	// meaning of the other bits is implementation-defined.
	//
	// The file mode creation mask of the process is used to turn off
	// permission bits in the mode argument supplied during calls to the
	// following functions:
	//
	// open(), openat(), creat(), mkdir(), mkdirat(), mkfifo(), and
	// mkfifoat()
	//
	// [XSI] ￼ mknod(), mknodat() ￼
	//
	// [MSG] ￼ mq_open() ￼
	//
	// sem_open()
	//
	// Bit positions that are set in cmask are cleared in the mode of the
	// created file.
	//
	// The file permission bits in the value returned by umask() shall be
	// the previous value of the file mode creation mask. The state of any
	// other bits in that value is unspecified, except that a subsequent
	// call to umask() with the returned value as cmask shall leave the
	// state of the mask the same as its state before the first call,
	// including any unspecified use of those bits.
	r, _, err := syscall.Syscall(syscall.SYS_UMASK, uintptr(cmask), 0, 0)
	if strace {
		fmt.Fprintf(TraceWriter, "umask(%#o) %#o %v\n", cmask, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return uint32(r)
}

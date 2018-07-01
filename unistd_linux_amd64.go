// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"syscall"
)

// off_t lseek(int fd, off_t offset, int whence);
func Xlseek(tls TLS, fd int32, offset int64, whence int32) int64 {
	return Xlseek64(tls, fd, offset, whence)
}

// int ftruncate(int fildes, off_t length);
func Xftruncate(tls TLS, fildes int32, length int64) int32 {
	panic("TODO")
}

// int ftruncate64(int fildes, off64_t length);
//
// If fildes is not a valid file descriptor open for writing, the ftruncate()
// function shall fail.
//
// If fildes refers to a regular file, the ftruncate() function shall cause the
// size of the file to be truncated to length. If the size of the file
// previously exceeded length, the extra data shall no longer be available to
// reads on the file. If the file previously was smaller than this size,
// ftruncate() shall increase the size of the file. If the file size is
// increased, the extended area shall appear as if it were zero-filled. The
// value of the seek pointer shall not be modified by a call to ftruncate().
//
// Upon successful completion, if fildes refers to a regular file, ftruncate()
// shall mark for update the last data modification and last file status change
// timestamps of the file and the S_ISUID and S_ISGID bits of the file mode may
// be cleared. If the ftruncate() function is unsuccessful, the file is
// unaffected.
//
// [XSI] ￼ If the request would cause the file size to exceed the soft file
// size limit for the process, the request shall fail and the implementation
// shall generate the SIGXFSZ signal for the thread. ￼
//
// If fildes refers to a directory, ftruncate() shall fail.
//
// If fildes refers to any other file type, except a shared memory object, the
// result is unspecified.
//
// [SHM] ￼ If fildes refers to a shared memory object, ftruncate() shall set
// the size of the shared memory object to length. ￼
//
// If the effect of ftruncate() is to decrease the size of a memory mapped file
// [SHM] ￼  or a shared memory object ￼  and whole pages beyond the new end
// were previously mapped, then the whole pages beyond the new end shall be
// discarded.
//
// References to discarded pages shall result in the generation of a SIGBUS
// signal.
//
// If the effect of ftruncate() is to increase the size of a memory object, it
// is unspecified whether the contents of any mapped pages between the old
// end-of-file and the new are flushed to the underlying object.
//
// Upon successful completion, ftruncate() shall return 0; otherwise, -1 shall
// be returned and errno set to indicate the error.
func Xftruncate64(tls TLS, fildes int32, length int64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FTRUNCATE, uintptr(fildes), uintptr(length), 0)
	if strace {
		fmt.Fprintf(TraceWriter, "ftruncate64(%#x, %#x) %v, %v\n", fildes, length, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

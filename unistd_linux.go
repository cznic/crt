// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/cznic/crt/errno"
	"github.com/cznic/crt/unistd"
	"golang.org/x/crypto/ssh/terminal"
)

// int close(int fd);
func Xclose(tls TLS, fd int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_CLOSE, uintptr(fd), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "close(%v) %v %v\n", fd, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	openFDsMu.Lock()
	delete(openFDs, fd)
	openFDsMu.Unlock()
	return int32(r)
}

// int access(const char *path, int amode);
func Xaccess(tls TLS, path uintptr, amode int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_ACCESS, path, uintptr(amode), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "access(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int unlink(const char *path);
func Xunlink(tls TLS, path uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_UNLINK, path, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "unlink(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int rmdir(const char *path);
//
// The rmdir() function shall remove a directory whose name is given by path.
// The directory shall be removed only if it is an empty directory.
//
// If the directory is the root directory or the current working directory of
// any process, it is unspecified whether the function succeeds, or whether it
// shall fail and set errno to [EBUSY].
//
// If path names a symbolic link, then rmdir() shall fail and set errno to
// [ENOTDIR].
//
// If the path argument refers to a path whose final component is either dot or
// dot-dot, rmdir() shall fail.
//
// If the directory's link count becomes 0 and no process has the directory
// open, the space occupied by the directory shall be freed and the directory
// shall no longer be accessible. If one or more processes have the directory
// open when the last link is removed, the dot and dot-dot entries, if present,
// shall be removed before rmdir() returns and no new entries may be created in
// the directory, but the directory shall not be removed until all references
// to the directory are closed.
//
// If the directory is not an empty directory, rmdir() shall fail and set errno
// to [EEXIST] or [ENOTEMPTY].
//
// Upon successful completion, rmdir() shall mark for update the last data
// modification and last file status change timestamps of the parent directory.
//
// Upon successful completion, the function rmdir() shall return 0. Otherwise,
// -1 shall be returned, and errno set to indicate the error. If -1 is
// returned, the named directory shall not be changed.
func Xrmdir(tls TLS, path uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_RMDIR, path, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "rmdir(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int chown(const char *pathname, uid_t owner, gid_t group);
func Xchown(tls TLS, pathname uintptr, owner, group uint32) int32 {
	panic("TODO chown")
}

// int fchown(int fd, uid_t owner, gid_t group);
func Xfchown(tls TLS, fd int32, owner, group uint32) int32 {
	panic("TODO fchown")
}

// uid_t getuid(void);
func Xgetuid(tls TLS) uint32 {
	r, _, err := syscall.RawSyscall(syscall.SYS_GETUID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getuid() %v\n", r)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return uint32(r)
}

// uid_t geteuid(void);
func Xgeteuid(tls TLS) uint32 {
	r, _, err := syscall.RawSyscall(syscall.SYS_GETEUID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "geteuid() %v\n", r)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return uint32(r)
}

// int fsync(int fildes);
func Xfsync(tls TLS, fildes int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FSYNC, uintptr(fildes), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fsync(%v) %v %v\n", fildes, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int fdatasync(int fd);
func Xfdatasync(tls TLS, fildes int32) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FDATASYNC, uintptr(fildes), 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fdatasync(%v) %v %v\n", fildes, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// pid_t getpid(void);
func Xgetpid(tls TLS) int32 {
	r, _, err := syscall.RawSyscall(syscall.SYS_GETPID, 0, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getpid() %v\n", r)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// unsigned sleep(unsigned seconds);
func Xsleep(tls TLS, seconds uint32) uint32 {
	time.Sleep(time.Duration(seconds) * time.Second)
	if strace {
		fmt.Fprintf(os.Stderr, "sleep(%#x)", seconds)
	}
	return 0
}

// off64_t lseek64(int fd, off64_t offset, int whence);
func Xlseek64(tls TLS, fd int32, offset int64, whence int32) int64 {
	r, _, err := syscall.Syscall(syscall.SYS_LSEEK, uintptr(fd), uintptr(offset), uintptr(whence))
	if strace {
		fmt.Fprintf(os.Stderr, "lseek64(%v, %v, %v) %v %v\n", fd, offset, whence, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int64(r)
}

// int usleep(useconds_t usec);
func Xusleep(tls TLS, usec uint32) int32 {
	time.Sleep(time.Duration(usec) * time.Microsecond)
	if strace {
		fmt.Fprintf(os.Stderr, "usleep(%#x)", usec)
	}
	return 0
}

// int chdir(const char *path);
func Xchdir(tls TLS, path uintptr) int32 {
	// The chdir() function shall cause the directory named by the pathname
	// pointed to by the path argument to become the current working
	// directory; that is, the starting point for path searches for
	// pathnames not beginning with '/'.
	//
	// Upon successful completion, 0 shall be returned. Otherwise, -1 shall
	// be returned, the current working directory shall remain unchanged,
	// and errno shall be set to indicate the error.
	r, _, err := syscall.Syscall(syscall.SYS_CHDIR, path, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "chdir(%q) %v %v\n", GoString(path), r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// ssize_t read(int fd, void *buf, size_t count);
func Xread(tls TLS, fd int32, buf uintptr, count size_t) ssize_t { //TODO stdin
	r, _, err := syscall.Syscall(syscall.SYS_READ, uintptr(fd), buf, uintptr(count))
	if strace {
		fmt.Fprintf(os.Stderr, "read(%v, %#x, %v) %v %v\n", fd, buf, count, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return ssize_t(r)
}

// char *getcwd(char *buf, size_t size);
func Xgetcwd(tls TLS, buf uintptr, size size_t) uintptr {
	r, _, err := syscall.Syscall(syscall.SYS_GETCWD, buf, uintptr(size), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "getcwd(%#x, %#x) %v %v %q\n", buf, size, r, err, GoString(buf))
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return r
}

// ssize_t write(int fd, const void *buf, size_t count);
func Xwrite(tls TLS, fd int32, buf uintptr, count size_t) ssize_t {
	switch fd {
	case unistd.XSTDIN_FILENO:
		panic("TODO")
	case unistd.XSTDOUT_FILENO:
		n, err := os.Stdout.Write((*rawmem)(unsafe.Pointer(buf))[:count])
		if err != nil {
			tls.setErrno(err)
		}
		return ssize_t(n)
	case unistd.XSTDERR_FILENO:
		n, err := os.Stderr.Write((*rawmem)(unsafe.Pointer(buf))[:count])
		if err != nil {
			tls.setErrno(err)
		}
		return ssize_t(n)
	}
	r, _, err := syscall.Syscall(syscall.SYS_WRITE, uintptr(fd), buf, uintptr(count))
	if strace {
		fmt.Fprintf(os.Stderr, "write(%v, %#x, %v) %v %v\n", fd, buf, count, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return ssize_t(r)
}

// ssize_t readlink(const char *pathname, char *buf, size_t bufsiz);
//
// The readlink() function shall place the contents of the symbolic link
// referred to by path in the buffer buf which has size bufsize. If the number
// of bytes in the symbolic link is less than bufsize, the contents of the
// remainder of buf are unspecified. If the buf argument is not large enough to
// contain the link content, the first bufsize bytes shall be placed in buf.
//
// If the value of bufsize is greater than {SSIZE_MAX}, the result is
// implementation-defined.
//
// Upon successful completion, readlink() shall mark for update the last data
// access timestamp of the symbolic link.
//
// Upon successful completion, these functions shall return the count of bytes
// placed in the buffer. Otherwise, these functions shall return a value of -1,
// leave the buffer unchanged, and set errno to indicate the error.
func Xreadlink(tls TLS, pathname, buf uintptr, bufsiz size_t) ssize_t {
	r, _, err := syscall.Syscall(syscall.SYS_READLINK, pathname, buf, uintptr(bufsiz))
	if strace {
		fmt.Fprintf(os.Stderr, "readlink(%q, %#x, %#x) %v %v\n", GoString(pathname), buf, bufsiz, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return ssize_t(r)
}

// long sysconf(int name);
func Xsysconf(tls TLS, name int32) int64 {
	switch name {
	case unistd.C_SC_PAGESIZE:
		return int64(os.Getpagesize())
	case unistd.C_SC_GETPW_R_SIZE_MAX:
		return -1
	default:
		panic(fmt.Errorf("%v(%#x)", name, name))
	}
}

// int isatty(int fd);
func Xisatty(tls TLS, fd int32) int32 {
	if terminal.IsTerminal(int(fd)) {
		return 1
	}

	tls.setErrno(errno.XENOTTY)
	return 0
}

// int symlink(const char *target, const char *linkpath);
func Xsymlink(tls TLS, target, linkpath uintptr) int32 {
	panic("TODO")
}

// int mknod(const char *pathname, mode_t mode, dev_t dev);
func Xmknod(tls TLS, pathname uintptr, mode uint32, dev uint64) int32 {
	panic("TODO")
}

// int link(const char *oldpath, const char *newpath);
func Xlink(tls TLS, oldpath, newpath uintptr) int32 {
	panic("TODO")
}

// int pipe(int fildes[2]);
//
// The pipe() function shall create a pipe and place two file descriptors, one
// each into the arguments fildes[0] and fildes[1], that refer to the open file
// descriptions for the read and write ends of the pipe. The file descriptors
// shall be allocated as described in File Descriptor Allocation. The
// O_NONBLOCK and FD_CLOEXEC flags shall be clear on both file descriptors.
// (The fcntl() function can be used to set both these flags.)
//
// Data can be written to the file descriptor fildes[1] and read from the file
// descriptor fildes[0]. A read on the file descriptor fildes[0] shall access
// data written to the file descriptor fildes[1] on a first-in-first-out basis.
// It is unspecified whether fildes[0] is also open for writing and whether
// fildes[1] is also open for reading.
//
// A process has the pipe open for reading (correspondingly writing) if it has
// a file descriptor open that refers to the read end, fildes[0] (write end,
// fildes[1]).
//
// The pipe's user ID shall be set to the effective user ID of the calling
// process.
//
// The pipe's group ID shall be set to the effective group ID of the calling
// process.
//
// Upon successful completion, pipe() shall mark for update the last data
// access, last data modification, and last file status change timestamps of
// the pipe.
//
// Upon successful completion, 0 shall be returned; otherwise, -1 shall be
// returned and errno set to indicate the error, no file descriptors shall be
// allocated and the contents of fildes shall be left unmodified.
func Xpipe(tls TLS, fildes uintptr) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_PIPE, fildes, 0, 0)
	if strace {
		fmt.Fprintf(os.Stderr, "pipe(%#x) %v %v\n", fildes, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int pthread_cond_timedwait(pthread_cond_t *restrict cond, pthread_mutex_t *restrict mutex, const struct timespec *restrict abstime);
func Xpthread_cond_timedwait(tls TLS, cond, mutex, abstime uintptr) int32 {
	panic("TODO")
}

// pid_t fork(void);
//
// The fork() function shall create a new process. The new process (child
// process) shall be an exact copy of the calling process (parent process)
// except as detailed below:
//
// The child process shall have a unique process ID.
//
// The child process ID also shall not match any active process group ID.
//
// The child process shall have a different parent process ID, which shall be
// the process ID of the calling process.
//
// The child process shall have its own copy of the parent's file descriptors.
// Each of the child's file descriptors shall refer to the same open file
// description with the corresponding file descriptor of the parent.
//
// The child process shall have its own copy of the parent's open directory
// streams. Each open directory stream in the child process may share directory
// stream positioning with the corresponding directory stream of the parent.
//
// The child process shall have its own copy of the parent's message catalog
// descriptors.
//
// The child process values of tms_utime, tms_stime, tms_cutime, and tms_cstime
// shall be set to 0.
//
// The time left until an alarm clock signal shall be reset to zero, and the
// alarm, if any, shall be canceled; see alarm.
//
// [XSI] ￼ All semadj values shall be cleared. ￼
//
// File locks set by the parent process shall not be inherited by the child
// process.
//
// The set of signals pending for the child process shall be initialized to the
// empty set.
//
// [XSI] ￼ Interval timers shall be reset in the child process. ￼
//
// Any semaphores that are open in the parent process shall also be open in the
// child process.
//
// [ML] ￼ The child process shall not inherit any address space memory locks
// established by the parent process via calls to mlockall() or mlock(). ￼
//
// Memory mappings created in the parent shall be retained in the child
// process. MAP_PRIVATE mappings inherited from the parent shall also be
// MAP_PRIVATE mappings in the child, and any modifications to the data in
// these mappings made by the parent prior to calling fork() shall be visible
// to the child. Any modifications to the data in MAP_PRIVATE mappings made by
// the parent after fork() returns shall be visible only to the parent.
// Modifications to the data in MAP_PRIVATE mappings made by the child shall be
// visible only to the child.
//
// [PS] ￼ For the SCHED_FIFO and SCHED_RR scheduling policies, the child
// process shall inherit the policy and priority settings of the parent process
// during a fork() function. For other scheduling policies, the policy and
// priority settings on fork() are implementation-defined. ￼
//
// Per-process timers created by the parent shall not be inherited by the child
// process.
//
// [MSG] ￼ The child process shall have its own copy of the message queue
// descriptors of the parent. Each of the message descriptors of the child
// shall refer to the same open message queue description as the corresponding
// message descriptor of the parent. ￼
//
// No asynchronous input or asynchronous output operations shall be inherited
// by the child process. Any use of asynchronous control blocks created by the
// parent produces undefined behavior.
//
// A process shall be created with a single thread. If a multi-threaded process
// calls fork(), the new process shall contain a replica of the calling thread
// and its entire address space, possibly including the states of mutexes and
// other resources. Consequently, to avoid errors, the child process may only
// execute async-signal-safe operations until such time as one of the exec
// functions is called.
//
// When the application calls fork() from a signal handler and any of the fork
// handlers registered by pthread_atfork() calls a function that is not
// async-signal-safe, the behavior is undefined.
//
// [OB TRC TRI] ￼ If the Trace option and the Trace Inherit option are both
// supported:
//
// If the calling process was being traced in a trace stream that had its
// inheritance policy set to POSIX_TRACE_INHERITED, the child process shall be
// traced into that trace stream, and the child process shall inherit the
// parent's mapping of trace event names to trace event type identifiers. If
// the trace stream in which the calling process was being traced had its
// inheritance policy set to POSIX_TRACE_CLOSE_FOR_CHILD, the child process
// shall not be traced into that trace stream. The inheritance policy is set by
// a call to the posix_trace_attr_setinherited() function. ￼
//
// [OB TRC] ￼ If the Trace option is supported, but the Trace Inherit option is
// not supported:
//
// The child process shall not be traced into any of the trace streams of its
// parent process. ￼
//
// [OB TRC] ￼ If the Trace option is supported, the child process of a trace
// controller process shall not control the trace streams controlled by its
// parent process. ￼
//
// [CPT] ￼ The initial value of the CPU-time clock of the child process shall
// be set to zero. ￼
//
// [TCT] ￼ The initial value of the CPU-time clock of the single thread of the
// child process shall be set to zero. ￼
//
// All other process characteristics defined by POSIX.1-2017 shall be the same
// in the parent and child processes. The inheritance of process
// characteristics not defined by POSIX.1-2017 is unspecified by POSIX.1-2017.
//
// After fork(), both the parent and the child processes shall be capable of
// executing independently before either one terminates.
//
// Upon successful completion, fork() shall return 0 to the child process and
// shall return the process ID of the child process to the parent process. Both
// processes shall continue to execute from the fork() function. Otherwise, -1
// shall be returned to the parent process, no child process shall be created,
// and errno shall be set to indicate the error.
//
// Note: Go programs cannot fork safely. The function sets errno to
// errno.XENOMEM and returns -1.
func Xfork(tls TLS) int32 {
	tls.setErrno(errno.XENOMEM)
	return -1
}

// int dup2(int oldfd, int newfd);
func Xdup2(tls TLS, oldfd, newfd int32) int32 {
	panic("TODO")
}

// void _exit(int status);
func X_exit(tls TLS, status int32) {
	X__builtin_exit(tls, status)
}

// int execvp(const char *file, char *const argv[]);
func Xexecvp(tls TLS, file, argv uintptr) int32 {
	panic("TODO")
}

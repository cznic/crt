// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//TODO go:generate go run generator.go -goos linux -goarch 386
//TODO go:generate go run generator.go -goos linux -goarch amd64
//TODO go:generate go run generator.go -goos linux -goarch arm
//TODO go:generate go run generator.go -goos windows -goarch 386
//TODO go:generate go run generator.go -goos windows -goarch amd64

// Package crt provides C-runtime services. Work In Progress. API unstable.
//
// Installation
//
//     $ go get github.com/cznic/crt
//
// Documentation: http://godoc.org/github.com/cznic/crt
//
// The vast majority of this package is a mechanical translation of the musl
// libc project:
//
//	https://www.musl-libc.org/
package crt

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/cznic/internal/buffer" //TODO-
	"github.com/cznic/memory"
	"github.com/cznic/strutil"
)

var (
	allocMu   sync.Mutex
	allocator memory.Allocator
	allocs    = map[uintptr][]byte{}
	env       = os.Environ()
	threadID  uintptr //TODO-

	// RepositoryPath is the path of this repository.
	RepositoryPath string

	// Nz32 holds the float32 value -0.0. R/O
	Nz32 float32

	// Nz64 holds the float64 value -0.0. R/O
	Nz64 float64
)

func init() {
	Nz32 = -Nz32
	Nz64 = -Nz64

	X__libc_start_main(0)

	self, err := strutil.ImportPath()
	if err != nil {
		panic(err)
	}

	gopath := strutil.Gopath()
	for _, v := range strings.Split(gopath, string(os.PathListSeparator)) {
		x := filepath.Join(v, "src", self)
		fi, err := os.Stat(x)
		if err != nil || !fi.IsDir() {
			continue
		}

		RepositoryPath = x
		return
	}

	panic("cannot determine repository path")
}

func Main(main func(TLS, int32, uintptr) int32) {
	tls := MainTLS()
	Xexit(tls, main(tls, int32(len(os.Args)), *(*uintptr)(unsafe.Pointer(X__ccgo_argv))))
}

func MainTLS() TLS { return TLS(*(*uintptr)(unsafe.Pointer(X__ccgo_main_tls))) }

// TLS represents a virtual C thread.
type TLS uintptr

// MustMalloc is like Malloc but panics if the allocation cannot be made.
func MustMalloc(size int) uintptr {
	p, err := Malloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// UsableSize reports the size of the memory block allocated at p.
func UsableSize(p uintptr) int {
	allocMu.Lock()
	n := memory.UintptrUsableSize(p)
	allocMu.Unlock()
	return n
}

// Malloc allocates uninitialized memory.
func Malloc(size int) (uintptr, error) {
	if size < 0 {
		panic("internal error")
	}

	if size == 0 {
		return 0, nil
	}

	allocMu.Lock()
	r, _ := malloc(size)
	allocMu.Unlock()
	return r, nil
}

func malloc(size int) (uintptr, error) {
	b := make([]byte, size+16)
	r := uintptr(unsafe.Pointer(&b[0]))
	if r%2*unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}
	if _, ok := allocs[r]; ok {
		panic("internal error")
	}
	allocs[r] = b
	return r, nil
}

// MustCalloc is like Calloc but panics if the allocation cannot be made.
func MustCalloc(size int) uintptr {
	p, err := Calloc(size)
	if err != nil {
		panic(fmt.Errorf("out of memory: %v", err))
	}

	return p
}

// Calloc allocates zeroed memory.
func Calloc(size int) (uintptr, error) { return Malloc(size) }

// Realloc reallocates memory.
func Realloc(p uintptr, size int) (uintptr, error) {
	if p == 0 {
		return Malloc(size)
	}

	if size == 0 {
		Free(p)
		return 0, nil
	}

	allocMu.Lock()
	b, ok := allocs[p]
	if !ok {
		panic("internal error")
	}

	switch {
	case cap(b) >= size:
		b = b[:size]
		allocs[p] = b
	default:
		r, _ := malloc(size)
		copy(allocs[r], b)
		free(p)
		p = r
	}
	allocMu.Unlock()
	return p, nil
}

// Free frees memory allocated by Calloc, Malloc or Realloc.
func Free(p uintptr) error {
	allocMu.Lock()
	err := free(p)
	allocMu.Unlock()
	return err
}

func X__builtin_free(tls TLS, p uintptr) { Xfree(tls, p) }

func free(p uintptr) error {
	if p == 0 {
		return nil
	}

	if _, ok := allocs[p]; !ok {
		panic("internal error")
	}

	delete(allocs, p)
	return nil
}

// CString allocates a C string initialized from s.
func CString(s string) (uintptr, error) {
	n := len(s)
	p, err := Malloc(n + 1)
	if p == 0 || err != nil {
		return 0, err
	}

	copy((*rawmem)(unsafe.Pointer(p))[:n], s)
	(*rawmem)(unsafe.Pointer(p))[n] = 0
	return p, nil
}

// MustCString is like CString but panic if the allocation cannot be made.
func MustCString(s string) uintptr {
	n := len(s)
	p := MustMalloc(n + 1)
	copy((*rawmem)(unsafe.Pointer(p))[:n], s)
	(*rawmem)(unsafe.Pointer(p))[n] = 0
	return p
}

// BSS allocates the the bss segment of a package/command.
func BSS(init *byte) uintptr {
	r := uintptr(unsafe.Pointer(init))
	if r%2*unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// TS allocates the R/O text segment of a package/command.
func TS(init string) uintptr { return (*reflect.StringHeader)(unsafe.Pointer(&init)).Data }

// DS allocates the the data segment of a package/command.
func DS(init []byte) uintptr {
	r := (*reflect.SliceHeader)(unsafe.Pointer(&init)).Data
	if r%2*unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// Copy copies n bytes form src to dest and returns n.
func Copy(dst, src uintptr, n int) int { //TODO-
	if n != 0 {
		return copy((*rawmem)(unsafe.Pointer(dst))[:n], (*rawmem)(unsafe.Pointer(src))[:n])
	}
	return n
}

// GoString returns a string from a C char* null terminated string s.
func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}

	var b buffer.Bytes
	for {
		ch := *(*byte)(unsafe.Pointer(s))
		if ch == 0 {
			r := string(b.Bytes())
			b.Close()
			return r
		}

		b.WriteByte(ch)
		s++
	}
}

func printAssertFail(expr, file uintptr, line int32, _func uintptr) {
	fmt.Fprintf(os.Stderr, "Assertion failed: %s (%s: %s: %d)\n", GoString(expr), GoString(file), GoString(_func), line)
	fmt.Fprintf(os.Stderr, "%s\n", debug.Stack()) //TODO-
	os.Exit(1)                                    //TODO-
}

var barier int32

func aBarier() {
	atomic.LoadInt32(&barier)
}

func Alloca(p *[]uintptr, n int) uintptr   { r := MustMalloc(n); *p = append(*p, r); return r }
func Preinc(p *uintptr, n uintptr) uintptr { *p += n; return *p }

var globalMutex sync.Mutex

// static inline int a_cas(volatile int *p, int t, int s)
func a_cas(p uintptr, t, s int32) int32 {
	globalMutex.Lock()
	old := *(*int32)(unsafe.Pointer(p))
	if *(*int32)(unsafe.Pointer(p)) == t {
		*(*int32)(unsafe.Pointer(p)) = s
	}
	globalMutex.Unlock()
	return old
}

// static inline void *a_cas_p(volatile void *p, void *t, void *s)
func a_cas_p(p, t, s uintptr) uintptr {
	globalMutex.Lock()
	old := *(*uintptr)(unsafe.Pointer(p))
	if *(*uintptr)(unsafe.Pointer(p)) == t {
		*(*uintptr)(unsafe.Pointer(p)) = s
	}
	globalMutex.Unlock()
	return old
}

//static inline void a_or_64(volatile uint64_t *p, uint64_t v)
func a_or_64(p uintptr, v uint64) {
	globalMutex.Lock()
	*(*uint64)(unsafe.Pointer(p)) |= v
	globalMutex.Unlock()
}

//static inline void a_and_64(volatile uint64_t *p, uint64_t v)
func a_and_64(p uintptr, v uint64) {
	globalMutex.Lock()
	*(*uint64)(unsafe.Pointer(p)) &= v
	globalMutex.Unlock()
}

// static inline void a_inc(volatile int *p)
func a_inc(p uintptr) {
	atomic.AddInt32((*int32)(unsafe.Pointer(p)), 1)
}

// static inline void a_dec(volatile int *p)
func a_dec(p uintptr) {
	atomic.AddInt32((*int32)(unsafe.Pointer(p)), -1)
}

//static inline int a_fetch_add(volatile int *p, int v)
func a_fetch_add(p uintptr, v int32) int32 {
	return atomic.AddInt32((*int32)(unsafe.Pointer(p)), v)
}

func debugStack() { fmt.Printf("%s\n", debug.Stack()) }

type int16watch = struct {
	s string
	v int16
}

type int32watch = struct {
	s string
	v int32
}

type int64watch = struct {
	s string
	v int64
}

type uint32watch = struct {
	s string
	v uint32
}

type ptrwatch = struct {
	s string
	v uintptr
}

type strwatch = struct {
	s string
	v string
}

var ( //TODO-
	watches = map[uintptr]interface{}{}
	trace   int
)

func WatchPtr(tls TLS, s string, p uintptr) {
	if s == "" {
		delete(watches, p)
		return
	}

	v := *(*uintptr)(unsafe.Pointer(p))
	watches[p] = &ptrwatch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %#x)", s, p, v)})
}

func WatchInt16(tls TLS, s string, p uintptr) {
	if s == "" {
		delete(watches, p)
		return
	}

	v := *(*int16)(unsafe.Pointer(p))
	watches[p] = &int16watch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %v(%#[3]x)", s, p, v)})
}

func WatchInt32(tls TLS, s string, p uintptr) {
	if s == "" {
		delete(watches, p)
		return
	}

	v := *(*int32)(unsafe.Pointer(p))
	watches[p] = &int32watch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %v(%#[3]x)", s, p, v)})
}

func WatchInt64(tls TLS, s string, p uintptr) {
	if s == "" {
		delete(watches, p)
		return
	}

	v := *(*int64)(unsafe.Pointer(p))
	watches[p] = &int64watch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %v(%#[3]x)", s, p, v)})
}

func WatchString(tls TLS, s string, p uintptr) {
	if s == "" {
		delete(watches, p)
		return
	}

	v := GoString(p)
	watches[p] = &strwatch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %q", s, p, v)})
}

func WatchUint32(tls TLS, s string, p uintptr) {
	if s == "" || p == 0 {
		delete(watches, p)
		return
	}

	v := *(*uint32)(unsafe.Pointer(p))
	watches[p] = &uint32watch{s, v}
	watching(tls, []string{fmt.Sprintf("%q @ %#x is initially %v(%#[3]x)", s, p, v)})
}

func watching(tls TLS, a []string) {
	sort.Strings(a)
	fmt.Printf("==== %s\n%s----\n", strings.Join(a, "\n"), debug.Stack())
	os.Stdout.Sync()
}

func TraceOn()  { trace++ }
func TraceOff() { trace-- }

func Watch(tls TLS) { //TODO-
	if trace > 0 {
		Xfflush(tls, 0)
		_, fn, fl, _ := runtime.Caller(1)
		fmt.Printf("# trace %s:%d:\n", filepath.Base(fn), fl)
	}
	var a []string
	for p, v := range watches {
		switch x := v.(type) {
		case *int16watch:
			if w := *(*int16)(unsafe.Pointer(p)); w != x.v {
				x.v = w
				a = append(a, fmt.Sprintf("%q @ %#x is now %v(%#[3]x)", x.s, p, w))
			}
		case *int32watch:
			if w := *(*int32)(unsafe.Pointer(p)); w != x.v {
				x.v = w
				a = append(a, fmt.Sprintf("%q @ %#x is now %v(%#[3]x)", x.s, p, w))
			}
		case *int64watch:
			if w := *(*int64)(unsafe.Pointer(p)); w != x.v {
				a = append(a, fmt.Sprintf("%q @ %#x is now %v(%#[3]x) was %[4]v(%#[4]x)", x.s, p, w, x.v))
				x.v = w
			}
		case *uint32watch:
			if w := *(*uint32)(unsafe.Pointer(p)); w != x.v {
				x.v = w
				a = append(a, fmt.Sprintf("%q @ %#x is now %v(%#[3]x)", x.s, p, w))
			}
		case *ptrwatch:
			if w := *(*uintptr)(unsafe.Pointer(p)); w != x.v {
				x.v = w
				a = append(a, fmt.Sprintf("%q @ %#x is now %#x)", x.s, p, w))
			}
		case *strwatch:
			if w := GoString(p); w != x.v {
				x.v = w
				a = append(a, fmt.Sprintf("%q @ %#x is now %q)", x.s, p, w))
			}
		}
	}
	if len(a) != 0 {
		watching(tls, a)
	}
}

func checkSyscall(n long) bool { //TODO- eventually after making the C code clean of "bad" syscalls
	switch n {
	case
		DSYS_access,
		DSYS_brk,
		DSYS_chdir,
		DSYS_chmod,
		DSYS_clock_gettime,
		DSYS_close,
		DSYS_fchmod,
		DSYS_fcntl,
		DSYS_fstat,
		DSYS_fsync,
		DSYS_ftruncate,
		DSYS_futex,
		DSYS_getcwd,
		DSYS_getdents,
		DSYS_geteuid,
		DSYS_getpid,
		DSYS_getsockname,
		DSYS_getuid,
		DSYS_ioctl,
		DSYS_lseek,
		DSYS_lstat,
		DSYS_madvise,
		DSYS_mkdir,
		DSYS_mmap,
		DSYS_mprotect,
		DSYS_mremap,
		DSYS_munmap,
		DSYS_nanosleep,
		DSYS_open,
		DSYS_pipe,
		DSYS_pipe2,
		DSYS_pread,
		DSYS_pwrite,
		DSYS_read,
		DSYS_readlink,
		DSYS_readv,
		DSYS_rename,
		DSYS_rmdir,
		DSYS_select,
		DSYS_stat,
		DSYS_symlink,
		DSYS_umask,
		DSYS_uname,
		DSYS_unlink,
		DSYS_utimensat,
		DSYS_wait4,
		DSYS_write,
		DSYS_writev:

		return true
	case
		DSYS_rt_sigaction,   //TODO
		DSYS_rt_sigprocmask: //TODO

		return false // ignore the syscall //TODO later
	default:
		panic(n)
	}
}

var (
	syscalls = map[int]string{
		DSYS__sysctl:                "_sysctl",
		DSYS_accept:                 "accept",
		DSYS_accept4:                "accept4",
		DSYS_access:                 "access",
		DSYS_acct:                   "acct",
		DSYS_add_key:                "add_key",
		DSYS_adjtimex:               "adjtimex",
		DSYS_afs_syscall:            "afs_syscall",
		DSYS_alarm:                  "alarm",
		DSYS_arch_prctl:             "arch_prctl",
		DSYS_bind:                   "bind",
		DSYS_bpf:                    "bpf",
		DSYS_brk:                    "brk",
		DSYS_capget:                 "capget",
		DSYS_capset:                 "capset",
		DSYS_chdir:                  "chdir",
		DSYS_chmod:                  "chmod",
		DSYS_chown:                  "chown",
		DSYS_chroot:                 "chroot",
		DSYS_clock_adjtime:          "clock_adjtime",
		DSYS_clock_getres:           "clock_getres",
		DSYS_clock_gettime:          "clock_gettime",
		DSYS_clock_nanosleep:        "clock_nanosleep",
		DSYS_clock_settime:          "clock_settime",
		DSYS_clone:                  "clone",
		DSYS_close:                  "close",
		DSYS_connect:                "connect",
		DSYS_copy_file_range:        "copy_file_range",
		DSYS_creat:                  "creat",
		DSYS_create_module:          "create_module",
		DSYS_delete_module:          "delete_module",
		DSYS_dup:                    "dup",
		DSYS_dup2:                   "dup2",
		DSYS_dup3:                   "dup3",
		DSYS_epoll_create:           "epoll_create",
		DSYS_epoll_create1:          "epoll_create1",
		DSYS_epoll_ctl:              "epoll_ctl",
		DSYS_epoll_ctl_old:          "epoll_ctl_old",
		DSYS_epoll_pwait:            "epoll_pwait",
		DSYS_epoll_wait:             "epoll_wait",
		DSYS_epoll_wait_old:         "epoll_wait_old",
		DSYS_eventfd:                "eventfd",
		DSYS_eventfd2:               "eventfd2",
		DSYS_execve:                 "execve",
		DSYS_execveat:               "execveat",
		DSYS_exit:                   "exit",
		DSYS_exit_group:             "exit_group",
		DSYS_faccessat:              "faccessat",
		DSYS_fadvise:                "fadvise",
		DSYS_fallocate:              "fallocate",
		DSYS_fanotify_init:          "fanotify_init",
		DSYS_fanotify_mark:          "fanotify_mark",
		DSYS_fchdir:                 "fchdir",
		DSYS_fchmod:                 "fchmod",
		DSYS_fchmodat:               "fchmodat",
		DSYS_fchown:                 "fchown",
		DSYS_fchownat:               "fchownat",
		DSYS_fcntl:                  "fcntl",
		DSYS_fdatasync:              "fdatasync",
		DSYS_fgetxattr:              "fgetxattr",
		DSYS_finit_module:           "finit_module",
		DSYS_flistxattr:             "flistxattr",
		DSYS_flock:                  "flock",
		DSYS_fork:                   "fork",
		DSYS_fremovexattr:           "fremovexattr",
		DSYS_fsetxattr:              "fsetxattr",
		DSYS_fstat:                  "fstat",
		DSYS_fstatat:                "fstatat",
		DSYS_fstatfs:                "fstatfs",
		DSYS_fsync:                  "fsync",
		DSYS_ftruncate:              "ftruncate",
		DSYS_futex:                  "futex",
		DSYS_futimesat:              "futimesat",
		DSYS_get_kernel_syms:        "get_kernel_syms",
		DSYS_get_mempolicy:          "get_mempolicy",
		DSYS_get_robust_list:        "get_robust_list",
		DSYS_get_thread_area:        "get_thread_area",
		DSYS_getcpu:                 "getcpu",
		DSYS_getcwd:                 "getcwd",
		DSYS_getdents:               "getdents",
		DSYS_getegid:                "getegid",
		DSYS_geteuid:                "geteuid",
		DSYS_getgid:                 "getgid",
		DSYS_getgroups:              "getgroups",
		DSYS_getitimer:              "getitimer",
		DSYS_getpeername:            "getpeername",
		DSYS_getpgid:                "getpgid",
		DSYS_getpgrp:                "getpgrp",
		DSYS_getpid:                 "getpid",
		DSYS_getpmsg:                "getpmsg",
		DSYS_getppid:                "getppid",
		DSYS_getpriority:            "getpriority",
		DSYS_getrandom:              "getrandom",
		DSYS_getresgid:              "getresgid",
		DSYS_getresuid:              "getresuid",
		DSYS_getrlimit:              "getrlimit",
		DSYS_getrusage:              "getrusage",
		DSYS_getsid:                 "getsid",
		DSYS_getsockname:            "getsockname",
		DSYS_getsockopt:             "getsockopt",
		DSYS_gettid:                 "gettid",
		DSYS_gettimeofday:           "gettimeofday",
		DSYS_getuid:                 "getuid",
		DSYS_getxattr:               "getxattr",
		DSYS_init_module:            "init_module",
		DSYS_inotify_add_watch:      "inotify_add_watch",
		DSYS_inotify_init:           "inotify_init",
		DSYS_inotify_init1:          "inotify_init1",
		DSYS_inotify_rm_watch:       "inotify_rm_watch",
		DSYS_io_cancel:              "io_cancel",
		DSYS_io_destroy:             "io_destroy",
		DSYS_io_getevents:           "io_getevents",
		DSYS_io_setup:               "io_setup",
		DSYS_io_submit:              "io_submit",
		DSYS_ioctl:                  "ioctl",
		DSYS_ioperm:                 "ioperm",
		DSYS_iopl:                   "iopl",
		DSYS_ioprio_get:             "ioprio_get",
		DSYS_ioprio_set:             "ioprio_set",
		DSYS_kcmp:                   "kcmp",
		DSYS_kexec_file_load:        "kexec_file_load",
		DSYS_kexec_load:             "kexec_load",
		DSYS_keyctl:                 "keyctl",
		DSYS_kill:                   "kill",
		DSYS_lchown:                 "lchown",
		DSYS_lgetxattr:              "lgetxattr",
		DSYS_link:                   "link",
		DSYS_linkat:                 "linkat",
		DSYS_listen:                 "listen",
		DSYS_listxattr:              "listxattr",
		DSYS_llistxattr:             "llistxattr",
		DSYS_lookup_dcookie:         "lookup_dcookie",
		DSYS_lremovexattr:           "lremovexattr",
		DSYS_lseek:                  "lseek",
		DSYS_lsetxattr:              "lsetxattr",
		DSYS_lstat:                  "lstat",
		DSYS_madvise:                "madvise",
		DSYS_mbind:                  "mbind",
		DSYS_membarrier:             "membarrier",
		DSYS_memfd_create:           "memfd_create",
		DSYS_migrate_pages:          "migrate_pages",
		DSYS_mincore:                "mincore",
		DSYS_mkdir:                  "mkdir",
		DSYS_mkdirat:                "mkdirat",
		DSYS_mknod:                  "mknod",
		DSYS_mknodat:                "mknodat",
		DSYS_mlock:                  "mlock",
		DSYS_mlock2:                 "mlock2",
		DSYS_mlockall:               "mlockall",
		DSYS_mmap:                   "mmap",
		DSYS_modify_ldt:             "modify_ldt",
		DSYS_mount:                  "mount",
		DSYS_move_pages:             "move_pages",
		DSYS_mprotect:               "mprotect",
		DSYS_mq_getsetattr:          "mq_getsetattr",
		DSYS_mq_notify:              "mq_notify",
		DSYS_mq_open:                "mq_open",
		DSYS_mq_timedreceive:        "mq_timedreceive",
		DSYS_mq_timedsend:           "mq_timedsend",
		DSYS_mq_unlink:              "mq_unlink",
		DSYS_mremap:                 "mremap",
		DSYS_msgctl:                 "msgctl",
		DSYS_msgget:                 "msgget",
		DSYS_msgrcv:                 "msgrcv",
		DSYS_msgsnd:                 "msgsnd",
		DSYS_msync:                  "msync",
		DSYS_munlock:                "munlock",
		DSYS_munlockall:             "munlockall",
		DSYS_munmap:                 "munmap",
		DSYS_name_to_handle_at:      "name_to_handle_at",
		DSYS_nanosleep:              "nanosleep",
		DSYS_nfsservctl:             "nfsservctl",
		DSYS_open:                   "open",
		DSYS_open_by_handle_at:      "open_by_handle_at",
		DSYS_openat:                 "openat",
		DSYS_pause:                  "pause",
		DSYS_perf_event_open:        "perf_event_open",
		DSYS_personality:            "personality",
		DSYS_pipe:                   "pipe",
		DSYS_pipe2:                  "pipe2",
		DSYS_pivot_root:             "pivot_root",
		DSYS_pkey_alloc:             "pkey_alloc",
		DSYS_pkey_free:              "pkey_free",
		DSYS_pkey_mprotect:          "pkey_mprotect",
		DSYS_poll:                   "poll",
		DSYS_ppoll:                  "ppoll",
		DSYS_prctl:                  "prctl",
		DSYS_pread:                  "pread",
		DSYS_preadv:                 "preadv",
		DSYS_preadv2:                "preadv2",
		DSYS_prlimit64:              "prlimit64",
		DSYS_process_vm_readv:       "process_vm_readv",
		DSYS_process_vm_writev:      "process_vm_writev",
		DSYS_pselect6:               "pselect6",
		DSYS_ptrace:                 "ptrace",
		DSYS_putpmsg:                "putpmsg",
		DSYS_pwrite:                 "pwrite",
		DSYS_pwritev:                "pwritev",
		DSYS_pwritev2:               "pwritev2",
		DSYS_query_module:           "query_module",
		DSYS_quotactl:               "quotactl",
		DSYS_read:                   "read",
		DSYS_readahead:              "readahead",
		DSYS_readlink:               "readlink",
		DSYS_readlinkat:             "readlinkat",
		DSYS_readv:                  "readv",
		DSYS_reboot:                 "reboot",
		DSYS_recvfrom:               "recvfrom",
		DSYS_recvmmsg:               "recvmmsg",
		DSYS_recvmsg:                "recvmsg",
		DSYS_remap_file_pages:       "remap_file_pages",
		DSYS_removexattr:            "removexattr",
		DSYS_rename:                 "rename",
		DSYS_renameat:               "renameat",
		DSYS_renameat2:              "renameat2",
		DSYS_request_key:            "request_key",
		DSYS_restart_syscall:        "restart_syscall",
		DSYS_rmdir:                  "rmdir",
		DSYS_rt_sigaction:           "rt_sigaction",
		DSYS_rt_sigpending:          "rt_sigpending",
		DSYS_rt_sigprocmask:         "rt_sigprocmask",
		DSYS_rt_sigqueueinfo:        "rt_sigqueueinfo",
		DSYS_rt_sigreturn:           "rt_sigreturn",
		DSYS_rt_sigsuspend:          "rt_sigsuspend",
		DSYS_rt_sigtimedwait:        "rt_sigtimedwait",
		DSYS_rt_tgsigqueueinfo:      "rt_tgsigqueueinfo",
		DSYS_sched_get_priority_max: "sched_get_priority_max",
		DSYS_sched_get_priority_min: "sched_get_priority_min",
		DSYS_sched_getaffinity:      "sched_getaffinity",
		DSYS_sched_getattr:          "sched_getattr",
		DSYS_sched_getparam:         "sched_getparam",
		DSYS_sched_getscheduler:     "sched_getscheduler",
		DSYS_sched_rr_get_interval:  "sched_rr_get_interval",
		DSYS_sched_setaffinity:      "sched_setaffinity",
		DSYS_sched_setattr:          "sched_setattr",
		DSYS_sched_setparam:         "sched_setparam",
		DSYS_sched_setscheduler:     "sched_setscheduler",
		DSYS_sched_yield:            "sched_yield",
		DSYS_seccomp:                "seccomp",
		DSYS_security:               "security",
		DSYS_select:                 "select",
		DSYS_semctl:                 "semctl",
		DSYS_semget:                 "semget",
		DSYS_semop:                  "semop",
		DSYS_semtimedop:             "semtimedop",
		DSYS_sendfile:               "sendfile",
		DSYS_sendmmsg:               "sendmmsg",
		DSYS_sendmsg:                "sendmsg",
		DSYS_sendto:                 "sendto",
		DSYS_set_mempolicy:          "set_mempolicy",
		DSYS_set_robust_list:        "set_robust_list",
		DSYS_set_thread_area:        "set_thread_area",
		DSYS_set_tid_address:        "set_tid_address",
		DSYS_setdomainname:          "setdomainname",
		DSYS_setfsgid:               "setfsgid",
		DSYS_setfsuid:               "setfsuid",
		DSYS_setgid:                 "setgid",
		DSYS_setgroups:              "setgroups",
		DSYS_sethostname:            "sethostname",
		DSYS_setitimer:              "setitimer",
		DSYS_setns:                  "setns",
		DSYS_setpgid:                "setpgid",
		DSYS_setpriority:            "setpriority",
		DSYS_setregid:               "setregid",
		DSYS_setresgid:              "setresgid",
		DSYS_setresuid:              "setresuid",
		DSYS_setreuid:               "setreuid",
		DSYS_setrlimit:              "setrlimit",
		DSYS_setsid:                 "setsid",
		DSYS_setsockopt:             "setsockopt",
		DSYS_settimeofday:           "settimeofday",
		DSYS_setuid:                 "setuid",
		DSYS_setxattr:               "setxattr",
		DSYS_shmat:                  "shmat",
		DSYS_shmctl:                 "shmctl",
		DSYS_shmdt:                  "shmdt",
		DSYS_shmget:                 "shmget",
		DSYS_shutdown:               "shutdown",
		DSYS_sigaltstack:            "sigaltstack",
		DSYS_signalfd:               "signalfd",
		DSYS_signalfd4:              "signalfd4",
		DSYS_socket:                 "socket",
		DSYS_socketpair:             "socketpair",
		DSYS_splice:                 "splice",
		DSYS_stat:                   "stat",
		DSYS_statfs:                 "statfs",
		DSYS_statx:                  "statx",
		DSYS_swapoff:                "swapoff",
		DSYS_swapon:                 "swapon",
		DSYS_symlink:                "symlink",
		DSYS_symlinkat:              "symlinkat",
		DSYS_sync:                   "sync",
		DSYS_sync_file_range:        "sync_file_range",
		DSYS_syncfs:                 "syncfs",
		DSYS_sysfs:                  "sysfs",
		DSYS_sysinfo:                "sysinfo",
		DSYS_syslog:                 "syslog",
		DSYS_tee:                    "tee",
		DSYS_tgkill:                 "tgkill",
		DSYS_time:                   "time",
		DSYS_timer_create:           "timer_create",
		DSYS_timer_delete:           "timer_delete",
		DSYS_timer_getoverrun:       "timer_getoverrun",
		DSYS_timer_gettime:          "timer_gettime",
		DSYS_timer_settime:          "timer_settime",
		DSYS_timerfd_create:         "timerfd_create",
		DSYS_timerfd_gettime:        "timerfd_gettime",
		DSYS_timerfd_settime:        "timerfd_settime",
		DSYS_times:                  "times",
		DSYS_tkill:                  "tkill",
		DSYS_truncate:               "truncate",
		DSYS_tuxcall:                "tuxcall",
		DSYS_umask:                  "umask",
		DSYS_umount2:                "umount2",
		DSYS_uname:                  "uname",
		DSYS_unlink:                 "unlink",
		DSYS_unlinkat:               "unlinkat",
		DSYS_unshare:                "unshare",
		DSYS_uselib:                 "uselib",
		DSYS_userfaultfd:            "userfaultfd",
		DSYS_ustat:                  "ustat",
		DSYS_utime:                  "utime",
		DSYS_utimensat:              "utimensat",
		DSYS_utimes:                 "utimes",
		DSYS_vfork:                  "vfork",
		DSYS_vhangup:                "vhangup",
		DSYS_vmsplice:               "vmsplice",
		DSYS_vserver:                "vserver",
		DSYS_wait4:                  "wait4",
		DSYS_waitid:                 "waitid",
		DSYS_write:                  "write",
		DSYS_writev:                 "writev",
	}

	fcntls = map[int]string{
		DF_ADD_SEALS: "F_ADD_SEALS",
		// DF_APP:              "F_APP",
		// DF_CANCELLK:         "F_CANCELLK",
		DF_DUPFD:         "F_DUPFD",
		DF_DUPFD_CLOEXEC: "F_DUPFD_CLOEXEC",
		// DF_EOF:              "F_EOF",
		// DF_ERR:              "F_ERR",
		DF_GETFD:    "F_GETFD",
		DF_GETFL:    "F_GETFL",
		DF_GETLEASE: "F_GETLEASE",
		DF_GETLK:    "F_GETLK",
		// DF_GETLK64:          "F_GETLK64",
		DF_GETOWN: "F_GETOWN",
		// DF_GETOWNER_UIDS:    "F_GETOWNER_UIDS",
		DF_GETOWN_EX:  "F_GETOWN_EX",
		DF_GETPIPE_SZ: "F_GETPIPE_SZ",
		DF_GETSIG:     "F_GETSIG",
		// DF_GET_FILE_RW_HINT: "F_GET_FILE_RW_HINT",
		// DF_GET_RW_HINT:      "F_GET_RW_HINT",
		DF_GET_SEALS: "F_GET_SEALS",
		// DF_LOCK:             "F_LOCK",
		// DF_NORD:             "F_NORD",
		DF_NOTIFY: "F_NOTIFY",
		// DF_NOWR:             "F_NOWR",
		DF_OFD_GETLK:  "F_OFD_GETLK",
		DF_OFD_SETLK:  "F_OFD_SETLK",
		DF_OFD_SETLKW: "F_OFD_SETLKW",
		// DF_OK:               "F_OK",
		// DF_OWNER_GID:        "F_OWNER_GID",
		// DF_OWNER_PGRP:       "F_OWNER_PGRP",
		// DF_OWNER_PID:        "F_OWNER_PID",
		// DF_OWNER_TID:        "F_OWNER_TID",
		// DF_PERM:             "F_PERM",
		// DF_RDLCK:            "F_RDLCK",
		// DF_SEAL_GROW:        "F_SEAL_GROW",
		// DF_SEAL_SEAL:        "F_SEAL_SEAL",
		// DF_SEAL_SHRINK:      "F_SEAL_SHRINK",
		// DF_SEAL_WRITE:       "F_SEAL_WRITE",
		DF_SETFD:    "F_SETFD",
		DF_SETFL:    "F_SETFL",
		DF_SETLEASE: "F_SETLEASE",
		DF_SETLK:    "F_SETLK",
		// DF_SETLK64:          "F_SETLK64",
		DF_SETLKW: "F_SETLKW",
		// DF_SETLKW64:         "F_SETLKW64",
		DF_SETOWN:     "F_SETOWN",
		DF_SETOWN_EX:  "F_SETOWN_EX",
		DF_SETPIPE_SZ: "F_SETPIPE_SZ",
		DF_SETSIG:     "F_SETSIG",
		// DF_SET_FILE_RW_HINT: "F_SET_FILE_RW_HINT",
		// DF_SET_RW_HINT:      "F_SET_RW_HINT",
		// DF_SVB:              "F_SVB",
		// DF_TEST:             "F_TEST",
		// DF_TLOCK:            "F_TLOCK",
		// DF_ULOCK:            "F_ULOCK",
		// DF_UNLCK:            "F_UNLCK",
		// DF_WRLCK:            "F_WRLCK",
	}
)

func log(s string, a ...interface{}) {
	fmt.Printf(s, a...)
	if n := len(s); n != 0 && s[n-1] != '\n' {
		fmt.Println()
	}
	os.Stderr.Sync()
}

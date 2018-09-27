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
	Xfflush(tls, 0)
	fmt.Printf("==== %s\n%s----\n", strings.Join(a, "\n"), debug.Stack())
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

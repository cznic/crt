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
	env       = os.Environ()
	files     = map[uintptr]*os.File{}
	filesMu   sync.Mutex
	logging   bool

	Log = func(s string, a ...interface{}) {}

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
	if tls := MainTLS(); (*s1__pthread)(unsafe.Pointer(tls)).Fself != uintptr(tls) { // sanity check
		panic("internal error")
	}

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
	if fn := os.Getenv("CCGOLOG"); fn != "" {
		logging = true
		f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			panic(err)
		}

		pid := fmt.Sprintf("[pid %v] ", os.Getpid())

		Log = func(s string, args ...interface{}) {
			if s == "" {
				s = strings.Repeat("%v ", len(args))
			}
			_, fn, fl, _ := runtime.Caller(1)
			s = fmt.Sprintf(pid+"%s:%d: "+s, append([]interface{}{filepath.Base(fn), fl}, args...)...)
			switch {
			case len(s) != 0 && s[len(s)-1] == '\n':
				fmt.Fprint(f, s)
			default:
				fmt.Fprintln(f, s)
			}
		}

		Log("==== start: %v", os.Args)
	}
	tls := MainTLS()
	Xexit(tls, main(tls, int32(len(os.Args)), *(*uintptr)(unsafe.Pointer(X__ccgo_argv))))
}

func MainTLS() TLS {
	return TLS(*(*uintptr)(unsafe.Pointer(X__ccgo_main_tls)))
}

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
	allocMu.Lock()
	p, err := allocator.UintptrMalloc(size)
	allocMu.Unlock()
	return p, err
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
func Calloc(size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrCalloc(size)
	allocMu.Unlock()
	return p, err
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
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
		panic("internal error")
	}

	return r
}

// TS allocates the R/O text segment of a package/command.
func TS(init string) uintptr { return (*reflect.StringHeader)(unsafe.Pointer(&init)).Data }

// Free frees memory allocated by Calloc, Malloc or Realloc.
func Free(p uintptr) error {
	allocMu.Lock()
	err := allocator.UintptrFree(p)
	allocMu.Unlock()
	return err
}

// DS allocates the the data segment of a package/command.
func DS(init []byte) uintptr {
	r := (*reflect.SliceHeader)(unsafe.Pointer(&init)).Data
	if r%unsafe.Sizeof(uintptr(0)) != 0 {
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

// GoStringLen returns a string from a C char* string s having length len bytes.
func GoStringLen(s uintptr, len int) string {
	if len == 0 {
		return ""
	}

	return string((*rawmem)(unsafe.Pointer(s))[:len])
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

// Realloc reallocates memory.
func Realloc(p uintptr, size int) (uintptr, error) {
	allocMu.Lock()
	p, err := allocator.UintptrRealloc(p, size)
	allocMu.Unlock()
	return p, err
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
	Log("%s\n%s----\n", strings.Join(a, "\n"), debug.Stack())
}

func TraceOn()  { trace++ }
func TraceOff() { trace-- }

func Watch(tls TLS) { //TODO-
	if trace > 0 {
		_, fn, fl, _ := runtime.Caller(1)
		Log("# trace %s:%d:\n", filepath.Base(fn), fl)
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

func File(fd uintptr) *os.File {
	if int(fd) < 0 {
		return nil
	}

	filesMu.Lock()

	defer filesMu.Unlock()

	f := files[fd]
	if f == nil {
		f = os.NewFile(fd, fmt.Sprintf("fd%v", fd))
		files[fd] = f
	}
	return f
}

func X__log(tls TLS, format uintptr, args ...interface{}) {
	const sz = 1 << 10
	if !logging {
		return
	}
	buf := Xmalloc(tls, sz)
	defer Xfree(tls, buf)
	ap := X__builtin_va_start(tls, args)
	n := Xvsnprintf(tls, buf, sz, format, ap)
	Log("%s", GoStringLen(buf, int(n)))
	X__builtin_va_end(tls, ap)
}

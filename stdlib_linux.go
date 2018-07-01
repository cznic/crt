// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file uses also code originally from
//
//	https://github.com/golang/go/tree/97677273532cf1a4e8b181c242d89c0be8c92bb6/src/io/ioutil
//
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO-LICENSE file.

package crt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/cznic/crt/errno"
	"github.com/cznic/crt/limits"
	"github.com/cznic/mathutil"
)

var (
	// Random number state.
	// We generate random temporary file names so that there's a good
	// chance the file doesn't exist yet - keeps the number of tries in
	// tempFile to a minimum.
	rand   uint32
	randmu sync.Mutex

	openFDsMu sync.Mutex
	openFDs   = map[int32]struct{}{}
)

// int atoi(const char *str);
//
// The call atoi(str) shall be equivalent to:
//
// (int) strtol(str, (char **)NULL, 10)
//
// except that the handling of errors may differ. If the value cannot be
// represented, the behavior is undefined.
//
// The atoi() function shall return the converted value if the value can be
// represented.
func Xatoi(tls TLS, str uintptr) int32 {
	e0 := int(tls.err())
	tls.setErrno(0)
	n := Xstrtol(tls, str, 0, 10)
	if tls.err() != 0 {
		switch {
		case n < 0:
			return limits.XINT_MIN
		default:
			return limits.XINT_MAX
		}
	}

	tls.setErrno(e0)
	return int32(n)
}

// long strtol(const char *restrict nptr, char **restrict endptr, int base);
//
// These functions shall convert the initial portion of the string pointed to
// by nptr to a type long and long long representation, respectively. First,
// they decompose the input string into three parts:
//
// An initial, possibly empty, sequence of white-space characters (as specified
// by isspace())
//
// A subject sequence interpreted as an integer represented in some radix
// determined by the value of base
//
// A final string of one or more unrecognized characters, including the
// terminating NUL character of the input string.
//
// Then they shall attempt to convert the subject sequence to an integer, and
// return the result.
//
// If the value of base is 0, the expected form of the subject sequence is that
// of a decimal constant, octal constant, or hexadecimal constant, any of which
// may be preceded by a '+' or '-' sign. A decimal constant begins with a
// non-zero digit, and consists of a sequence of decimal digits. An octal
// constant consists of the prefix '0' optionally followed by a sequence of the
// digits '0' to '7' only. A hexadecimal constant consists of the prefix 0x or
// 0X followed by a sequence of the decimal digits and letters 'a' (or 'A' ) to
// 'f' (or 'F' ) with values 10 to 15 respectively.
//
// If the value of base is between 2 and 36, the expected form of the subject
// sequence is a sequence of letters and digits representing an integer with
// the radix specified by base, optionally preceded by a '+' or '-' sign. The
// letters from 'a' (or 'A' ) to 'z' (or 'Z' ) inclusive are ascribed the
// values 10 to 35; only letters whose ascribed values are less than that of
// base are permitted. If the value of base is 16, the characters 0x or 0X may
// optionally precede the sequence of letters and digits, following the sign if
// present.
//
// The subject sequence is defined as the longest initial subsequence of the
// input string, starting with the first non-white-space character that is of
// the expected form. The subject sequence shall contain no characters if the
// input string is empty or consists entirely of white-space characters, or if
// the first non-white-space character is other than a sign or a permissible
// letter or digit.
//
// If the subject sequence has the expected form and the value of base is 0,
// the sequence of characters starting with the first digit shall be
// interpreted as an integer constant. If the subject sequence has the expected
// form and the value of base is between 2 and 36, it shall be used as the base
// for conversion, ascribing to each letter its value as given above. If the
// subject sequence begins with a <hyphen-minus>, the value resulting from the
// conversion shall be negated. A pointer to the final string shall be stored
// in the object pointed to by endptr, provided that endptr is not a null
// pointer.
//
// In other than the C [CX] ￼  or POSIX ￼ locale, additional locale-specific
// subject sequence forms may be accepted.
//
// If the subject sequence is empty or does not have the expected form, no
// conversion is performed; the value of nptr shall be stored in the object
// pointed to by endptr, provided that endptr is not a null pointer.
//
// These functions shall not change the setting of errno if successful.
//
// Since 0, {LONG_MIN} or {LLONG_MIN}, and {LONG_MAX} or {LLONG_MAX} are
// returned on error and are also valid returns on success, an application
// wishing to check for error situations should set errno to 0, then call
// strtol() or strtoll(), then check errno.
//
// RETURN VALUE Upon successful completion, these functions shall return the
// converted value, if any. If no conversion could be performed, 0 shall be
// returned [CX] ￼  and errno may be set to [EINVAL]. ￼
//
// [CX] ￼ If the value of base is not supported, 0 shall be returned and errno
// shall be set to [EINVAL]. ￼
//
// If the correct value is outside the range of representable values,
// {LONG_MIN}, {LONG_MAX}, {LLONG_MIN}, or {LLONG_MAX} shall be returned
// (according to the sign of the value), and errno set to [ERANGE].
func Xstrtol(tls TLS, nptr, endptr uintptr, base int32) (r long_t) {
	switch base {
	case 10:
		ovf := false
		p := nptr
		for {
			ch := *(*byte)(unsafe.Pointer(p))
			switch c := int32(ch); {
			case Xisspace(tls, c) != 0:
				p++
			case c >= '0' && c <= '9':
				for {
					p++
					r0 := r
					r := 10*r + long_t(c) - '0'
					if r < r0 {
						ovf = true
					}
					switch c = int32(*(*byte)(unsafe.Pointer(p))); {
					case c >= '0' && c <= '9':
						// ok
					default:
						if endptr != 0 {
							*(*uintptr)(unsafe.Pointer(endptr)) = p
						}
						if ovf {
							tls.setErrno(errno.XERANGE)
							r = limits.XLONG_MAX
						}
						return r
					}
				}
			default:
				if endptr != 0 {
					*(*uintptr)(unsafe.Pointer(endptr)) = nptr
					tls.setErrno(errno.XEINVAL)
					return 0
				}
			}
		}
	case 16:
		ovf := false
		p := nptr
		var num long_t
		for {
			ch := *(*byte)(unsafe.Pointer(p))
			switch c := int32(ch); {
			case Xisspace(tls, c) != 0:
				p++
			case c == '+':
				panic("TODO")
			case c == '-':
				panic("TODO")
			case isHex(c, &num):
				for {
					p++
					r0 := r
					r := 10*r + num
					if r < r0 {
						ovf = true
					}
					switch c = int32(*(*byte)(unsafe.Pointer(p))); {
					case c == '+':
						panic("TODO")
					case c == '-':
						panic("TODO")
					case isHex(c, &num):
						// ok
					default:
						if endptr != 0 {
							*(*uintptr)(unsafe.Pointer(endptr)) = p
						}
						if ovf {
							tls.setErrno(errno.XERANGE)
							r = limits.XLONG_MAX
						}
						return r
					}
				}
			default:
				if endptr != 0 {
					*(*uintptr)(unsafe.Pointer(endptr)) = nptr
					tls.setErrno(errno.XEINVAL)
					return 0
				}
			}
		}
	default:
		panic(fmt.Errorf("%q %#x %v", GoString(nptr), endptr, base))
	}
}

func isHex(c int32, out *long_t) bool {
	if c >= '0' && c <= '0' {
		*out = long_t(c - '0')
		return true
	}

	if c >= 'a' && c <= 'f' {
		*out = long_t(c - 'a' + 10)
		return true
	}

	if c >= 'A' && c <= 'F' {
		*out = long_t(c - 'A' + 10)
		return true
	}

	return false
}

// unsigned long int strtoul(const char *str, char **endptr, int base);
//
// These functions shall convert the initial portion of the string pointed to
// by str to a type unsigned long and unsigned long long representation,
// respectively. First, they decompose the input string into three parts:
//
// An initial, possibly empty, sequence of white-space characters (as specified
// by isspace())
//
// A subject sequence interpreted as an integer represented in some radix
// determined by the value of base
//
// A final string of one or more unrecognized characters, including the
// terminating NUL character of the input string
//
// Then they shall attempt to convert the subject sequence to an unsigned
// integer, and return the result.
//
// If the value of base is 0, the expected form of the subject sequence is that
// of a decimal constant, octal constant, or hexadecimal constant, any of which
// may be preceded by a '+' or '-' sign. A decimal constant begins with a
// non-zero digit, and consists of a sequence of decimal digits. An octal
// constant consists of the prefix '0' optionally followed by a sequence of the
// digits '0' to '7' only. A hexadecimal constant consists of the prefix 0x or
// 0X followed by a sequence of the decimal digits and letters 'a' (or 'A' ) to
// 'f' (or 'F' ) with values 10 to 15 respectively.
//
// If the value of base is between 2 and 36, the expected form of the subject
// sequence is a sequence of letters and digits representing an integer with
// the radix specified by base, optionally preceded by a '+' or '-' sign. The
// letters from 'a' (or 'A' ) to 'z' (or 'Z' ) inclusive are ascribed the
// values 10 to 35; only letters whose ascribed values are less than that of
// base are permitted. If the value of base is 16, the characters 0x or 0X may
// optionally precede the sequence of letters and digits, following the sign if
// present.
//
// The subject sequence is defined as the longest initial subsequence of the
// input string, starting with the first non-white-space character that is of
// the expected form. The subject sequence shall contain no characters if the
// input string is empty or consists entirely of white-space characters, or if
// the first non-white-space character is other than a sign or a permissible
// letter or digit.
//
// If the subject sequence has the expected form and the value of base is 0,
// the sequence of characters starting with the first digit shall be
// interpreted as an integer constant. If the subject sequence has the expected
// form and the value of base is between 2 and 36, it shall be used as the base
// for conversion, ascribing to each letter its value as given above. If the
// subject sequence begins with a <hyphen-minus>, the value resulting from the
// conversion shall be negated. A pointer to the final string shall be stored
// in the object pointed to by endptr, provided that endptr is not a null
// pointer.
//
// In other than the C [CX] [Option Start]  or POSIX [Option End] locale,
// additional locale-specific subject sequence forms may be accepted.
//
// If the subject sequence is empty or does not have the expected form, no
// conversion shall be performed; the value of str shall be stored in the
// object pointed to by endptr, provided that endptr is not a null pointer.
//
// These functions shall not change the setting of errno if successful.
//
// Since 0, {ULONG_MAX}, and {ULLONG_MAX} are returned on error and are also
// valid returns on success, an application wishing to check for error
// situations should set errno to 0, then call strtoul() or strtoull(), then
// check errno.
//
// Upon successful completion, these functions shall return the converted
// value, if any. If no conversion could be performed, 0 shall be returned [CX]
// [Option Start]  and errno may be set to [EINVAL]. [Option End]
//
// If the correct value is outside the range of representable values,
// {ULONG_MAX} or {ULLONG_MAX} shall be returned and errno set to [ERANGE].
func Xstrtoul(tls TLS, str, endptr uintptr, base int32) (r ulong_t) {
	switch base {
	case 10:
		ovf := false
		p := str
		for {
			ch := *(*byte)(unsafe.Pointer(p))
			switch c := int32(ch); {
			case Xisspace(tls, c) != 0:
				p++
			case c >= '0' && c <= '9':
				for {
					p++
					r0 := r
					r := 10*r + ulong_t(c) - '0'
					if r < r0 {
						ovf = true
					}
					switch c = int32(*(*byte)(unsafe.Pointer(p))); {
					case c >= '0' && c <= '9':
						// ok
					default:
						if endptr != 0 {
							*(*uintptr)(unsafe.Pointer(endptr)) = p
						}
						if ovf {
							tls.setErrno(errno.XERANGE)
							r = limits.XULONG_MAX
						}
						return r
					}
				}
			default:
				if endptr != 0 {
					*(*uintptr)(unsafe.Pointer(endptr)) = str
					tls.setErrno(errno.XEINVAL)
					return 0
				}
			}
		}
	default:
		panic(fmt.Errorf("%q %#x %v", GoString(str), endptr, base))
	}
}

// void exit(int);
func Xexit(tls TLS, n int32) { X__builtin_exit(tls, n) }

// // void exit(int);
// func X__builtin_exit(tls TLS, n int32) {
// 	os.Exit(int(n))
// }

// void free(void *ptr);
func Xfree(tls TLS, ptr uintptr) {
	free(tls, ptr)
	// if strace {
	// 	fmt.Fprintf(os.Stderr, "free(%#x)\n", ptr)
	// }
}

// void abort();
func Xabort(tls TLS) { X__builtin_abort(tls) }

// void __builtin_trap();
func X__builtin_trap(tls TLS) { os.Exit(1) }

// void abort();
func X__builtin_abort(tls TLS) { X__builtin_trap(tls) }

// char *getenv(const char *name);
func Xgetenv(tls TLS, name uintptr) uintptr {
	// The getenv() function shall search the environment of the calling
	// process (see XBD Environment Variables) for the environment variable
	// name if it exists and return a pointer to the value of the
	// environment variable. If the specified environment variable cannot
	// be found, a null pointer shall be returned. The application shall
	// ensure that it does not modify the string pointed to by the getenv()
	// function.

	// Upon successful completion, getenv() shall return a pointer to a
	// string containing the value for the specified name. If the specified
	// name cannot be found in the environment of the calling process, a
	// null pointer shall be returned.
	env := *(*uintptr)(unsafe.Pointer(penviron))
	if env == 0 {
		return 0
	}

	n := Xstrlen(tls, name)
	for {
		s := *(*uintptr)(unsafe.Pointer(env))
		if s == 0 {
			return 0
		}

		if Xstrncmp(tls, s, name, n) == 0 {
			if *(*byte)(unsafe.Pointer(s + uintptr(n))) == '=' {
				return s + uintptr(n) + 1
			}
		}

		env += unsafe.Sizeof(uintptr(0))
	}
}

// int abs(int j);
func X__builtin_abs(tls TLS, j int32) int32 {
	if j < 0 {
		return -j
	}

	return j
}

// int abs(int j);
func Xabs(tls TLS, j int32) int32 { return X__builtin_abs(tls, j) }

type sorter struct {
	base   uintptr
	compar func(tls TLS, a, b uintptr) int32
	nmemb  int
	size   uintptr
	tls    TLS
	buf    []byte
}

func (s *sorter) Len() int { return s.nmemb }

func (s *sorter) Less(i, j int) bool {
	return s.compar(s.tls, s.base+uintptr(i)*s.size, s.base+uintptr(j)*s.size) < 0
}

func (s *sorter) Swap(i, j int) {
	p := s.base + uintptr(i)*s.size
	q := s.base + uintptr(j)*s.size
	switch s.size {
	case 1:
		*(*int8)(unsafe.Pointer(p)), *(*int8)(unsafe.Pointer(q)) = *(*int8)(unsafe.Pointer(q)), *(*int8)(unsafe.Pointer(p))
	case 2:
		*(*int16)(unsafe.Pointer(p)), *(*int16)(unsafe.Pointer(q)) = *(*int16)(unsafe.Pointer(q)), *(*int16)(unsafe.Pointer(p))
	case 4:
		*(*int32)(unsafe.Pointer(p)), *(*int32)(unsafe.Pointer(q)) = *(*int32)(unsafe.Pointer(q)), *(*int32)(unsafe.Pointer(p))
	case 8:
		*(*int64)(unsafe.Pointer(p)), *(*int64)(unsafe.Pointer(q)) = *(*int64)(unsafe.Pointer(q)), *(*int64)(unsafe.Pointer(p))
	default:
		size := int(s.size)
		if s.buf == nil {
			s.buf = make([]byte, size)
		}
		Copy(uintptr(unsafe.Pointer(&s.buf[0])), p, size)
		Copy(p, q, size)
		Copy(q, uintptr(unsafe.Pointer(&s.buf[0])), size)
	}
}

// void qsort(void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));
func qsort(tls TLS, base uintptr, nmemb, size size_t, compar uintptr) {
	if size > mathutil.MaxInt {
		panic("size overflow")
	}

	if nmemb > mathutil.MaxInt {
		panic("nmemb overflow")
	}

	s := sorter{base, *(*func(TLS, uintptr, uintptr) int32)(unsafe.Pointer(&compar)), int(nmemb), uintptr(size), tls, nil}
	sort.Sort(&s)
}

// int system(const char *command);
func Xsystem(tls TLS, command uintptr) int32 {
	if command == 0 || *(*int8)(unsafe.Pointer(command)) == 0 {
		return 1
	}

	cmd := exec.Command("sh", "-c", GoString(command))
	if err := cmd.Run(); err != nil {
		return int32(cmd.ProcessState.Sys().(syscall.WaitStatus))
	}

	return 0
}

// void *calloc(size_t nmemb, size_t size);
func Xcalloc(tls TLS, nmemb, size size_t) (p uintptr) {
	hi, lo := mathutil.MulUint128_64(uint64(nmemb), uint64(size))
	if hi == 0 && lo <= mathutil.MaxInt {
		p = calloc(tls, int(lo))
	}
	// if strace {
	// 	fmt.Fprintf(os.Stderr, "calloc(%#x) %#x\n", size, p)
	// }
	return p
}

// void *malloc(size_t size);
func X__builtin_malloc(tls TLS, size size_t) (p uintptr) {
	if size < mathutil.MaxInt {
		p = malloc(tls, int(size))
	}
	// if strace {
	// 	fmt.Fprintf(os.Stderr, "malloc(%#x) %#x\n", size, p)
	// }
	return p
}

// void *malloc(size_t size);
func Xmalloc(tls TLS, size size_t) uintptr { return X__builtin_malloc(tls, size) }

// void *realloc(void *ptr, size_t size);
func Xrealloc(tls TLS, ptr uintptr, size size_t) (p uintptr) {
	if size < mathutil.MaxInt {
		p = realloc(tls, ptr, int(size))
	}
	// if strace {
	// 	fmt.Fprintf(os.Stderr, "realloc(%#x, %#x) %#x\n", ptr, size, p)
	// }
	return p
}

// size_t malloc_usable_size (void *ptr);
// func Xmalloc_usable_size(tls TLS, ptr uintptr) size_t { return size_t(memory.UintptrUsableSize(ptr)) }

// void qsort(void *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));
func Xqsort(tls TLS, base uintptr, nmemb, size size_t, compar uintptr) {
	qsort(tls, base, nmemb, size, compar)
}

// int mkstemps(char *template, int suffixlen);
func Xmkstemps(tls TLS, template uintptr, suffixlen int32) int32 {
	panic("TODO")
}

// int mkstemp(char *template);
//
// The mkstemp() function shall create a regular file with a unique name
// derived from template and return a file descriptor for the file open for
// reading and writing. The application shall ensure that the string provided
// in template is a pathname ending with at least six trailing 'X' characters.
// The mkstemp() function shall modify the contents of template by replacing
// six or more 'X' characters at the end of the pathname with the same number
// of characters from the portable filename character set. The characters shall
// be chosen such that the resulting pathname does not duplicate the name of an
// existing file at the time of the call to mkstemp(). The mkstemp() function
// shall use the resulting pathname to create the file, and obtain a file
// descriptor for it, as if by a call to:
//
// open(pathname, O_RDWR|O_CREAT|O_EXCL, S_IRUSR|S_IWUSR)
//
// By behaving as if the O_EXCL flag for open() is set, the function prevents
// any possible race condition between testing whether the file exists and
// opening it for use.
//
// Upon successful completion, the mkstemp() function shall return an open file
// descriptor. Otherwise, it shall return -1 and shall set errno to indicate
// the error.
func Xmkstemp(tls TLS, template uintptr) int32 {
	s := GoString(template)
	n := 0
	for len(s) != 0 && s[len(s)-1-n] == 'X' {
		n++
	}
	s = s[:len(s)-n]
	dir, prefix := filepath.Split(s)
	f, err := tempFile(dir, prefix, n)
	if err != nil {
		tls.setErrno(err)
		return -1
	}

	fd := int32(f.Fd())
	openFDsMu.Lock()
	openFDs[fd] = struct{}{}
	openFDsMu.Unlock()
	return fd
}

// char *realpath(const char *path, char *resolved_path);
func Xrealpath(tls TLS, path, resolved_path uintptr) uintptr {
	// realpath()  expands  all  symbolic links and resolves references to
	// /./, /../ and extra '/' characters in the null-terminated string
	// named by path to produce a canonicalized absolute path       name.
	// The resulting pathname is stored as a null-terminated string, up to
	// a maximum of PATH_MAX bytes, in the buffer pointed to by
	// resolved_path.  The resulting path will have no symbolic link, /./
	// or /../ components.

	pth := GoString(path)
	resolved, err := filepath.EvalSymlinks(pth)
	// If there is no error, realpath() returns a pointer to the
	// resolved_path.

	// Otherwise, it returns NULL, the contents of the array resolved_path
	// are undefined, and errno is set to indicate the error.
	if err != nil {
		switch {
		case os.IsNotExist(err):
			tls.setErrno(errno.XENOENT)
		default:
			panic(fmt.Errorf("%q: %v", pth, err))
		}
		return 0
	}

	// If resolved_path is specified as NULL, then realpath() uses
	// malloc(3) to allocate a buffer of up to PATH_MAX bytes to hold the
	// resolved pathname, and returns a pointer to this buffer.  The caller
	// should deallocate this buffer using free(3).
	if resolved_path == 0 {
		resolved_path = CString(resolved)
		return 0
	}

	if len(resolved) >= limits.XPATH_MAX {
		panic("TODO")
	}

	b := make([]byte, len(resolved)+1)
	copy(b, resolved)
	Copy(resolved_path, uintptr(unsafe.Pointer(&b[0])), len(b))
	return resolved_path
}

// double atof(const char *nptr);
func Xatof(tls TLS, nptr uintptr) float64 {
	panic("TODO")
}

// tempFile creates a new temporary file in the directory dir with a name
// beginning with prefix appended by n characters, opens the file for reading
// and writing, and returns the resulting *os.File.  If dir is the empty
// string, TempFile uses the default directory for temporary files (see
// os.TempDir).  Multiple programs calling TempFile simultaneously will not
// choose the same file. The caller can use f.Name() to find the pathname of
// the file. It is the caller's responsibility to remove the file when no
// longer needed.
func tempFile(dir, prefix string, n int) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nconflict := 0
	for i := 0; i < 10000; i++ {
		name := filepath.Join(dir, prefix+nextSuffix(n))
		f, err = os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
		if os.IsExist(err) {
			if nconflict++; nconflict > 10 {
				randmu.Lock()
				rand = reseed()
				randmu.Unlock()
			}
			continue
		}
		break
	}
	return
}

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextSuffix(n int) string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*1664525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1 : n+1]
}

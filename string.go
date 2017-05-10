// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"unsafe"
)

// // void *memmove(void *dest, const void *src, size_t n);
// func (c *cpu) memmove() {
// 	sp, n := popLong(c.sp)
// 	sp, src := popPtr(sp)
// 	dest := readPtr(sp)
// 	movemem(dest, src, int(n))
// 	writePtr(c.rp, dest)
// }
//
// // void *mempcpy(void *dest, const void *src, size_t n);
// func (c *cpu) mempcpy() {
// 	sp, n := popLong(c.sp)
// 	sp, src := popPtr(sp)
// 	dest := readPtr(sp)
// 	movemem(dest, src, int(n))
// 	writePtr(c.rp, dest+uintptr(n))
// }

// char *strcat(char *dest, const char *src)
func Xstrcat(dest, src *int8) *int8 {
	ret := dest
	for *dest != 0 {
		*(*uintptr)(unsafe.Pointer(&dest))++
	}
	for {
		ch := *src
		*(*uintptr)(unsafe.Pointer(&src))++
		*dest = ch
		*(*uintptr)(unsafe.Pointer(&dest))++
		if ch == 0 {
			return ret
		}
	}
}

// char *index(const char *s, int c)
func Xindex(s *int8, ch int32) *int8 { return Xstrchr(s, ch) }

// char *strchr(const char *s, int c)
func Xstrchr(s *int8, ch int32) *int8 {
	for {
		ch2 := byte(*s)
		if ch2 == byte(ch) {
			return s
		}

		if ch2 == 0 {
			return nil
		}

		*(*uintptr)(unsafe.Pointer(&s))++
	}
}

// int strcmp(const char *s1, const char *s2)
func Xstrcmp(s1, s2 *int8) int32 {
	for {
		ch1 := byte(*s1)
		*(*uintptr)(unsafe.Pointer(&s1))++
		ch2 := byte(*s2)
		*(*uintptr)(unsafe.Pointer(&s2))++
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			return int32(ch1) - int32(ch2)
		}
	}
}

// char *strcpy(char *dest, const char *src)
func Xstrcpy(dest, src *int8) *int8 {
	r := dest
	for {
		ch := *src
		*(*uintptr)(unsafe.Pointer(&src))++
		*dest = ch
		*(*uintptr)(unsafe.Pointer(&dest))++
		if ch == 0 {
			return r
		}
	}
}

// // char *strncpy(char *dest, const char *src, size_t n)
// func (c *cpu) strncpy() {
// 	sp, n := popLong(c.sp)
// 	sp, src := popPtr(sp)
// 	dest := readPtr(sp)
// 	ret := dest
// 	var ch int8
// 	for ch = readI8(src); ch != 0 && n > 0; n-- {
// 		writeI8(dest, ch)
// 		dest++
// 		src++
// 		ch = readI8(src)
// 	}
// 	for ; n > 0; n-- {
// 		writeI8(dest, 0)
// 		dest++
// 	}
// 	writePtr(c.rp, ret)
// }

// char *rindex(const char *s, int c)
func Xrindex(s *int8, ch int32) *int8 { return Xstrrchr(s, ch) }

// char *strrchr(const char *s, int c)
func Xstrrchr(s *int8, ch int32) *int8 {
	var ret *int8
	for {
		ch2 := byte(*s)
		if ch2 == 0 {
			return ret
		}

		if ch2 == byte(ch) {
			ret = s
		}
		*(*uintptr)(unsafe.Pointer(&s))++
	}
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64

package crt

import (
	"unsafe"
)

// char *strncpy(char *dest, const char *src, size_t n)
func Xstrncpy(dest, src *int8, n uint64) *int8 {
	ret := dest
	for c := *src; c != 0 && n > 0; n-- {
		*dest = c
		*(*uintptr)(unsafe.Pointer(&dest))++
		*(*uintptr)(unsafe.Pointer(&src))++
		c = *src
	}
	for ; n > 0; n-- {
		*dest = 0
		*(*uintptr)(unsafe.Pointer(&dest))++
	}
	return ret
}

// size_t strlen(const char *s)
func Xstrlen(s *int8) uint64 {
	var n uint64
	for ; *s != 0; *(*uintptr)(unsafe.Pointer(&s))++ {
		n++
	}
	return n
}

// int strncmp(const char *s1, const char *s2, size_t n)
func Xstrncmp(s1, s2 *int8, n uint64) int32 {
	var ch1, ch2 byte
	for n != 0 {
		ch1 = byte(*s1)
		*(*uintptr)(unsafe.Pointer(&s1))++
		ch2 = byte(*s2)
		*(*uintptr)(unsafe.Pointer(&s2))++
		n--
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			break
		}
	}
	if n != 0 {
		return int32(ch1) - int32(ch2)
	}

	return 0
}

// void *memset(void *s, int c, size_t n)
func Xmemset(s unsafe.Pointer, c int32, n uint64) unsafe.Pointer { return X__builtin_memset(s, c, n) }

// void *memset(void *s, int c, size_t n)
func X__builtin_memset(s unsafe.Pointer, c int32, n uint64) unsafe.Pointer {
	for d := (*int8)(unsafe.Pointer(s)); n > 0; n-- {
		*d = int8(c)
		*(*uintptr)(unsafe.Pointer(&d))++
	}
	return s
}

// void *memcpy(void *dest, const void *src, size_t n)
func Xmemcpy(dest, src unsafe.Pointer, n uint64) unsafe.Pointer {
	movemem(dest, src, int(n))
	return dest
}

// int memcmp(const void *s1, const void *s2, size_t n)
func Xmemcmp(s1, s2 unsafe.Pointer, n uint64) int32 {
	var ch1, ch2 byte
	for n != 0 {
		ch1 = *(*byte)(unsafe.Pointer(s1))
		*(*uintptr)(unsafe.Pointer(&s1))++
		ch2 = *(*byte)(unsafe.Pointer(s2))
		*(*uintptr)(unsafe.Pointer(&s2))++
		if ch1 != ch2 {
			break
		}

		n--
	}
	if n != 0 {
		return int32(ch1) - int32(ch2)
	}

	return 0
}

// void *memmove(void *dest, const void *src, size_t n);
func Xmemmove(dest, src unsafe.Pointer, n uint64) unsafe.Pointer {
	movemem(dest, src, int(n))
	return dest
}

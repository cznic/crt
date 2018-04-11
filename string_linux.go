// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"unsafe"
)

// char *strcat(char *dest, const char *src)
func Xstrcat(tls TLS, dest, src uintptr) uintptr {
	ret := dest
	for *(*int8)(unsafe.Pointer(dest)) != 0 {
		dest++
	}
	for {
		c := *(*int8)(unsafe.Pointer(src))
		src++
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		if c == 0 {
			return ret
		}
	}
}

// char *index(const char *s, int c)
func Xindex(tls TLS, s uintptr, c int32) uintptr { return Xstrchr(tls, s, c) }

// char *strchr(const char *s, int c)
func Xstrchr(tls TLS, s uintptr, c int32) uintptr {
	// The strchr() function returns a pointer to the first occurrence of
	// the character c in the string s.

	// The  strchr()  and strrchr() functions return a pointer to the
	// matched character or NULL if the character is not found.  The
	// terminating null byte is considered part of the string, so that
	// if c is specified as '\0', these functions return a pointer to the
	// terminator.
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == byte(c) {
			return s
		}

		if ch2 == 0 {
			return 0
		}

		s++
	}
}

// char *strchrnul(const char *s, int c);
func Xstrchrnul(tls TLS, s uintptr, c int32) uintptr {
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == 0 || ch2 == byte(c) {
			return s
		}

		s++
	}
}

// int strcmp(const char *s1, const char *s2)
func X__builtin_strcmp(tls TLS, s1, s2 uintptr) int32 {
	for {
		ch1 := *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 := *(*byte)(unsafe.Pointer(s2))
		s2++
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			return int32(ch1) - int32(ch2)
		}
	}
}

// int strcmp(const char *s1, const char *s2)
func Xstrcmp(tls TLS, s1, s2 uintptr) int32 { return X__builtin_strcmp(tls, s1, s2) }

// char *strcpy(char *dest, const char *src)
func X__builtin_strcpy(tls TLS, dest, src uintptr) uintptr {
	r := dest
	for {
		c := *(*int8)(unsafe.Pointer(src))
		src++
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		if c == 0 {
			return r
		}
	}
}

// char *strcpy(char *dest, const char *src)
func Xstrcpy(tls TLS, dest, src uintptr) uintptr { return X__builtin_strcpy(tls, dest, src) }

// char *rindex(const char *s, int c)
func Xrindex(tls TLS, s uintptr, c int32) uintptr { return Xstrrchr(tls, s, c) }

// char *strrchr(const char *s, int c)
func Xstrrchr(tls TLS, s uintptr, c int32) uintptr {
	var ret uintptr
	for {
		ch2 := *(*byte)(unsafe.Pointer(s))
		if ch2 == 0 {
			return ret
		}

		if ch2 == byte(c) {
			ret = s
		}
		s++
	}
}

// char *strstr(const char *haystack, const char *needle);
func Xstrstr(tls TLS, haystack, needle uintptr) uintptr {
	// Public domain code from http://clc-wiki.net/wiki/C_standard_library:string.h:strstr#Implementation
	//
	//	#include <string.h> /* size_t memcmp() strlen() */
	//	char *strstr(const char *s1, const char *s2)
	//	{
	//	    size_t n = strlen(s2);
	//	    while(*s1)
	//	        if(!memcmp(s1++,s2,n))
	//	            return s1-1;
	//	    return 0;
	//	}
	n := Xstrlen(tls, needle)
	for *(*byte)(unsafe.Pointer(haystack)) != 0 {
		if X__builtin_memcmp(tls, haystack, needle, n) == 0 {
			return haystack
		}

		haystack++
	}
	return 0
}

// char *strncpy(char *dest, const char *src, size_t n)
func Xstrncpy(tls TLS, dest, src uintptr, n size_t) uintptr {
	ret := dest
	for c := *(*int8)(unsafe.Pointer(src)); c != 0 && n > 0; n-- {
		*(*int8)(unsafe.Pointer(dest)) = c
		dest++
		src++
		c = *(*int8)(unsafe.Pointer(src))
	}
	for ; n > 0; n-- {
		*(*int8)(unsafe.Pointer(dest)) = 0
		dest++
	}
	return ret
}

// size_t strlen(const char *s)
func X__builtin_strlen(tls TLS, s uintptr) size_t {
	var n size_t
	for ; *(*int8)(unsafe.Pointer(s)) != 0; s++ {
		n++
	}
	return n
}

// size_t strlen(const char *s)
func Xstrlen(tls TLS, s uintptr) size_t { return X__builtin_strlen(tls, s) }

// int strncmp(const char *s1, const char *s2, size_t n)
//
// The strncmp() function shall compare not more than n bytes (bytes that
// follow a NUL character are not compared) from the array pointed to by s1 to
// the array pointed to by s2.
//
// The sign of a non-zero return value is determined by the sign of the
// difference between the values of the first pair of bytes (both interpreted
// as type unsigned char) that differ in the strings being compared.
//
// Upon successful completion, strncmp() shall return an integer greater than,
// equal to, or less than 0, if the possibly null-terminated array pointed to
// by s1 is greater than, equal to, or less than the possibly null-terminated
// array pointed to by s2 respectively.
func Xstrncmp(tls TLS, s1, s2 uintptr, n size_t) int32 {
	// Public domain code from http://clc-wiki.net/wiki/C_standard_library:string.h:strncmp
	//
	//	#include <stddef.h>
	//	int strncmp(const char* s1, const char* s2, size_t n)
	//	{
	//	    while(n--)
	//	        if(*s1++!=*s2++)
	//	            return *(unsigned char*)(s1 - 1) - *(unsigned char*)(s2 - 1);
	//	    return 0;
	//	}
	for n != 0 {
		n--
		ch1 := *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 := *(*byte)(unsafe.Pointer(s2))
		s2++
		if ch1 != ch2 || ch1 == 0 || ch2 == 0 {
			return int32(ch1) - int32(ch2)
		}
	}
	return 0
}

// void *memset(void *s, int c, size_t n)
func Xmemset(tls TLS, s uintptr, c int32, n size_t) uintptr {
	return X__builtin_memset(tls, s, c, n)
}

// void *memset(void *s, int c, size_t n)
func X__builtin_memset(tls TLS, s uintptr, c int32, n size_t) uintptr {
	for d := s; n > 0; n-- {
		*(*int8)(unsafe.Pointer(d)) = int8(c)
		d++
	}
	return s
}

// void *memcpy(void *dest, const void *src, size_t n)
func X__builtin_memcpy(tls TLS, dest, src uintptr, n size_t) uintptr {
	Copy(dest, src, int(n))
	return dest
}

// void *memcpy(void *dest, const void *src, size_t n)
func Xmemcpy(tls TLS, dest, src uintptr, n size_t) uintptr {
	return X__builtin_memcpy(tls, dest, src, n)
}

// int memcmp(const void *s1, const void *s2, size_t n)
func X__builtin_memcmp(tls TLS, s1, s2 uintptr, n size_t) int32 {
	var ch1, ch2 byte
	for n != 0 {
		ch1 = *(*byte)(unsafe.Pointer(s1))
		s1++
		ch2 = *(*byte)(unsafe.Pointer(s2))
		s2++
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

// int memcmp(const void *s1, const void *s2, size_t n)
func Xmemcmp(tls TLS, s1, s2 uintptr, n size_t) int32 {
	return X__builtin_memcmp(tls, s1, s2, n)
}

// void *memmove(void *dest, const void *src, size_t n);
func Xmemmove(tls TLS, dest, src uintptr, n size_t) uintptr {
	Copy(dest, src, int(n))
	return dest
}

// void *mempcpy(void *dest, const void *src, size_t n);
func Xmempcpy(tls TLS, dest, src uintptr, n size_t) uintptr {
	return dest + uintptr(Copy(dest, src, int(n)))
}

// char *strpbrk(const char *s, const char *accept);
func Xstrpbrk(tls TLS, s, accept uintptr) uintptr {
	// The strpbrk() function locates the first occurrence in the string s
	// of any of the bytes in the string accept.

	// The strpbrk() function returns a pointer to the byte in s that
	// matches one of the bytes in accept, or NULL if no such byte is
	// found.

	// Public domain code from http://clc-wiki.net/wiki/strpbrk
	//
	//	#include <string.h> /* strchr */
	//	char *strpbrk(const char *s1, const char *s2)
	//	{
	//	    while(*s1)
	//	        if(strchr(s2, *s1++))
	//	            return (char*)--s1;
	//	    return 0;
	//	}
	for ; ; s++ {
		switch c := *(*byte)(unsafe.Pointer(s)); c {
		case 0:
			return 0
		default:
			if Xstrchr(tls, accept, int32(c)) != 0 {
				return s
			}
		}
	}
}

// void *memchr(const void *s, int c, size_t n);
func Xmemchr(tls TLS, s uintptr, c int32, n size_t) uintptr {
	// Public domain code from http://clc-wiki.net/wiki/C_standard_library:string.h:memchr
	//
	//	#include <stddef.h>
	//	void *memchr(const void *s, int c, size_t n)
	//	{
	//	    unsigned char *p = (unsigned char*)s;
	//	    while( n-- )
	//	        if( *p != (unsigned char)c )
	//	            p++;
	//	        else
	//	            return p;
	//	    return 0;
	//	}
	for n != 0 {
		n--
		if *(*byte)(unsafe.Pointer(s)) != byte(c) {
			s++
		} else {
			return s
		}
	}
	return 0
}

// char *strerror(int errnum);
func Xstrerror(tls TLS, errnum int32) uintptr {
	panic("TODO")
}

// char *strdup(const char *s);
func Xstrdup(tls TLS, s uintptr) uintptr {
	panic("TODO")
}

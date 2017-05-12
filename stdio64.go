// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64
// +build !windows

package crt

const (
	longBits = 64
)

type file *struct {
	X0  int32
	X1  *int8
	X2  *int8
	X3  *int8
	X4  *int8
	X5  *int8
	X6  *int8
	X7  *int8
	X8  *int8
	X9  *int8
	X10 *int8
	X11 *int8
	X12 uintptr
	X13 uintptr
	X14 int32
	X15 int32
	X16 int64
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 uintptr
	X21 int64
	X22 uintptr
	X23 uintptr
	X24 uintptr
	X25 uintptr
	X26 uint64
	X27 int32
	X28 [20]int8
}

func (r *argsReader) readLong() int64 {
	s := *r
	v := s[0].(int64)
	*r = s[1:]
	return v
}

func (r *argsReader) readULong() uint64 {
	s := *r
	v := s[0].(uint64)
	*r = s[1:]
	return v
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfwrite(ptr uintptr, size, nmemb uint64, stream file) uint64 {
	return fwrite(ptr, size, nmemb, stream)
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfread(ptr uintptr, size, nmemb uint64, stream file) uint64 {
	return fread(ptr, size, nmemb, stream)
}

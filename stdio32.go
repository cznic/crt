// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 arm arm64be armbe mips mipsle ppc ppc64le s390 s390x sparc

// +build !windows

package crt

import (
	"unsafe"
)

const (
	TFILE = "struct{int32,*int8,*int8,*int8,*int8,*int8,*int8,*int8,*int8,*int8,*int8,*int8,*struct{},*struct{},int32,int32,int32,uint16,int8,[1]int8,*struct{},int64,*struct{},*struct{},*struct{},*struct{},uint32,int32,[40]int8}"
)

type XFILE struct {
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
	X12 unsafe.Pointer
	X13 unsafe.Pointer
	X14 int32
	X15 int32
	X16 int32
	X17 uint16
	X18 int8
	X19 [1]int8
	X20 unsafe.Pointer
	X21 int64
	X22 unsafe.Pointer
	X23 unsafe.Pointer
	X24 unsafe.Pointer
	X25 unsafe.Pointer
	X26 uint32
	X27 int32
	X28 [40]int8
}

// size_t fwrite(const void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfwrite(tls *TLS, ptr unsafe.Pointer, size, nmemb uint32, stream *XFILE) uint32 {
	return uint32(fwrite(tls, ptr, uint64(size), uint64(nmemb), stream))
}

// size_t fread(void *ptr, size_t size, size_t nmemb, FILE *stream);
func Xfread(tls *TLS, ptr unsafe.Pointer, size, nmemb uint32, stream *XFILE) uint32 {
	return uint32(fread(tls, ptr, uint64(size), uint64(nmemb), stream))
}

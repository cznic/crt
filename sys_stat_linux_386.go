// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 arm arm64be armbe mips mipsle ppc ppc64le s390 s390x sparc

package crt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const Tstruct_stat64 = "struct{uint64,uint32,uint32,uint32,uint32,uint32,uint32,uint64,uint32,int64,int32,int64,int32,uint32,int32,uint32,int32,uint32,uint64}"

type Xstruct_stat64 struct {
	X0  uint64
	X1  uint32
	X2  uint32
	X3  uint32
	X4  uint32
	X5  uint32
	X6  uint32
	X7  uint64
	X8  uint32
	X9  int64
	X10 int32
	X11 int64
	X12 int32
	X13 uint32
	X14 int32
	X15 uint32
	X16 int32
	X17 uint32
	X18 uint64
}

// extern int stat64(char *__file, struct stat64 *__buf);
func Xstat64(tls *TLS, file *int8, buf *Xstruct_stat64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_STAT64, uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(buf)), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "stat(%q, %#x) %v %v\n", GoString(file), buf, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// int fstat64(int fildes, struct stat64 *buf);
func Xfstat64(tls *TLS, fildes int32, buf *Xstruct_stat64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_FSTAT64, uintptr(fildes), uintptr(unsafe.Pointer(buf)), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "fstat(%v, %#x) %v %v\n", fildes, buf, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// extern int lstat64(char *__file, struct stat64 *__buf);
func Xlstat64(tls *TLS, file *int8, buf *Xstruct_stat64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_LSTAT64, uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(buf)), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "lstat(%q, %#x) %v %v\n", GoString(file), buf, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

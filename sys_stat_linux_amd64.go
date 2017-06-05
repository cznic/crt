// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64

package crt

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const Tstruct_stat64 = "struct{uint64,uint64,uint64,uint32,uint32,uint32,int32,uint64,int64,int64,int64,int64,uint64,int64,uint64,int64,uint64,[3]int64}"

type Xstruct_stat64 struct {
	X0  uint64
	X1  uint64
	X2  uint64
	X3  uint32
	X4  uint32
	X5  uint32
	X6  int32
	X7  uint64
	X8  int64
	X9  int64
	X10 int64
	X11 int64
	X12 uint64
	X13 int64
	X14 uint64
	X15 int64
	X16 uint64
	X17 [3]int64
}

// extern int stat64(char *__file, struct stat64 *__buf);
func Xstat64(tls *TLS, file *int8, buf *Xstruct_stat64) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_STAT, uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(buf)), 0)
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
	r, _, err := syscall.Syscall(syscall.SYS_FSTAT, uintptr(fildes), uintptr(unsafe.Pointer(buf)), 0)
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
	r, _, err := syscall.Syscall(syscall.SYS_LSTAT, uintptr(unsafe.Pointer(file)), uintptr(unsafe.Pointer(buf)), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "lstat(%q, %#x) %v %v\n", GoString(file), buf, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64

package crt

import (
	"fmt"
	"os"
	"unsafe"
)

const (
	Tpthread_attr_t  = "union{[56]int8,int64}"
	Tpthread_mutex_t = "union{struct{int32,uint32,int32,uint32,int32,int16,int16,struct{*struct{},*struct{}}},[40]int8,int64}"
)

type Xpthread_mutex_t struct {
	X [0]struct {
		X0 struct {
			X0 int32
			X1 uint32
			X2 int32
			X3 uint32
			X4 int32
			X5 int16
			X6 int16
			X7 struct {
				X0 unsafe.Pointer
				X1 unsafe.Pointer
			}
		}
		X1 [40]int8
		X2 int64
	}
	U [40]byte
}

type Xpthread_attr_t struct {
	X [0]struct {
		X0 [56]int8
		X1 int64
	}
	U [56]byte
}

// pthread_t pthread_self(void);
func Xpthread_self(tls *TLS) uint64 {
	threadID := tls.threadID
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_self() %v\n", threadID)
	}
	return uint64(threadID)
}

// extern int pthread_equal(pthread_t __thread1, pthread_t __thread2);
func Xpthread_equal(tls *TLS, thread1, thread2 uint64) int32 {
	if thread1 == thread2 {
		return 1
	}

	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_equal(%v, %v) %v\n", thread1, thread2, r)
	}
	return r
}

// int pthread_join(pthread_t thread, void **value_ptr);
func Xpthread_join(tls *TLS, thread uint64, value_ptr *unsafe.Pointer) int32 {
	panic("TODO")
}

// int pthread_create(pthread_t *restrict thread, const pthread_attr_t *restrict attr, void *(*start_routine)(void*), void *restrict arg);
func Xpthread_create(tls *TLS, thread *uint64, attr *Xpthread_attr_t, start_routine func(*TLS, unsafe.Pointer) unsafe.Pointer, arg unsafe.Pointer) int32 {
	panic("TODO")
}

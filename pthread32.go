// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386 arm arm64be armbe mips mipsle ppc ppc64le s390 s390x sparc

package crt

import (
	"fmt"
	"os"
	"unsafe"
)

const (
	Tpthread_attr_t  = "union{[36]int8,int32}"
	Tpthread_mutex_t = "union{struct{int32,uint32,int32,int32,uint32,union{struct{int16,int16},struct{*struct{}}}},[24]int8,int32}"
)

type Xpthread_mutex_t struct {
	X [0]struct {
		X0 struct {
			X0 int32
			X1 uint32
			X2 int32
			X3 int32
			X4 uint32
			X5 struct {
				X [0]struct {
					X0 struct {
						X0 int16
						X1 int16
					}
					X1 struct{ X0 unsafe.Pointer }
				}
				U [4]byte
			}
		}
		X1 [24]int8
		X2 int32
	}
	U [24]byte
}

type Xpthread_attr_t struct {
	X [0]struct {
		X0 [36]int8
		X1 int32
	}
	U [36]byte
}

// pthread_t pthread_self(void);
func Xpthread_self(tls *TLS) uint32 {
	threadID := tls.threadID
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_self() %v\n", threadID)
	}
	return uint32(threadID)
}

// extern int pthread_equal(pthread_t __thread1, pthread_t __thread2);
func Xpthread_equal(tls *TLS, thread1, thread2 uint32) int32 {
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
func Xpthread_join(tls *TLS, thread uint32, value_ptr *unsafe.Pointer) int32 {
	panic("TODO")
}

// int pthread_create(pthread_t *restrict thread, const pthread_attr_t *restrict attr, void *(*start_routine)(void*), void *restrict arg);
func Xpthread_create(tls *TLS, thread *uint32, attr *Xpthread_attr_t, start_routine func(*TLS, unsafe.Pointer) unsafe.Pointer, arg unsafe.Pointer) int32 {
	panic("TODO")
}

// int pthread_detach(pthread_t thread);
func Xpthread_detach(tls *TLS, thread uint32) int32 {
	panic("TODO")
}

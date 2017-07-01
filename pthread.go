// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/cznic/ccir/libc/pthread"
)

const Tpthread_mutexattr_t = "union{[4]int8,int32}"

type Xpthread_mutexattr_t struct {
	X [0]struct {
		X0 [4]int8
		X1 int32
	}
	U [4]byte
}

type mu struct {
	attr  int32
	count int
	inner sync.Mutex
	outer sync.Mutex
	owner uintptr
}

type mutexMap struct {
	m map[unsafe.Pointer]*mu
	sync.Mutex
}

func (m *mutexMap) mu(p unsafe.Pointer) *mu {
	m.Lock()
	r := m.m[p]
	if r == nil {
		r = &mu{}
		m.m[p] = r
	}
	m.Unlock()
	return r
}

var (
	mutexes = &mutexMap{m: map[unsafe.Pointer]*mu{}}
)

// extern int pthread_mutexattr_init(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_init(tls *TLS, attr *Xpthread_mutexattr_t) int32 {
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_init(%#x) %v\n", attr, r)
	}
	return r
}

// extern int pthread_mutexattr_settype(pthread_mutexattr_t * __attr, int __kind);
func Xpthread_mutexattr_settype(tls *TLS, attr *Xpthread_mutexattr_t, kind int32) int32 {
	*(*int32)(unsafe.Pointer(attr)) = kind
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_settype(%#x, %v) %v\n", attr, kind, r)
	}
	return r
}

// extern int pthread_mutex_init(pthread_mutex_t * __mutex, pthread_mutexattr_t * __mutexattr);
func Xpthread_mutex_init(tls *TLS, mutex *Xpthread_mutex_t, mutexattr *Xpthread_mutexattr_t) int32 {
	attr := int32(pthread.XPTHREAD_MUTEX_NORMAL)
	if mutexattr != nil {
		attr = *(*int32)(unsafe.Pointer(mutexattr))
	}
	mutexes.mu(unsafe.Pointer(mutex)).attr = attr
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_init(%p, %#x) %v\n", mutex, mutexattr, r)
	}
	return 0
}

// extern int pthread_mutexattr_destroy(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_destroy(tls *TLS, attr *Xpthread_mutexattr_t) int32 {
	*(*int32)(unsafe.Pointer(attr)) = -1
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_destroy(%#x) %v\n", attr, r)
	}
	return r
}

// extern int pthread_mutex_destroy(pthread_mutex_t * __mutex);
func Xpthread_mutex_destroy(tls *TLS, mutex *Xpthread_mutex_t) int32 {
	mutexes.Lock()
	delete(mutexes.m, unsafe.Pointer(mutex))
	mutexes.Unlock()
	var r int32
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_destroy(%p) %v\n", mutex, r)
	}
	return r
}

// extern int pthread_mutex_lock(pthread_mutex_t * __mutex);
func Xpthread_mutex_lock(tls *TLS, mutex *Xpthread_mutex_t) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(unsafe.Pointer(mutex))
	var r int32
	mu.outer.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		mu.owner = threadID
		mu.count = 1
		mu.inner.Lock()
	case pthread.XPTHREAD_MUTEX_RECURSIVE:
		switch mu.owner {
		case 0:
			mu.owner = threadID
			mu.count = 1
			mu.inner.Lock()
		case threadID:
			mu.count++
		default:
			panic("TODO105")
		}
	default:
		panic(fmt.Errorf("attr %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_lock(%p: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.outer.Unlock()
	return r
}

// int pthread_mutex_trylock(pthread_mutex_t *mutex);
func Xpthread_mutex_trylock(tls *TLS, mutex *Xpthread_mutex_t) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(unsafe.Pointer(mutex))
	var r int32
	mu.outer.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		switch mu.owner {
		case 0:
			mu.owner = threadID
			mu.count = 1
			mu.inner.Lock()
		case threadID:
			panic("TODO127")
		default:
			panic("TODO129")
		}
	default:
		panic(fmt.Errorf("attr %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_trylock(%p: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.outer.Unlock()
	return r
}

// extern int pthread_mutex_unlock(pthread_mutex_t * __mutex);
func Xpthread_mutex_unlock(tls *TLS, mutex *Xpthread_mutex_t) int32 {
	threadID := tls.threadID
	mu := mutexes.mu(unsafe.Pointer(mutex))
	var r int32
	mu.outer.Lock()
	switch mu.attr {
	case pthread.XPTHREAD_MUTEX_NORMAL:
		mu.owner = 0
		mu.count = 0
		mu.inner.Unlock()
	case pthread.XPTHREAD_MUTEX_RECURSIVE:
		switch mu.owner {
		case 0:
			panic("TODO140")
		case threadID:
			mu.count--
			if mu.count != 0 {
				break
			}

			mu.owner = 0
			mu.inner.Unlock()
		default:
			panic("TODO144")
		}
	default:
		panic(fmt.Errorf("TODO %#x", mu.attr))
	}
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_unlock(%p: %+v [thread id %v]) %v\n", mutex, mu, threadID, r)
	}
	mu.outer.Unlock()
	return r
}

// int pthread_cond_wait(pthread_cond_t *cond, pthread_mutex_t *mutex);
func Xpthread_cond_wait(tls *TLS, cons *Xpthread_cond_t, mutex *Xpthread_mutex_t) int32 {
	panic("TODO")
}

// int pthread_cond_signal(pthread_cond_t *cond);
func Xpthread_cond_signal(tls *TLS, cond *Xpthread_cond_t) int32 {
	panic("TODO")
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"sync"
	"unsafe"

	"github.com/cznic/ccir/libc/pthread"
)

type pthread_mutexattr_t *struct {
	X [0]struct {
		X0 [4]int8
		X1 int32
	}
	U [4]byte
}

type pthread_mutex_t *struct {
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

type mu struct {
	attr  int32
	count int
	inner sync.Mutex
	outer sync.Mutex
	owner int
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
func Xpthread_mutexattr_init(attr pthread_mutexattr_t) int32 {
	return 0
}

// extern int pthread_mutexattr_settype(pthread_mutexattr_t * __attr, int __kind);
func Xpthread_mutexattr_settype(attr pthread_mutexattr_t, kind int32) int32 {
	*(*int32)(unsafe.Pointer(attr)) = kind
	return 0
}

// extern int pthread_mutex_init(pthread_mutex_t * __mutex, pthread_mutexattr_t * __mutexattr);
func Xpthread_mutex_init(mutex pthread_mutex_t, mutexattr pthread_mutexattr_t) int32 {
	attr := int32(pthread.XPTHREAD_MUTEX_NORMAL)
	if mutexattr != nil {
		attr = *(*int32)(unsafe.Pointer(mutexattr))
	}
	mutexes.mu(unsafe.Pointer(mutex)).attr = attr
	return 0
}

// extern int pthread_mutexattr_destroy(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_destroy(attr pthread_mutexattr_t) int32 {
	*(*int32)(unsafe.Pointer(attr)) = -1
	return 0
}

// extern int pthread_mutex_destroy(pthread_mutex_t * __mutex);
func Xpthread_mutex_destroy(mutex pthread_mutex_t) int32 {
	mutexes.Lock()
	delete(mutexes.m, unsafe.Pointer(mutex))
	mutexes.Unlock()
	return 0
}

// extern int pthread_mutex_lock(pthread_mutex_t * __mutex);
func Xpthread_mutex_lock(mutex pthread_mutex_t) int32 {
	TODO("")
	panic("TODO")
	//TODO	threadID := c.tlsp.threadID
	//TODO	mutex := readPtr(c.sp)
	//TODO	mu := mutexes.mu(mutex)
	//TODO	var r int32
	//TODO	mu.outer.Lock()
	//TODO	switch mu.attr {
	//TODO	case pthread.XPTHREAD_MUTEX_NORMAL:
	//TODO		mu.owner = threadID
	//TODO		mu.count = 1
	//TODO		mu.inner.Lock()
	//TODO	case pthread.XPTHREAD_MUTEX_RECURSIVE:
	//TODO		switch mu.owner {
	//TODO		case 0:
	//TODO			mu.owner = threadID
	//TODO			mu.count = 1
	//TODO			mu.inner.Lock()
	//TODO		case threadID:
	//TODO			mu.count++
	//TODO		default:
	//TODO			panic("TODO105")
	//TODO		}
	//TODO	default:
	//TODO		panic(fmt.Errorf("attr %#x", mu.attr))
	//TODO	}
	//TODO	mu.outer.Unlock()
	//TODO	writeI32(c.rp, r)
	//TODO	if ptrace {
	//TODO		fmt.Fprintf(os.Stderr, "pthread_mutex_lock(%#x [thread id %v]) %v\n", mutex, threadID, r)
	//TODO	}
}

// pthread_t pthread_self(void);
func Xpthread_self() uint64 {
	TODO("")
	panic("TODO")
	//TODO	threadID := uint64(c.tlsp.threadID)
	//TODO	writeULong(c.rp, threadID)
	//TODO	if ptrace {
	//TODO		fmt.Fprintf(os.Stderr, "pthread_self() %v\n", threadID)
	//TODO	}
}

// extern int pthread_equal(pthread_t __thread1, pthread_t __thread2);
func Xpthread_equal(thread1, thread2 uint64) int32 {
	if thread1 == thread2 {
		return 1
	}

	return 0
}

// int pthread_mutex_trylock(pthread_mutex_t *mutex);
func Xpthread_mutex_trylock(mutex pthread_mutex_t) int32 {
	TODO("")
	panic("TODO")
	//TODO	threadID := c.tlsp.threadID
	//TODO	mutex := readPtr(c.sp)
	//TODO	mu := mutexes.mu(mutex)
	//TODO	var r int32
	//TODO	mu.outer.Lock()
	//TODO	switch mu.attr {
	//TODO	case pthread.XPTHREAD_MUTEX_NORMAL:
	//TODO		switch mu.owner {
	//TODO		case 0:
	//TODO			mu.owner = threadID
	//TODO			mu.count = 1
	//TODO			mu.inner.Lock()
	//TODO		case threadID:
	//TODO			panic("TODO127")
	//TODO		default:
	//TODO			panic("TODO129")
	//TODO		}
	//TODO	default:
	//TODO		panic(fmt.Errorf("attr %#x", mu.attr))
	//TODO	}
	//TODO	mu.outer.Unlock()
	//TODO	writeI32(c.rp, r)
	//TODO	if ptrace {
	//TODO		fmt.Fprintf(os.Stderr, "pthread_mutex_trylock(%#x [thread id %v]) %v\n", mutex, threadID, r)
	//TODO	}
}

// extern int pthread_mutex_unlock(pthread_mutex_t * __mutex);
func Xpthread_mutex_unlock(mutex pthread_mutex_t) int32 {
	TODO("")
	panic("TODO")
	//TODO	threadID := c.tlsp.threadID
	//TODO	mutex := readPtr(c.sp)
	//TODO	mu := mutexes.mu(readPtr(c.sp))
	//TODO	var r int32
	//TODO	mu.outer.Lock()
	//TODO	switch mu.attr {
	//TODO	case pthread.XPTHREAD_MUTEX_NORMAL:
	//TODO		mu.owner = 0
	//TODO		mu.count = 0
	//TODO		mu.inner.Unlock()
	//TODO	case pthread.XPTHREAD_MUTEX_RECURSIVE:
	//TODO		switch mu.owner {
	//TODO		case 0:
	//TODO			panic("TODO140")
	//TODO		case threadID:
	//TODO			mu.count--
	//TODO			if mu.count != 0 {
	//TODO				break
	//TODO			}
	//TODO
	//TODO			mu.owner = 0
	//TODO			mu.inner.Unlock()
	//TODO		default:
	//TODO			panic("TODO144")
	//TODO		}
	//TODO	default:
	//TODO		panic(fmt.Errorf("TODO %#x", mu.attr))
	//TODO	}
	//TODO	mu.outer.Unlock()
	//TODO	writeI32(c.rp, r)
	//TODO	if ptrace {
	//TODO		fmt.Fprintf(os.Stderr, "pthread_mutex_unlock(%#x [thread id %v]) %v\n", mutex, threadID, r)
	//TODO	}
}

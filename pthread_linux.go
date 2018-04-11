// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/cznic/crt/errno"
	"github.com/cznic/crt/pthread"
)

var (
	conds   = &condMap{m: map[uintptr]*cond{}}
	mutexes = &mutexMap{m: map[uintptr]*mu{}}

	// thread-specific data management
	pthreadDB            = map[pthreadDBKey]uintptr{}
	pthreadDBDestructors = map[uint32]uintptr{} // key: destructor
	pthreadDBMu          sync.Mutex
	pthreadDBNextKey     pthread_key_t
)

type pthreadDBKey struct {
	thread uintptr
	key    pthread_key_t
}

type cond struct {
	sync.Mutex // Go
}

type condMap struct {
	m map[uintptr]*cond
	sync.Mutex
}

func (m *condMap) cond(p uintptr) *cond {
	m.Lock()
	r := m.m[p]
	if r == nil {
		r = &cond{}
		m.m[p] = r
	}
	m.Unlock()
	return r
}

type mu struct {
	*sync.Cond
	attr       Spthread_mutexattr_t // PTHREAD_MUTEX_NORMAL, ...
	count      int
	owner      uintptr
	sync.Mutex // Go
}

type mutexMap struct {
	m map[uintptr]*mu
	sync.Mutex
}

func (m *mutexMap) mu(p uintptr) *mu {
	m.Lock()
	r := m.m[p]
	if r == nil {
		r = &mu{attr: Spthread_mutexattr_t{X: pthread.CPTHREAD_MUTEX_DEFAULT}}
		r.Cond = sync.NewCond(&r.Mutex)
		m.m[p] = r
	}
	m.Unlock()
	return r
}

// extern int pthread_mutexattr_init(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_init(tls TLS, attr uintptr) (r int32) {
	// The pthread_mutexattr_init() function shall initialize a mutex
	// attributes object attr with the default value for all of the
	// attributes defined by the implementation.
	//
	// Results are undefined if pthread_mutexattr_init() is called
	// specifying an already initialized attr attributes object.
	//
	// After a mutex attributes object has been used to initialize one or
	// more mutexes, any function affecting the attributes object
	// (including destruction) shall not affect any previously initialized
	// mutexes.
	//
	//
	// Upon successful completion, pthread_mutexattr_destroy() and
	// pthread_mutexattr_init() shall return zero; otherwise, an error
	// number shall be returned to indicate the error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_init.%v(%#x) ", tid, attr)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	*(*Spthread_mutexattr_t)(unsafe.Pointer(attr)) = Spthread_mutexattr_t{X: pthread.CPTHREAD_MUTEX_DEFAULT}
	return 0
}

// extern int pthread_mutexattr_destroy(pthread_mutexattr_t * __attr);
func Xpthread_mutexattr_destroy(tls TLS, attr uintptr) (r int32) {
	// The pthread_mutexattr_destroy() function shall destroy a mutex
	// attributes object; the object becomes, in effect, uninitialized. An
	// implementation may cause pthread_mutexattr_destroy() to set the
	// object referenced by attr to an invalid value.
	//
	// A destroyed attr attributes object can be reinitialized using
	// pthread_mutexattr_init(); the results of otherwise referencing the
	// object after it has been destroyed are undefined.
	//
	//
	// After a mutex attributes object has been used to initialize one or
	// more mutexes, any function affecting the attributes object
	// (including destruction) shall not affect any previously initialized
	// mutexes.
	//
	// The behavior is undefined if the value specified by the attr
	// argument to pthread_mutexattr_destroy() does not refer to an
	// initialized mutex attributes object.
	//
	// Upon successful completion, pthread_mutexattr_destroy() and
	// pthread_mutexattr_init() shall return zero; otherwise, an error
	// number shall be returned to indicate the error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_destroy.%v(%#x) ", tid, attr)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	(*Spthread_mutexattr_t)(unsafe.Pointer(attr)).X = -1
	return 0
}

// extern int pthread_mutexattr_settype(pthread_mutexattr_t * __attr, int __kind);
func Xpthread_mutexattr_settype(tls TLS, attr uintptr, kind int32) (r int32) {
	// The pthread_mutexattr_gettype() and pthread_mutexattr_settype()
	// functions, respectively, shall get and set the mutex type attribute.
	// This attribute is set in the type parameter to these functions. The
	// default value of the type attribute is PTHREAD_MUTEX_DEFAULT.
	//
	// The type of mutex is contained in the type attribute of the mutex
	// attributes. Valid mutex types include:
	//
	//  PTHREAD_MUTEX_NORMAL PTHREAD_MUTEX_ERRORCHECK
	//  PTHREAD_MUTEX_RECURSIVE PTHREAD_MUTEX_DEFAULT
	//
	// The mutex type affects the behavior of calls which lock and unlock
	// the mutex. See pthread_mutex_lock for details. An implementation may
	// map PTHREAD_MUTEX_DEFAULT to one of the other mutex types.
	//
	// The behavior is undefined if the value specified by the attr
	// argument to pthread_mutexattr_gettype() or
	// pthread_mutexattr_settype() does not refer to an initialized mutex
	// attributes object.
	//
	// Upon successful completion, the pthread_mutexattr_gettype() function
	// shall return zero and store the value of the type attribute of attr
	// into the object referenced by the type parameter. Otherwise, an
	// error shall be returned to indicate the error.
	//
	// If successful, the pthread_mutexattr_settype() function shall return
	// zero; otherwise, an error number shall be returned to indicate the
	// error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutexattr_settype.%v(%#x, %v) ", tid, attr, kind)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	(*Spthread_mutexattr_t)(unsafe.Pointer(attr)).X = kind
	return 0
}

// extern int pthread_mutex_init(pthread_mutex_t * __mutex, pthread_mutexattr_t * __mutexattr);
func Xpthread_mutex_init(tls TLS, mutex, mutexattr uintptr) (r int32) {
	// The pthread_mutex_init() function shall initialize the mutex
	// referenced by mutex with attributes specified by attr. If attr is
	// NULL, the default mutex attributes are used; the effect shall be the
	// same as passing the address of a default mutex attributes object.
	// Upon successful initialization, the state of the mutex becomes
	// initialized and unlocked.
	//
	// See Synchronization Object Copies and Alternative Mappings for
	// further requirements.
	//
	// Attempting to initialize an already initialized mutex results in
	// undefined behavior.
	//
	// The behavior is undefined if the value specified by the attr
	// argument to pthread_mutex_init() does not refer to an initialized
	// mutex attributes object.
	//
	// If successful, the pthread_mutex_destroy() and pthread_mutex_init()
	// functions shall return zero; otherwise, an error number shall be
	// returned to indicate the error.
	tid := tls.getThreadID()
	mu := mutexes.mu(mutex)
	mu.Lock()

	defer mu.Unlock()

	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_init.%v(%#x, %#x) %+v, ", tid, mutex, mutexattr, mu)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v, %+v\n", r, mu)
		}()
	}

	attr := Spthread_mutexattr_t{X: pthread.CPTHREAD_MUTEX_DEFAULT}
	if mutexattr != 0 {
		attr = *(*Spthread_mutexattr_t)(unsafe.Pointer(mutexattr))
	}
	mu.attr = attr
	return 0
}

// extern int pthread_mutex_destroy(pthread_mutex_t * __mutex);
func Xpthread_mutex_destroy(tls TLS, mutex uintptr) (r int32) {
	// The pthread_mutex_destroy() function shall destroy the mutex object
	// referenced by mutex; the mutex object becomes, in effect,
	// uninitialized. An implementation may cause pthread_mutex_destroy()
	// to set the object referenced by mutex to an invalid value.
	//
	// A destroyed mutex object can be reinitialized using
	// pthread_mutex_init(); the results of otherwise referencing the
	// object after it has been destroyed are undefined.
	//
	// It shall be safe to destroy an initialized mutex that is unlocked.
	// Attempting to destroy a locked mutex, or a mutex that another thread
	// is attempting to lock, or a mutex that is being used in a
	// pthread_cond_timedwait() or pthread_cond_wait() call by another
	// thread, results in undefined behavior.
	//
	// The pthread_mutex_init() function shall initialize the mutex
	// referenced by mutex with attributes specified by attr. If attr is
	// NULL, the default mutex attributes are used; the effect shall be the
	// same as passing the address of a default mutex attributes object.
	// Upon successful initialization, the state of the mutex becomes
	// initialized and unlocked.
	//
	// See Synchronization Object Copies and Alternative Mappings for
	// further requirements.
	//
	// Attempting to initialize an already initialized mutex results in
	// undefined behavior.
	//
	// The behavior is undefined if the value specified by the mutex
	// argument to pthread_mutex_destroy() does not refer to an initialized
	// mutex.
	//
	// The behavior is undefined if the value specified by the attr
	// argument to pthread_mutex_init() does not refer to an initialized
	// mutex attributes object.
	//
	// If successful, the pthread_mutex_destroy() and pthread_mutex_init()
	// functions shall return zero; otherwise, an error number shall be
	// returned to indicate the error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_destroy.%v(%#x) ", tid, mutex)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	mutexes.Lock()
	delete(mutexes.m, mutex)
	mutexes.Unlock()
	return 0
}

// extern int pthread_mutex_lock(pthread_mutex_t * __mutex);
func Xpthread_mutex_lock(tls TLS, mutex uintptr) (r int32) {
	// The mutex object referenced by mutex shall be locked by a call to
	// pthread_mutex_lock() that returns zero or [EOWNERDEAD]. If the mutex
	// is already locked by another thread, the calling thread shall block
	// until the mutex becomes available. This operation shall return with
	// the mutex object referenced by mutex in the locked state with the
	// calling thread as its owner. If a thread attempts to relock a mutex
	// that it has already locked, pthread_mutex_lock() shall behave as
	// described in the Relock column of the following table. If a thread
	// attempts to unlock a mutex that it has not locked or a mutex which
	// is unlocked, pthread_mutex_unlock() shall behave as described in the
	// Unlock When Not Owner column of the following table.
	//
	//  +------------+------------+--------------------+-----------------------+
	//  | Mutex Type | Robustness | Relock             | Unlock When Not Owner |
	//  +------------+------------+--------------------+-----------------------+
	//  | NORMAL     | non-robust | deadlock           | undefined behavior    |
	//  +------------+------------+--------------------+-----------------------+
	//  | NORMAL     | robust     | deadlock           | error returned        |
	//  +------------+------------+--------------------+-----------------------+
	//  | ERRORCHECK | either     | error returned     | error returned        |
	//  +------------+------------+--------------------+-----------------------+
	//  | RECURSIVE  | either     | recursive          | error returned        |
	//  +------------+------------+--------------------+-----------------------+
	//  | DEFAULT    | non-robust | undefined behavior | undefined behavior    |
	//  +------------+------------+--------------------+-----------------------+
	//  | DEFAULT    | robust     | undefined behavior | error returned        |
	//  +------------+------------+--------------------+-----------------------+
	//
	// Where the table indicates recursive behavior, the mutex shall
	// maintain the concept of a lock count. When a thread successfully
	// acquires a mutex for the first time, the lock count shall be set to
	// one. Every time a thread relocks this mutex, the lock count shall be
	// incremented by one. Each time the thread unlocks the mutex, the lock
	// count shall be decremented by one. When the lock count reaches zero,
	// the mutex shall become available for other threads to acquire.
	//
	// The pthread_mutex_trylock() function shall be equivalent to
	// pthread_mutex_lock(), except that if the mutex object referenced by
	// mutex is currently locked (by any thread, including the current
	// thread), the call shall return immediately. If the mutex type is
	// PTHREAD_MUTEX_RECURSIVE and the mutex is currently owned by the
	// calling thread, the mutex lock count shall be incremented by one and
	// the pthread_mutex_trylock() function shall immediately return
	// success.
	//
	// The pthread_mutex_unlock() function shall release the mutex object
	// referenced by mutex. The manner in which a mutex is released is
	// dependent upon the mutex's type attribute. If there are threads
	// blocked on the mutex object referenced by mutex when
	// pthread_mutex_unlock() is called, resulting in the mutex becoming
	// available, the scheduling policy shall determine which thread shall
	// acquire the mutex.
	//
	// (In the case of PTHREAD_MUTEX_RECURSIVE mutexes, the mutex shall
	// become available when the count reaches zero and the calling thread
	// no longer has any locks on this mutex.)
	//
	// If a signal is delivered to a thread waiting for a mutex, upon
	// return from the signal handler the thread shall resume waiting for
	// the mutex as if it was not interrupted.
	//
	// If mutex is a robust mutex and the process containing the owning
	// thread terminated while holding the mutex lock, a call to
	// pthread_mutex_lock() shall return the error value [EOWNERDEAD]. If
	// mutex is a robust mutex and the owning thread terminated while
	// holding the mutex lock, a call to pthread_mutex_lock() may return
	// the error value [EOWNERDEAD] even if the process in which the owning
	// thread resides has not terminated. In these cases, the mutex is
	// locked by the thread but the state it protects is marked as
	// inconsistent. The application should ensure that the state is made
	// consistent for reuse and when that is complete call
	// pthread_mutex_consistent(). If the application is unable to recover
	// the state, it should unlock the mutex without a prior call to
	// pthread_mutex_consistent(), after which the mutex is marked
	// permanently unusable.
	//
	// If mutex does not refer to an initialized mutex object, the behavior
	// of pthread_mutex_lock(), pthread_mutex_trylock(), and
	// pthread_mutex_unlock() is undefined.
	//
	// If successful, the pthread_mutex_lock(), pthread_mutex_trylock(),
	// and pthread_mutex_unlock() functions shall return zero; otherwise,
	// an error number shall be returned to indicate the error.
	tid := tls.getThreadID()
	mu := mutexes.mu(mutex)
	mu.Lock()

	defer mu.Unlock()

	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_lock.%v(%#x) %+v, ", tid, mutex, mu)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v, %+v\n", r, mu)
		}()
	}

	switch mu.attr.X {
	case pthread.CPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			mu.owner = tid
			mu.count = 1
			return 0
		}

		if mu.owner == tid {
			return errno.XEDEADLK
		}

		for mu.count != 0 {
			mu.Wait()
		}
		mu.owner = tid
		mu.count = 1
		return 0
	case pthread.CPTHREAD_MUTEX_RECURSIVE:
		if mu.count == 0 {
			mu.owner = tid
			mu.count = 1
			return 0
		}

		if mu.owner == tid {
			mu.count++
			return 0
		}

		panic("TODO")
	default:
		panic(mu.attr.X)
	}
}

// extern int pthread_mutex_unlock(pthread_mutex_t * __mutex);
func Xpthread_mutex_unlock(tls TLS, mutex uintptr) (r int32) {
	// Documentation: see Xpthread_mutex_lock
	tid := tls.getThreadID()
	mu := mutexes.mu(mutex)
	mu.Lock()

	defer mu.Unlock()

	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_unlock.%v(%#x) %+v, ", tid, mutex, mu)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v, %+v\n", r, mu)
		}()
	}

	switch mu.attr.X {
	case pthread.CPTHREAD_MUTEX_NORMAL:
		if mu.owner == tid {
			mu.owner = 0
			mu.count = 0
			mu.Broadcast()
			return 0
		}

		panic("TODO")
	case pthread.CPTHREAD_MUTEX_RECURSIVE:
		if mu.owner == tid {
			if mu.count == 0 {
				panic("TODO")
			}

			mu.count--
			if mu.count != 0 {
				return 0
			}

			mu.owner = 0
			mu.Broadcast()
			return 0
		}

		panic("TODO")
	default:
		panic(mu.attr.X)
	}
}

// int pthread_mutex_trylock(pthread_mutex_t *mutex);
func Xpthread_mutex_trylock(tls TLS, mutex uintptr) (r int32) {
	// Documentation: see Xpthread_mutex_lock
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_mutex_trylock.%v(%#x) ", tid, mutex)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	mu := mutexes.mu(mutex)
	mu.Lock()

	defer mu.Unlock()

	switch mu.attr.X {
	case pthread.CPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			mu.count = 1
			mu.owner = threadID
			return
		}

		return errno.XEBUSY
	default:
		panic(mu.attr.X)
	}
}

// pthread_t pthread_self(void);
func Xpthread_self(tls TLS) (r pthread_t) {
	// The pthread_self() function shall return the thread ID of the
	// calling thread.
	//
	// The pthread_self() function shall always be successful and no return
	// value is reserved to indicate an error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_self.%v() ", tid)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	return pthread_t(tid)
}

// extern int pthread_equal(pthread_t __thread1, pthread_t __thread2);
func Xpthread_equal(tls TLS, thread1, thread2 pthread_t) (r int32) {
	// This function shall compare the thread IDs t1 and t2.
	//
	// The pthread_equal() function shall return a non-zero value if t1 and
	// t2 are equal; otherwise, zero shall be returned.
	//
	// If either t1 or t2 are not valid thread IDs, the behavior is
	// undefined.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_equal.%v(%v, %v) ", tid, thread1, thread2)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v\n", r)
		}()
	}
	if thread1 == thread2 {
		r = 1
	}
	return r
}

// int pthread_join(pthread_t thread, void **value_ptr);
func Xpthread_join(tls TLS, thread pthread_t, value_ptr uintptr) int32 {
	panic("TODO")
}

// int pthread_create(pthread_t *restrict thread, const pthread_attr_t *restrict attr, void *(*start_routine)(void*), void *restrict arg);
func Xpthread_create(tls TLS, thread, attr, start_routine, arg uintptr) int32 {
	panic("TODO")
}

// int pthread_key_create(pthread_key_t *pkey, void (*destructor)(void*));
func Xpthread_key_create(tls TLS, pkey, destructor uintptr) (r int32) {
	// The pthread_key_create() function shall create a thread-specific
	// data key visible to all threads in the process. Key values provided
	// by pthread_key_create() are opaque objects used to locate
	// thread-specific data. Although the same key value may be used by
	// different threads, the values bound to the key by
	// pthread_setspecific() are maintained on a per-thread basis and
	// persist for the life of the calling thread.
	//
	// Upon key creation, the value NULL shall be associated with the new
	// key in all active threads. Upon thread creation, the value NULL
	// shall be associated with all defined keys in the new thread.
	//
	// An optional destructor function may be associated with each key
	// value. At thread exit, if a key value has a non-NULL destructor
	// pointer, and the thread has a non-NULL value associated with that
	// key, the value of the key is set to NULL, and then the function
	// pointed to is called with the previously associated value as its
	// sole argument. The order of destructor calls is unspecified if more
	// than one destructor exists for a thread when it exits.
	//
	// If, after all the destructors have been called for all non-NULL
	// values with associated destructors, there are still some non-NULL
	// values with associated destructors, then the process is repeated.
	// If, after at least {PTHREAD_DESTRUCTOR_ITERATIONS} iterations of
	// destructor calls for outstanding non-NULL values, there are still
	// some non-NULL values with associated destructors, implementations
	// may stop calling destructors, or they may continue calling
	// destructors until no non-NULL values with associated destructors
	// exist, even though this might result in an infinite loop.
	//
	// If successful, the pthread_key_create() function shall store the
	// newly created key value at *key and shall return zero. Otherwise, an
	// error number shall be returned to indicate the error.
	tid := tls.getThreadID()
	var key pthread_key_t
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_key_create(%v, %#x, %#x) ", tid, pkey, destructor)
		defer func() {
			fmt.Fprintf(os.Stderr, "%v (key %#x)\n", r, key)
		}()
	}
	pthreadDBMu.Lock()
	pthreadDBNextKey++
	key = pthreadDBNextKey
	if destructor != 0 {
		pthreadDBDestructors[key] = destructor
	}
	*(*pthread_key_t)(unsafe.Pointer(pkey)) = key
	pthreadDBMu.Unlock()
	return 0
}

// void *pthread_getspecific(pthread_key_t key);
func Xpthread_getspecific(tls TLS, key pthread_key_t) (r uintptr) {
	// The pthread_getspecific() function shall return the value currently
	// bound to the specified key on behalf of the calling thread.
	//
	// The pthread_setspecific() function shall associate a thread-specific
	// value with a key obtained via a previous call to
	// pthread_key_create(). Different threads may bind different values to
	// the same key. These values are typically pointers to blocks of
	// dynamically allocated memory that have been reserved for use by the
	// calling thread.
	//
	// The effect of calling pthread_getspecific() or pthread_setspecific()
	// with a key value not obtained from pthread_key_create() or after key
	// has been deleted with pthread_key_delete() is undefined.
	//
	// Both pthread_getspecific() and pthread_setspecific() may be called
	// from a thread-specific data destructor function. A call to
	// pthread_getspecific() for the thread-specific data key being
	// destroyed shall return the value NULL, unless the value is changed
	// (after the destructor starts) by a call to pthread_setspecific().
	// Calling pthread_setspecific() from a thread-specific data destructor
	// routine may result either in lost storage (after at least
	// PTHREAD_DESTRUCTOR_ITERATIONS attempts at destruction) or in an
	// infinite loop.
	//
	// The pthread_getspecific() function shall return the thread-specific
	// data value associated with the given key. If no thread-specific data
	// value is associated with key, then the value NULL shall be returned.
	//
	// The pthread_getspecific() function shall return the thread-specific
	// data value associated with the given key. If no thread-specific data
	// value is associated with key, then the value NULL shall be returned.
	//
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_getspecific.%v(%#x) ", tid, key)
		defer func() {
			fmt.Fprintf(os.Stderr, "%#x\n", r)
		}()
	}
	pthreadDBMu.Lock()
	r = pthreadDB[pthreadDBKey{thread: tid, key: key}]
	pthreadDBMu.Unlock()
	return r
}

// int pthread_setspecific(pthread_key_t key, const void *value);
func Xpthread_setspecific(tls TLS, key uint32, value uintptr) (r int32) {
	// The pthread_setspecific() function shall associate a thread-specific
	// value with a key obtained via a previous call to
	// pthread_key_create(). Different threads may bind different values to
	// the same key. These values are typically pointers to blocks of
	// dynamically allocated memory that have been reserved for use by the
	// calling thread.
	//
	// The effect of calling pthread_getspecific() or pthread_setspecific()
	// with a key value not obtained from pthread_key_create() or after key
	// has been deleted with pthread_key_delete() is undefined.
	//
	// Both pthread_getspecific() and pthread_setspecific() may be called
	// from a thread-specific data destructor function. A call to
	// pthread_getspecific() for the thread-specific data key being
	// destroyed shall return the value NULL, unless the value is changed
	// (after the destructor starts) by a call to pthread_setspecific().
	// Calling pthread_setspecific() from a thread-specific data destructor
	// routine may result either in lost storage (after at least
	// PTHREAD_DESTRUCTOR_ITERATIONS attempts at destruction) or in an
	// infinite loop.
	//
	// If successful, the pthread_setspecific() function shall return zero;
	// otherwise, an error number shall be returned to indicate the error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_setspecific.%v(%#x, %#x) ", tid, key, value)
		defer func() {
			fmt.Fprintf(os.Stderr, "%#x\n", r)
		}()
	}
	pthreadDBMu.Lock()
	pthreadDB[pthreadDBKey{thread: tls.getThreadID(), key: key}] = value
	pthreadDBMu.Unlock()
	return 0
}

// int pthread_cond_broadcast(pthread_cond_t *cond);
func Xpthread_cond_broadcast(tls TLS, cond uintptr) int32 {
	panic("TODO")
}

// int pthread_cond_wait(pthread_cond_t *restrict cond, pthread_mutex_t *restrict mutex);
func Xpthread_cond_wait(tls TLS, cond, mutex uintptr) int32 {
	panic("TODO")
}

// int pthread_cond_init(pthread_cond_t *restrict cond, const pthread_condattr_t *restrict attr);
func Xpthread_cond_init(tls TLS, cond, attr uintptr) (r int32) {
	// The pthread_cond_init() function shall initialize the condition
	// variable referenced by cond with attributes referenced by attr. If
	// attr is NULL, the default condition variable attributes shall be
	// used; the effect is the same as passing the address of a default
	// condition variable attributes object. Upon successful
	// initialization, the state of the condition variable shall become
	// initialized.
	//
	// See Synchronization Object Copies and Alternative Mappings for
	// further requirements.
	//
	// Attempting to initialize an already initialized condition variable
	// results in undefined behavior.
	//
	// The behavior is undefined if the value specified by the attr
	// argument to pthread_cond_init() does not refer to an initialized
	// condition variable attributes object.
	//
	// If successful, the pthread_cond_destroy() and pthread_cond_init()
	// functions shall return zero; otherwise, an error number shall be
	// returned to indicate the error.
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(os.Stderr, "pthread_cond_init.%v(%#x, %#x) ", tid, cond, attr)
		defer func() {
			fmt.Fprintf(os.Stderr, "%#x\n", r)
		}()
	}
	// c := conds.cond(cond) //TODO for non-nil attr
	if attr != 0 {
		panic("TODO")
	}
	return 0
}

// int pthread_cond_destroy(pthread_cond_t *cond);
func Xpthread_cond_destroy(tls TLS, cond uintptr) int32 {
	//TODO when pthread_cond_init starts to do something
	return 0
}

// int pthread_atfork(void (*prepare)(void), void (*parent)(void), void (*child)(void));
func Xpthread_atfork(tls TLS, prepare, parent, child uintptr) int32 {
	panic("TODO")
}

// int pthread_key_delete(pthread_key_t key);
func Xpthread_key_delete(tls TLS, key uint32) int32 {
	panic("TODO")
}

// int pthread_attr_init(pthread_attr_t *attr);
func Xpthread_attr_init(tls TLS, attr uintptr) int32 {
	panic("TODO")
}

// int pthread_attr_setscope(pthread_attr_t *attr, int scope);
func Xpthread_attr_setscope(tls TLS, attr uintptr, scope int32) int32 {
	panic("TODO")
}

// int pthread_attr_setstacksize(pthread_attr_t *attr, size_t stacksize);
func Xpthread_attr_setstacksize(tls TLS, attr uintptr, stacksize size_t) int32 {
	panic("TODO")
}

// int pthread_attr_setdetachstate(pthread_attr_t *attr, int detachstate);
func Xpthread_attr_setdetachstate(tls TLS, attr uintptr, detached int32) int32 {
	panic("TODO")
}

// int pthread_attr_destroy(pthread_attr_t *attr);
func Xpthread_attr_destroy(tls TLS, attr uintptr) int32 {
	panic("TODO")
}

// void pthread_exit(void *retval);
func Xpthread_exit(tls TLS, retval uintptr) {
	panic("TODO")
}

// int pthread_detach(pthread_t thread);
func Xpthread_detach(tls TLS, thread pthread_t) int32 {
	panic("TODO")
}

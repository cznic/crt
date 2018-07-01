// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/cznic/crt/errno"
	"github.com/cznic/crt/pthread"
)

var (
	_ sync.Locker = (*locker)(nil)

	attrs   = attrMap{m: map[uintptr]*attr{}}
	conds   = condMap{m: map[uintptr]*cond{}}
	mutexes = mutexMap{m: map[uintptr]*mutex{}}
	threads = &threadMap{m: map[pthread_t]*threadState{}}

	// thread-specific data management
	pthreadDB            = map[pthreadDBKey]uintptr{}
	pthreadDBDestructors = map[uint32]uintptr{} // key: destructor
	pthreadDBMu          sync.Mutex
	pthreadDBNextKey     pthread_key_t
)

type threadExited uintptr

type threadMap struct {
	m map[pthread_t]*threadState
	sync.Mutex
}

type threadState struct {
	attr     *attr
	c        chan struct{}
	detached bool
	retval   uintptr
}

type attr struct {
	scope int32
}

type attrMap struct {
	m map[uintptr]*attr
	sync.Mutex
}

func (m *attrMap) attr(p uintptr) *attr {
	m.Lock()
	r := m.m[p]
	if r == nil {
		r = &attr{}
		m.m[p] = r
	}
	m.Unlock()
	return r
}

type cond struct {
	conds      map[*sync.Cond]struct{}
	signals    int
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
		r = &cond{conds: map[*sync.Cond]struct{}{}}
		m.m[p] = r
	}
	r.Lock()
	m.Unlock()
	return r
}

type pthreadDBKey struct {
	thread uintptr
	key    pthread_key_t
}

type mutexMap struct {
	m  map[uintptr]*mutex
	mu sync.Mutex
}

func (m *mutexMap) mutex(p uintptr) *mutex {
	m.mu.Lock()
	r := m.m[p]
	if r == nil {
		r = &mutex{p: p}
		r.Cond = sync.NewCond(r)
		m.m[p] = r
	}
	r.Lock()
	m.mu.Unlock()
	return r
}

type mutex struct {
	*sync.Cond // Go
	count      int
	owner      uintptr
	p          uintptr
	sync.Mutex // Go
}

func (mu *mutex) setTyp(typ int32) { (*pthread_mutex_t)(unsafe.Pointer(mu.p)).X__kind = typ }
func (mu *mutex) typ() int32       { return (*pthread_mutex_t)(unsafe.Pointer(mu.p)).X__kind }

func (mu *mutex) lock(tid uintptr) int32 {
	defer mu.Unlock()

	switch mu.typ() {
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
		panic(mu.typ())
	}
}

func (mu *mutex) unlock(tid uintptr) int32 {
	defer mu.Unlock()

	switch mu.typ() {
	case pthread.CPTHREAD_MUTEX_NORMAL:
		if mu.owner == tid {
			mu.owner = 0
			mu.count = 0
			mu.Signal()
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
			mu.Signal()
			return 0
		}

		panic("TODO")
	default:
		panic(mu.typ())
	}
}

type locker struct {
	m   *mutex
	tid uintptr
}

func (l *locker) Lock() {
	l.m.Lock()
	l.m.lock(l.tid)
}

func (l *locker) Unlock() {
	l.m.Lock()
	l.m.unlock(l.tid)
}

//  C: int pthread_mutexattr_destroy(pthread_mutexattr_t *attr);
//  C: int pthread_mutexattr_init(pthread_mutexattr_t *attr);
//
// The pthread_mutexattr_destroy() function shall destroy a mutex attributes
// object; the object becomes, in effect, uninitialized. An implementation may
// cause pthread_mutexattr_destroy() to set the object referenced by attr to an
// invalid value.
//
// A destroyed attr attributes object can be reinitialized using
// pthread_mutexattr_init(); the results of otherwise referencing the object
// after it has been destroyed are undefined.
//
// The pthread_mutexattr_init() function shall initialize a mutex attributes
// object attr with the default value for all of the attributes defined by the
// implementation.
//
// Results are undefined if pthread_mutexattr_init() is called specifying an
// already initialized attr attributes object.
//
// After a mutex attributes object has been used to initialize one or more
// mutexes, any function affecting the attributes object (including
// destruction) shall not affect any previously initialized mutexes.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_mutexattr_destroy() does not refer to an initialized mutex
// attributes object.
//
// Upon successful completion, pthread_mutexattr_destroy() and
// pthread_mutexattr_init() shall return zero; otherwise, an error number shall
// be returned to indicate the error.
//
// The pthread_mutexattr_init() function shall fail if:
//
// [ENOMEM] Insufficient memory exists to initialize the mutex attributes
// object.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_mutexattr_destroy(tcl TLS, attr uintptr) (r int32) {
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_mutexattr_destroy(%#x) ", attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	*(*int32)(unsafe.Pointer(attr)) = -1
	return 0
}

// Documentation: see Xpthread_mutexattr_destroy.
func Xpthread_mutexattr_init(tcl TLS, attr uintptr) (r int32) {
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_mutexattr_init(%#x) ", attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	*(*int32)(unsafe.Pointer(attr)) = pthread.CPTHREAD_MUTEX_DEFAULT
	return 0
}

//  C: int pthread_mutexattr_gettype(const pthread_mutexattr_t *restrict attr, int *restrict type);
//  C: int pthread_mutexattr_settype(pthread_mutexattr_t *attr, int type);
//
// The pthread_mutexattr_gettype() and pthread_mutexattr_settype() functions,
// respectively, shall get and set the mutex type attribute. This attribute is
// set in the type parameter to these functions. The default value of the type
// attribute is PTHREAD_MUTEX_DEFAULT.
//
// The type of mutex is contained in the type attribute of the mutex
// attributes. Valid mutex types include:
//
//	PTHREAD_MUTEX_NORMAL PTHREAD_MUTEX_ERRORCHECK PTHREAD_MUTEX_RECURSIVE
//	PTHREAD_MUTEX_DEFAULT
//
// The mutex type affects the behavior of calls which lock and unlock the
// mutex. See pthread_mutex_lock for details. An implementation may map
// PTHREAD_MUTEX_DEFAULT to one of the other mutex types.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_mutexattr_gettype() or pthread_mutexattr_settype() does not refer to
// an initialized mutex attributes object.
//
// Upon successful completion, the pthread_mutexattr_gettype() function shall
// return zero and store the value of the type attribute of attr into the
// object referenced by the type parameter. Otherwise, an error shall be
// returned to indicate the error.
//
// If successful, the pthread_mutexattr_settype() function shall return zero;
// otherwise, an error number shall be returned to indicate the error.
//
// The pthread_mutexattr_settype() function shall fail if:
//
// [EINVAL] The value type is invalid.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_mutexattr_gettype(tls TLS, attr pthread_mutexattr_t, typ int32) {
	panic("TODO")
}

// Documentation: see Xpthread_mutexattr_gettype.
func Xpthread_mutexattr_settype(tls TLS, attr uintptr, typ int32) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_mutexattr_settype.%v(%#x, %v) ", tid, attr, typ)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	switch typ {
	case
		pthread.CPTHREAD_MUTEX_ERRORCHECK,
		pthread.CPTHREAD_MUTEX_NORMAL,
		pthread.CPTHREAD_MUTEX_RECURSIVE:
		*(*int32)(unsafe.Pointer(attr)) = typ
	default:
		r = errno.XEINVAL
	}
	return r
}

//  C: int pthread_mutex_destroy(pthread_mutex_t *mutex);
//  C: int pthread_mutex_init(pthread_mutex_t *restrict mutex, const pthread_mutexattr_t *restrict attr);
//
// The pthread_mutex_destroy() function shall destroy the mutex object
// referenced by mutex; the mutex object becomes, in effect, uninitialized. An
// implementation may cause pthread_mutex_destroy() to set the object
// referenced by mutex to an invalid value.
//
// A destroyed mutex object can be reinitialized using pthread_mutex_init();
// the results of otherwise referencing the object after it has been destroyed
// are undefined.
//
// It shall be safe to destroy an initialized mutex that is unlocked.
// Attempting to destroy a locked mutex, or a mutex that another thread is
// attempting to lock, or a mutex that is being used in a
// pthread_cond_timedwait() or pthread_cond_wait() call by another thread,
// results in undefined behavior.
//
// The pthread_mutex_init() function shall initialize the mutex referenced by
// mutex with attributes specified by attr. If attr is NULL, the default mutex
// attributes are used; the effect shall be the same as passing the address of
// a default mutex attributes object. Upon successful initialization, the state
// of the mutex becomes initialized and unlocked.
//
// See Synchronization Object Copies and Alternative Mappings for further
// requirements.
//
// Attempting to initialize an already initialized mutex results in undefined
// behavior.
//
// In cases where default mutex attributes are appropriate, the macro
// PTHREAD_MUTEX_INITIALIZER can be used to initialize mutexes. The effect
// shall be equivalent to dynamic initialization by a call to
// pthread_mutex_init() with parameter attr specified as NULL, except that no
// error checks are performed.
//
// The behavior is undefined if the value specified by the mutex argument to
// pthread_mutex_destroy() does not refer to an initialized mutex.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_mutex_init() does not refer to an initialized mutex attributes
// object.
//
// If successful, the pthread_mutex_destroy() and pthread_mutex_init()
// functions shall return zero; otherwise, an error number shall be returned to
// indicate the error.
//
// The pthread_mutex_init() function shall fail if:
//
// [EAGAIN] The system lacked the necessary resources (other than memory) to
// initialize another mutex.
//
// [ENOMEM] Insufficient memory exists to initialize the mutex.
//
// [EPERM] The caller does not have the privilege to perform the operation.
//
// The pthread_mutex_init() function may fail if:
//
// [EINVAL] The attributes object referenced by attr has the robust mutex
// attribute set without the process-shared attribute being set.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_mutex_destroy(tls TLS, mutex uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_mutex_destroy.%v(%#x) , ", tid, mutex)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	mu := mutexes.mutex(mutex)
	mu.setTyp(-1)
	mu.Unlock()
	return 0
}

// Documentation: see Xpthread_mutex_destroy
func Xpthread_mutex_init(tls TLS, mutex, attr uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_mutex_init.%v(%#x, %#x) ", tid, mutex, attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	mu := mutexes.mutex(mutex)
	typ := int32(pthread.CPTHREAD_MUTEX_DEFAULT)
	if attr != 0 {
		typ = *(*int32)(unsafe.Pointer(attr))
	}
	mu.setTyp(typ)
	mu.Unlock()
	return 0
}

//  C: int pthread_mutex_lock(pthread_mutex_t *mutex);
//  C: int pthread_mutex_trylock(pthread_mutex_t *mutex);
//  C: int pthread_mutex_unlock(pthread_mutex_t *mutex);
//
// The mutex object referenced by mutex shall be locked by a call to
// pthread_mutex_lock() that returns zero or [EOWNERDEAD]. If the mutex is
// already locked by another thread, the calling thread shall block until the
// mutex becomes available. This operation shall return with the mutex object
// referenced by mutex in the locked state with the calling thread as its
// owner. If a thread attempts to relock a mutex that it has already locked,
// pthread_mutex_lock() shall behave as described in the Relock column of the
// following table. If a thread attempts to unlock a mutex that it has not
// locked or a mutex which is unlocked, pthread_mutex_unlock() shall behave as
// described in the Unlock When Not Owner column of the following table.
//
//	+------------+------------+----------------+-----------------------+
//	| Mutex Type | Robustness | Relock         | Unlock When Not Owner |
//	+------------+------------+----------------+-----------------------+
//	| NORMAL     | non-robust | deadlock       | undefined behavior    |
//	| NORMAL     | robust     | deadlock       | error returned        |
//	| ERRORCHECK | either     | error returned | error returned        |
//	| RECURSIVE  | either     | recursive      | error returned        |
//	| DEFAULT    | non-robust | undefined      | undefined behavior    |
//	| DEFAULT    | robust     | undefined      | error returned        |
//	+------------+------------+----------------+-----------------------+
//
// If the mutex type is PTHREAD_MUTEX_DEFAULT, the behavior of
// pthread_mutex_lock() may correspond to one of the three other standard mutex
// types as described in the table above. If it does not correspond to one of
// those three, the behavior is undefined for the cases marked .  Where the
// table indicates recursive behavior, the mutex shall maintain the concept of
// a lock count. When a thread successfully acquires a mutex for the first
// time, the lock count shall be set to one. Every time a thread relocks this
// mutex, the lock count shall be incremented by one. Each time the thread
// unlocks the mutex, the lock count shall be decremented by one. When the lock
// count reaches zero, the mutex shall become available for other threads to
// acquire.
//
// The pthread_mutex_trylock() function shall be equivalent to
// pthread_mutex_lock(), except that if the mutex object referenced by mutex is
// currently locked (by any thread, including the current thread), the call
// shall return immediately. If the mutex type is PTHREAD_MUTEX_RECURSIVE and
// the mutex is currently owned by the calling thread, the mutex lock count
// shall be incremented by one and the pthread_mutex_trylock() function shall
// immediately return success.
//
// The pthread_mutex_unlock() function shall release the mutex object
// referenced by mutex. The manner in which a mutex is released is dependent
// upon the mutex's type attribute. If there are threads blocked on the mutex
// object referenced by mutex when pthread_mutex_unlock() is called, resulting
// in the mutex becoming available, the scheduling policy shall determine which
// thread shall acquire the mutex.
//
// (In the case of PTHREAD_MUTEX_RECURSIVE mutexes, the mutex shall become
// available when the count reaches zero and the calling thread no longer has
// any locks on this mutex.)
//
// If a signal is delivered to a thread waiting for a mutex, upon return from
// the signal handler the thread shall resume waiting for the mutex as if it
// was not interrupted.
//
// If mutex is a robust mutex and the process containing the owning thread
// terminated while holding the mutex lock, a call to pthread_mutex_lock()
// shall return the error value [EOWNERDEAD]. If mutex is a robust mutex and
// the owning thread terminated while holding the mutex lock, a call to
// pthread_mutex_lock() may return the error value [EOWNERDEAD] even if the
// process in which the owning thread resides has not terminated. In these
// cases, the mutex is locked by the thread but the state it protects is marked
// as inconsistent. The application should ensure that the state is made
// consistent for reuse and when that is complete call
// pthread_mutex_consistent(). If the application is unable to recover the
// state, it should unlock the mutex without a prior call to
// pthread_mutex_consistent(), after which the mutex is marked permanently
// unusable.
//
// If mutex does not refer to an initialized mutex object, the behavior of
// pthread_mutex_lock(), pthread_mutex_trylock(), and pthread_mutex_unlock() is
// undefined.
//
// If successful, the pthread_mutex_lock(), pthread_mutex_trylock(), and
// pthread_mutex_unlock() functions shall return zero; otherwise, an error
// number shall be returned to indicate the error.
//
// The pthread_mutex_lock() and pthread_mutex_trylock() functions shall fail
// if:
//
// [EAGAIN] The mutex could not be acquired because the maximum number of
// recursive locks for mutex has been exceeded.
//
// [EINVAL] The mutex was created with the protocol attribute having the value
// PTHREAD_PRIO_PROTECT and the calling thread's priority is higher than the
// mutex's current priority ceiling.
//
// [ENOTRECOVERABLE] The state protected by the mutex is not recoverable.
//
// [EOWNERDEAD] The mutex is a robust mutex and the process containing the
// previous owning thread terminated while holding the mutex lock. The mutex
// lock shall be acquired by the calling thread and it is up to the new owner
// to make the state consistent.
//
// The pthread_mutex_lock() function shall fail if:
//
// [EDEADLK] The mutex type is PTHREAD_MUTEX_ERRORCHECK and the current thread
// already owns the mutex.
//
// The pthread_mutex_trylock() function shall fail if:
//
// [EBUSY] The mutex could not be acquired because it was already locked.
//
// The pthread_mutex_unlock() function shall fail if:
//
// [EPERM] The mutex type is PTHREAD_MUTEX_ERRORCHECK or
// PTHREAD_MUTEX_RECURSIVE, or the mutex is a robust mutex, and the current
// thread does not own the mutex.
//
// The pthread_mutex_lock() and pthread_mutex_trylock() functions may fail if:
//
// [EOWNERDEAD] The mutex is a robust mutex and the previous owning thread
// terminated while holding the mutex lock. The mutex lock shall be acquired by
// the calling thread and it is up to the new owner to make the state
// consistent.
//
// The pthread_mutex_lock() function may fail if:
//
// [EDEADLK] A deadlock condition was detected.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_mutex_lock(tls TLS, mutex uintptr) (r int32) {
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_mutex_lock.%v(%#x) ", tls.getThreadID(), mutex)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	return mutexes.mutex(mutex).lock(tls.getThreadID())
}

// Documentation: see Xpthread_mutex_lock.
func Xpthread_mutex_trylock(tls TLS, mutex uintptr) (r int32) {
	tid := tls.getThreadID()
	mu := mutexes.mutex(mutex)

	defer mu.Unlock()

	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_mutex_trylock.%v(%#x) ", tid, mutex)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	switch mu.typ() {
	case pthread.CPTHREAD_MUTEX_NORMAL:
		if mu.count == 0 {
			mu.count = 1
			mu.owner = threadID
			return
		}

		return errno.XEBUSY
	default:
		panic(mu.typ())
	}
}

// Documentation: see Xpthread_mutex_lock.
func Xpthread_mutex_unlock(tls TLS, mutex uintptr) (r int32) {
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_mutex_unlock.%v(%#x) ", tls.getThreadID(), mutex)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	return mutexes.mutex(mutex).unlock(tls.getThreadID())
}

//  C: pthread_t pthread_self(void);
//
// The pthread_self() function shall return the thread ID of the calling
// thread.
//
// The pthread_self() function shall always be successful and no return value
// is reserved to indicate an error.
//
// No errors are defined.
func Xpthread_self(tls TLS) (r pthread_t) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_self.%v() ", tid)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	return pthread_t(tid)
}

//  C: int pthread_equal(pthread_t t1, pthread_t t2);
//
// This function shall compare the thread IDs t1 and t2.
//
// The pthread_equal() function shall return a non-zero value if t1 and t2 are
// equal; otherwise, zero shall be returned.
//
// If either t1 or t2 are not valid thread IDs, the behavior is undefined.
//
// No errors are defined.
//
// The pthread_equal() function shall not return an error code of [EINTR].
func Xpthread_equal(tls TLS, t1, t2 pthread_t) (r int32) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_equal.%v(%#x, %#x) ", tid, t1, t2)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	if t1 == t2 {
		r = 1
	}
	return r
}

//  C: int pthread_join(pthread_t thread, void **value_ptr);
//
// The pthread_join() function shall suspend execution of the calling thread
// until the target thread terminates, unless the target thread has already
// terminated. On return from a successful pthread_join() call with a non-NULL
// value_ptr argument, the value passed to pthread_exit() by the terminating
// thread shall be made available in the location referenced by value_ptr. When
// a pthread_join() returns successfully, the target thread has been
// terminated. The results of multiple simultaneous calls to pthread_join()
// specifying the same target thread are undefined. If the thread calling
// pthread_join() is canceled, then the target thread shall not be detached.
//
// It is unspecified whether a thread that has exited but remains unjoined
// counts against {PTHREAD_THREADS_MAX}.
//
// The behavior is undefined if the value specified by the thread argument to
// pthread_join() does not refer to a joinable thread.
//
// The behavior is undefined if the value specified by the thread argument to
// pthread_join() refers to the calling thread.
//
// If successful, the pthread_join() function shall return zero; otherwise, an
// error number shall be returned to indicate the error.
//
// The pthread_join() function may fail if:
//
// [EDEADLK] A deadlock was detected.
//
// The pthread_join() function shall not return an error code of [EINTR].
func Xpthread_join(tls TLS, thread pthread_t, value_ptr uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_join.%v(%v, %#x) ", tid, thread, value_ptr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	threads.Lock()
	t := threads.m[thread]
	threads.Unlock()
	switch {
	case t != nil:
		<-t.c
		if value_ptr != 0 {
			*(*uintptr)(unsafe.Pointer(value_ptr)) = t.retval
		}
	default:
		if value_ptr != 0 {
			panic("TODO")
		}
	}
	return 0
}

func startRoutine(p uintptr) func(TLS, uintptr) uintptr {
	return *(*func(TLS, uintptr) uintptr)(unsafe.Pointer(&p))
}

//  C: int pthread_create(pthread_t *restrict thread, const pthread_attr_t *restrict attr, void *(*start_routine)(void*), void *restrict arg);
//
// The pthread_create() function shall create a new thread, with attributes
// specified by attr, within a process. If attr is NULL, the default attributes
// shall be used. If the attributes specified by attr are modified later, the
// thread's attributes shall not be affected. Upon successful completion,
// pthread_create() shall store the ID of the created thread in the location
// referenced by thread.
//
// The thread is created executing start_routine with arg as its sole argument.
// If the start_routine returns, the effect shall be as if there was an
// implicit call to pthread_exit() using the return value of start_routine as
// the exit status. Note that the thread in which main() was originally invoked
// differs from this. When it returns from main(), the effect shall be as if
// there was an implicit call to exit() using the return value of main() as the
// exit status.
//
// The signal state of the new thread shall be initialized as follows:
//
// The signal mask shall be inherited from the creating thread.
//
// The set of signals pending for the new thread shall be empty.
//
// The thread-local current locale [XSI] ￼  and the alternate stack ￼  shall
// not be inherited.
//
// The floating-point environment shall be inherited from the creating thread.
//
// If pthread_create() fails, no new thread is created and the contents of the
// location referenced by thread are undefined.
//
// _POSIX_THREAD_CPUTIME is defined, the new thread shall have a CPU-time clock
// accessible, and the initial value of this clock shall be set to zero.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_create() does not refer to an initialized thread attributes object.
//
// If successful, the pthread_create() function shall return zero; otherwise,
// an error number shall be returned to indicate the error.
//
// The pthread_create() function shall fail if:
//
// [EAGAIN] The system lacked the necessary resources to create another thread,
// or the system-imposed limit on the total number of threads in a process
// {PTHREAD_THREADS_MAX} would be exceeded.
//
// [EPERM] The caller does not have appropriate privileges to set the required
// scheduling parameters or scheduling policy.
//
// The pthread_create() function shall not return an error code of [EINTR].
func Xpthread_create(tls TLS, thread, pattr, start_routine, arg uintptr) (r int32) {
	newTLS := NewTLS()
	tid := pthread_t(newTLS.getThreadID())
	*(*pthread_t)(unsafe.Pointer(thread)) = tid
	if ptrace {
		id := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_create.%v(%#x, %#x, %#x, %#x) ", id, thread, pattr, start_routine, arg)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v %v\n", tid, r)
		}()
	}

	var a *attr
	switch {
	case pattr != 0:
		a = attrs.attr(pattr)
	default:
		a = &attr{}
	}
	t := &threadState{c: make(chan struct{}), attr: a}
	threads.Lock()
	threads.m[tid] = t
	threads.Unlock()
	ch := make(chan struct{})
	go func() {
		defer func() {
			close(t.c)
			if t.detached {
				threads.Lock()
				delete(threads.m, tid)
				threads.Unlock()
				if ptrace {
					fmt.Fprintf(TraceWriter, "thread #%v was detached", tid)
				}
			}
		}()

		defer func() {
			switch x := recover().(type) {
			case nil:
			case threadExited:
				t.retval = uintptr(x)
				if ptrace {
					fmt.Fprintf(TraceWriter, "thread #%v exited: %#x\n", tid, t.retval)
				}
				return
			default:
				panic(x)
			}
		}()

		close(ch)
		t.retval = startRoutine(start_routine)(newTLS, arg)
		if ptrace {
			fmt.Fprintf(TraceWriter, "thread #%v returned: %#x\n", tid, t.retval)
		}
	}()
	<-ch
	return r

}

//  C: int pthread_cond_timedwait(pthread_cond_t *restrict cond, pthread_mutex_t *restrict mutex, const struct timespec *restrict abstime);
//  C: int pthread_cond_wait(pthread_cond_t *restrict cond, pthread_mutex_t *restrict mutex);
//
// The pthread_cond_timedwait() and pthread_cond_wait() functions shall block
// on a condition variable. The application shall ensure that these functions
// are called with mutex locked by the calling thread; otherwise, an error (for
// PTHREAD_MUTEX_ERRORCHECK and robust mutexes) or undefined behavior (for
// other mutexes) results.
//
// These functions atomically release mutex and cause the calling thread to
// block on the condition variable cond; atomically here means "atomically with
// respect to access by another thread to the mutex and then the condition
// variable". That is, if another thread is able to acquire the mutex after the
// about-to-block thread has released it, then a subsequent call to
// pthread_cond_broadcast() or pthread_cond_signal() in that thread shall
// behave as if it were issued after the about-to-block thread has blocked.
//
// Upon successful return, the mutex shall have been locked and shall be owned
// by the calling thread. If mutex is a robust mutex where an owner terminated
// while holding the lock and the state is recoverable, the mutex shall be
// acquired even though the function returns an error code.
//
// When using condition variables there is always a Boolean predicate involving
// shared variables associated with each condition wait that is true if the
// thread should proceed. Spurious wakeups from the pthread_cond_timedwait() or
// pthread_cond_wait() functions may occur. Since the return from
// pthread_cond_timedwait() or pthread_cond_wait() does not imply anything
// about the value of this predicate, the predicate should be re-evaluated upon
// such return.
//
// When a thread waits on a condition variable, having specified a particular
// mutex to either the pthread_cond_timedwait() or the pthread_cond_wait()
// operation, a dynamic binding is formed between that mutex and condition
// variable that remains in effect as long as at least one thread is blocked on
// the condition variable. During this time, the effect of an attempt by any
// thread to wait on that condition variable using a different mutex is
// undefined. Once all waiting threads have been unblocked (as by the
// pthread_cond_broadcast() operation), the next wait operation on that
// condition variable shall form a new dynamic binding with the mutex specified
// by that wait operation. Even though the dynamic binding between condition
// variable and mutex may be removed or replaced between the time a thread is
// unblocked from a wait on the condition variable and the time that it returns
// to the caller or begins cancellation cleanup, the unblocked thread shall
// always re-acquire the mutex specified in the condition wait operation call
// from which it is returning.
//
// A condition wait (whether timed or not) is a cancellation point. When the
// cancelability type of a thread is set to PTHREAD_CANCEL_DEFERRED, a
// side-effect of acting upon a cancellation request while in a condition wait
// is that the mutex is (in effect) re-acquired before calling the first
// cancellation cleanup handler. The effect is as if the thread were unblocked,
// allowed to execute up to the point of returning from the call to
// pthread_cond_timedwait() or pthread_cond_wait(), but at that point notices
// the cancellation request and instead of returning to the caller of
// pthread_cond_timedwait() or pthread_cond_wait(), starts the thread
// cancellation activities, which includes calling cancellation cleanup
// handlers.
//
// A thread that has been unblocked because it has been canceled while blocked
// in a call to pthread_cond_timedwait() or pthread_cond_wait() shall not
// consume any condition signal that may be directed concurrently at the
// condition variable if there are other threads blocked on the condition
// variable.
//
// The pthread_cond_timedwait() function shall be equivalent to
// pthread_cond_wait(), except that an error is returned if the absolute time
// specified by abstime passes (that is, system time equals or exceeds abstime)
// before the condition cond is signaled or broadcasted, or if the absolute
// time specified by abstime has already been passed at the time of the call.
// When such timeouts occur, pthread_cond_timedwait() shall nonetheless release
// and re-acquire the mutex referenced by mutex, and may consume a condition
// signal directed concurrently at the condition variable.
//
// The condition variable shall have a clock attribute which specifies the
// clock that shall be used to measure the time specified by the abstime
// argument. The pthread_cond_timedwait() function is also a cancellation
// point.
//
// If a signal is delivered to a thread waiting for a condition variable, upon
// return from the signal handler the thread resumes waiting for the condition
// variable as if it was not interrupted, or it shall return zero due to
// spurious wakeup.
//
// The behavior is undefined if the value specified by the cond or mutex
// argument to these functions does not refer to an initialized condition
// variable or an initialized mutex object, respectively.
//
// Except for [ETIMEDOUT], [ENOTRECOVERABLE], and [EOWNERDEAD], all these error
// checks shall act as if they were performed immediately at the beginning of
// processing for the function and shall cause an error return, in effect,
// prior to modifying the state of the mutex specified by mutex or the
// condition variable specified by cond.
//
// Upon successful completion, a value of zero shall be returned; otherwise, an
// error number shall be returned to indicate the error.
//
// These functions shall fail if:
//
// [ENOTRECOVERABLE] The state protected by the mutex is not recoverable.
//
// [EOWNERDEAD] The mutex is a robust mutex and the process containing the
// previous owning thread terminated while holding the mutex lock. The mutex
// lock shall be acquired by the calling thread and it is up to the new owner
// to make the state consistent.
//
// [EPERM] The mutex type is PTHREAD_MUTEX_ERRORCHECK or the mutex is a robust
// mutex, and the current thread does not own the mutex.
//
// The pthread_cond_timedwait() function shall fail if:
//
// [ETIMEDOUT] The time specified by abstime to pthread_cond_timedwait() has
// passed.
//
// [EINVAL] The abstime argument specified a nanosecond value less than zero or
// greater than or equal to 1000 million.
//
// These functions may fail if:
//
// [EOWNERDEAD] The mutex is a robust mutex and the previous owning thread
// terminated while holding the mutex lock. The mutex lock shall be acquired by
// the calling thread and it is up to the new owner to make the state
// consistent.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_cond_timedwait(tls TLS, cond, mutex, abstime uintptr) (r int32) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_cond_timedwait.%v(%#x, %#x, %#x) ", tid, cond, mutex, abstime)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	ts := *(*Stimespec)(unsafe.Pointer(abstime))
	if ts.Xtv_nsec < 0 || ts.Xtv_nsec > 1e9 {
	}
	deadline := time.Unix(ts.Xtv_sec, ts.Xtv_nsec)
	if r = Xpthread_cond_wait(tls, cond, mutex); r != 0 {
		return r
	}

	if time.Now().After(deadline) {
		return errno.XETIMEDOUT
	}

	return 0
}

// Documentation: see Xpthread_cond_timedwait.
func Xpthread_cond_wait(tls TLS, cond, mutex uintptr) (r int32) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_cond_wait.%v(%#x, %#x) ", tid, cond, mutex)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	mu := mutexes.mutex(mutex)
	mu.Unlock()
	c := conds.cond(cond)
	sc := sync.NewCond(&locker{m: mu, tid: tid})
	c.conds[sc] = struct{}{}
	for c.signals == 0 {
		c.Unlock()
		sc.Wait()
		c.Lock()
	}
	c.signals--
	delete(c.conds, sc)
	c.Unlock()
	return 0
}

//  C: int pthread_cond_destroy(pthread_cond_t *cond);
//  C: int pthread_cond_init(pthread_cond_t *restrict cond, const pthread_condattr_t *restrict attr);
//
// The pthread_cond_destroy() function shall destroy the given condition
// variable specified by cond; the object becomes, in effect, uninitialized. An
// implementation may cause pthread_cond_destroy() to set the object referenced
// by cond to an invalid value. A destroyed condition variable object can be
// reinitialized using pthread_cond_init(); the results of otherwise
// referencing the object after it has been destroyed are undefined.
//
// It shall be safe to destroy an initialized condition variable upon which no
// threads are currently blocked. Attempting to destroy a condition variable
// upon which other threads are currently blocked results in undefined
// behavior.
//
// The pthread_cond_init() function shall initialize the condition variable
// referenced by cond with attributes referenced by attr. If attr is NULL, the
// default condition variable attributes shall be used; the effect is the same
// as passing the address of a default condition variable attributes object.
// Upon successful initialization, the state of the condition variable shall
// become initialized.
//
// See Synchronization Object Copies and Alternative Mappings for further
// requirements.
//
// Attempting to initialize an already initialized condition variable results
// in undefined behavior.
//
// In cases where default condition variable attributes are appropriate, the
// macro PTHREAD_COND_INITIALIZER can be used to initialize condition
// variables. The effect shall be equivalent to dynamic initialization by a
// call to pthread_cond_init() with parameter attr specified as NULL, except
// that no error checks are performed.
//
// The behavior is undefined if the value specified by the cond argument to
// pthread_cond_destroy() does not refer to an initialized condition variable.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_cond_init() does not refer to an initialized condition variable
// attributes object.
//
// If successful, the pthread_cond_destroy() and pthread_cond_init() functions
// shall return zero; otherwise, an error number shall be returned to indicate
// the error.
//
// The pthread_cond_init() function shall fail if:
//
// [EAGAIN] The system lacked the necessary resources (other than memory) to
// initialize another condition variable.
//
// [ENOMEM] Insufficient memory exists to initialize the condition variable.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_cond_destroy(tls TLS, cond uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_cond_destroy.%v(%#x) ", tid, cond)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	conds.Lock()
	delete(conds.m, cond)
	conds.Unlock()
	return 0
}

// Documentation: see Xpthread_cond_destroy.
func Xpthread_cond_init(tls TLS, cond, attr uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_cond_init.%v(%#x, %#x) ", tid, cond, attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	if attr != 0 {
		panic("TODO")
	}
	conds.cond(cond).Unlock()
	return 0
}

//  C: int pthread_cond_broadcast(pthread_cond_t *cond);
//  C: int pthread_cond_signal(pthread_cond_t *cond);
//
// These functions shall unblock threads blocked on a condition variable.
//
// The pthread_cond_broadcast() function shall unblock all threads currently
// blocked on the specified condition variable cond.
//
// The pthread_cond_signal() function shall unblock at least one of the threads
// that are blocked on the specified condition variable cond (if any threads
// are blocked on cond).
//
// If more than one thread is blocked on a condition variable, the scheduling
// policy shall determine the order in which threads are unblocked. When each
// thread unblocked as a result of a pthread_cond_broadcast() or
// pthread_cond_signal() returns from its call to pthread_cond_wait() or
// pthread_cond_timedwait(), the thread shall own the mutex with which it
// called pthread_cond_wait() or pthread_cond_timedwait(). The thread(s) that
// are unblocked shall contend for the mutex according to the scheduling policy
// (if applicable), and as if each had called pthread_mutex_lock().
//
// The pthread_cond_broadcast() or pthread_cond_signal() functions may be
// called by a thread whether or not it currently owns the mutex that threads
// calling pthread_cond_wait() or pthread_cond_timedwait() have associated with
// the condition variable during their waits; however, if predictable
// scheduling behavior is required, then that mutex shall be locked by the
// thread calling pthread_cond_broadcast() or pthread_cond_signal().
//
// The pthread_cond_broadcast() and pthread_cond_signal() functions shall have
// no effect if there are no threads currently blocked on cond.
//
// The behavior is undefined if the value specified by the cond argument to
// pthread_cond_broadcast() or pthread_cond_signal() does not refer to an
// initialized condition variable.
//
// If successful, the pthread_cond_broadcast() and pthread_cond_signal()
// functions shall return zero; otherwise, an error number shall be returned to
// indicate the error.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_cond_broadcast(tls TLS, cond uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_cond_broadcast.%v(%#x) ", tid, cond)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	c := conds.cond(cond)
	for k := range c.conds {
		c.signals++
		defer k.Signal()
	}
	c.Unlock()
	return 0
}

// Documentation: see Xpthread_cond_broadcast.
func Xpthread_cond_signal(tls TLS, cond uintptr) int32 {
	panic("TODO")
}

//  C. int pthread_key_create(pthread_key_t *key, void (*destructor)(void*));
//
// The pthread_key_create() function shall create a thread-specific data key
// visible to all threads in the process. Key values provided by
// pthread_key_create() are opaque objects used to locate thread-specific data.
// Although the same key value may be used by different threads, the values
// bound to the key by pthread_setspecific() are maintained on a per-thread
// basis and persist for the life of the calling thread.
//
// Upon key creation, the value NULL shall be associated with the new key in
// all active threads. Upon thread creation, the value NULL shall be associated
// with all defined keys in the new thread.
//
// An optional destructor function may be associated with each key value. At
// thread exit, if a key value has a non-NULL destructor pointer, and the
// thread has a non-NULL value associated with that key, the value of the key
// is set to NULL, and then the function pointed to is called with the
// previously associated value as its sole argument. The order of destructor
// calls is unspecified if more than one destructor exists for a thread when it
// exits.
//
// If, after all the destructors have been called for all non-NULL values with
// associated destructors, there are still some non-NULL values with associated
// destructors, then the process is repeated. If, after at least
// {PTHREAD_DESTRUCTOR_ITERATIONS} iterations of destructor calls for
// outstanding non-NULL values, there are still some non-NULL values with
// associated destructors, implementations may stop calling destructors, or
// they may continue calling destructors until no non-NULL values with
// associated destructors exist, even though this might result in an infinite
// loop.
//
// If successful, the pthread_key_create() function shall store the newly
// created key value at *key and shall return zero. Otherwise, an error number
// shall be returned to indicate the error.
//
// The pthread_key_create() function shall fail if:
//
// [EAGAIN] The system lacked the necessary resources to create another
// thread-specific data key, or the system-imposed limit on the total number of
// keys per process {PTHREAD_KEYS_MAX} has been exceeded.
//
// [ENOMEM] Insufficient memory exists to create the key.
//
// The pthread_key_create() function shall not return an error code of [EINTR].
func Xpthread_key_create(tls TLS, key, destructor uintptr) (r int32) {
	var k pthread_key_t
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_key_create(%v, %#x, %#x) ", tid, key, destructor)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v (key %#x)\n", r, k)
		}()
	}
	pthreadDBMu.Lock()
	pthreadDBNextKey++
	k = pthreadDBNextKey
	if destructor != 0 {
		pthreadDBDestructors[k] = destructor //TODO actually call destructrors
	}
	*(*pthread_key_t)(unsafe.Pointer(key)) = k
	pthreadDBMu.Unlock()
	return 0

}

//  C: void pthread_exit(void *value_ptr);
//
// The pthread_exit() function shall terminate the calling thread and make the
// value value_ptr available to any successful join with the terminating
// thread. Any cancellation cleanup handlers that have been pushed and not yet
// popped shall be popped in the reverse order that they were pushed and then
// executed. After all cancellation cleanup handlers have been executed, if the
// thread has any thread-specific data, appropriate destructor functions shall
// be called in an unspecified order. Thread termination does not release any
// application visible process resources, including, but not limited to,
// mutexes and file descriptors, nor does it perform any process-level cleanup
// actions, including, but not limited to, calling any atexit() routines that
// may exist.
//
// An implicit call to pthread_exit() is made when a thread other than the
// thread in which main() was first invoked returns from the start routine that
// was used to create it. The function's return value shall serve as the
// thread's exit status.
//
// The behavior of pthread_exit() is undefined if called from a cancellation
// cleanup handler or destructor function that was invoked as a result of
// either an implicit or explicit call to pthread_exit().
//
// After a thread has terminated, the result of access to local (auto)
// variables of the thread is undefined. Thus, references to local variables of
// the exiting thread should not be used for the pthread_exit() value_ptr
// parameter value.
//
// The process shall exit with an exit status of 0 after the last thread has
// been terminated. The behavior shall be as if the implementation called
// exit() with a zero argument at thread termination time.
//
// The pthread_exit() function cannot return to its caller.
//
// No errors are defined.
func Xpthread_exit(tls TLS, value_ptr uintptr) {
	if ptrace {
		tid := pthread_t(tls.getThreadID())
		fmt.Fprintf(TraceWriter, "pthread_exit.%v(%#x) ", tid, value_ptr)
	}
	panic(threadExited(value_ptr))
}

//  C: int pthread_attr_destroy(pthread_attr_t *attr);
//  C: int pthread_attr_init(pthread_attr_t *attr);
//
// The pthread_attr_destroy() function shall destroy a thread attributes
// object. An implementation may cause pthread_attr_destroy() to set attr to an
// implementation-defined invalid value. A destroyed attr attributes object can
// be reinitialized using pthread_attr_init(); the results of otherwise
// referencing the object after it has been destroyed are undefined.
//
// The pthread_attr_init() function shall initialize a thread attributes object
// attr with the default value for all of the individual attributes used by a
// given implementation.
//
// The resulting attributes object (possibly modified by setting individual
// attribute values) when used by pthread_create() defines the attributes of
// the thread created. A single attributes object can be used in multiple
// simultaneous calls to pthread_create(). Results are undefined if
// pthread_attr_init() is called specifying an already initialized attr
// attributes object.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_attr_destroy() does not refer to an initialized thread attributes
// object.
//
// Upon successful completion, pthread_attr_destroy() and pthread_attr_init()
// shall return a value of 0; otherwise, an error number shall be returned to
// indicate the error.
//
// The pthread_attr_init() function shall fail if:
//
// [ENOMEM] Insufficient memory exists to initialize the thread attributes
// object.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_attr_destroy(tls TLS, attr uintptr) (r int32) {
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_attr_destroy.%v(%#x) ", tid, attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	attrs.Lock()
	delete(attrs.m, attr)
	attrs.Unlock()
	return 0
}

// Documentation: see Xpthread_attr_destroy.
func Xpthread_attr_init(tls TLS, attr uintptr) (r int32) {
	attrs.attr(attr)
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_attr_init.%v(%#x) ", tid, attr)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	return 0
}

//  C: int pthread_attr_getscope(const pthread_attr_t *restrict attr, int *restrict contentionscope);
//  C: int pthread_attr_setscope(pthread_attr_t *attr, int contentionscope);
//
// The pthread_attr_getscope() and pthread_attr_setscope() functions,
// respectively, shall get and set the contentionscope attribute in the attr
// object.
//
// The contentionscope attribute may have the values PTHREAD_SCOPE_SYSTEM,
// signifying system scheduling contention scope, or PTHREAD_SCOPE_PROCESS,
// signifying process scheduling contention scope. The symbols
// PTHREAD_SCOPE_SYSTEM and PTHREAD_SCOPE_PROCESS are defined in the
// <pthread.h> header.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_attr_getscope() or pthread_attr_setscope() does not refer to an
// initialized thread attributes object.
//
// If successful, the pthread_attr_getscope() and pthread_attr_setscope()
// functions shall return zero; otherwise, an error number shall be returned to
// indicate the error.
//
// The pthread_attr_setscope() function shall fail if:
//
// [ENOTSUP] An attempt was made to set the attribute to an unsupported value.
//
// The pthread_attr_setscope() function may fail if:
//
// [EINVAL] The value of contentionscope is not valid.
//
// These functions shall not return an error code of [EINTR].
func Xpthread_attr_getscope(tls TLS, attr, contentionscope uintptr) int32 {
	panic("TODO")
}

// Documentatio: see pthread_attr_getscope.
func Xpthread_attr_setscope(tls TLS, attr uintptr, contentionscope int32) (r int32) {
	a := attrs.attr(attr)
	if ptrace {
		tid := tls.getThreadID()
		fmt.Fprintf(TraceWriter, "pthread_attr_setscope.%v(%#x, %v) ", tid, a, contentionscope)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	switch contentionscope {
	case
		pthread.CPTHREAD_SCOPE_PROCESS,
		pthread.CPTHREAD_SCOPE_SYSTEM:
		a.scope = contentionscope
	default:
		r = errno.XEINVAL
	}
	return r
}

//  C: int pthread_attr_getdetachstate(const pthread_attr_t *attr, int *detachstate);
//  C: int pthread_attr_setdetachstate(pthread_attr_t *attr, int detachstate);
//
// The detachstate attribute controls whether the thread is created in a
// detached state. If the thread is created detached, then use of the ID of the
// newly created thread by the pthread_detach() or pthread_join() function is
// an error.
//
// The pthread_attr_getdetachstate() and pthread_attr_setdetachstate()
// functions, respectively, shall get and set the detachstate attribute in the
// attr object.
//
// For pthread_attr_getdetachstate(), detachstate shall be set to either
// PTHREAD_CREATE_DETACHED or PTHREAD_CREATE_JOINABLE.
//
// For pthread_attr_setdetachstate(), the application shall set detachstate to
// either PTHREAD_CREATE_DETACHED or PTHREAD_CREATE_JOINABLE.
//
// A value of PTHREAD_CREATE_DETACHED shall cause all threads created with attr
// to be in the detached state, whereas using a value of
// PTHREAD_CREATE_JOINABLE shall cause all threads created with attr to be in
// the joinable state. The default value of the detachstate attribute shall be
// PTHREAD_CREATE_JOINABLE.
//
// The behavior is undefined if the value specified by the attr argument to
// pthread_attr_getdetachstate() or pthread_attr_setdetachstate() does not
// refer to an initialized thread attributes object.
//
// Upon successful completion, pthread_attr_getdetachstate() and
// pthread_attr_setdetachstate() shall return a value of 0; otherwise, an error
// number shall be returned to indicate the error.
//
// The pthread_attr_getdetachstate() function stores the value of the
// detachstate attribute in detachstate if successful.
//
// The pthread_attr_setdetachstate() function shall fail if:
//
// [EINVAL] The value of detachstate was not valid
//
// These functions shall not return an error code of [EINTR].
func Xpthread_attr_getdetachstate(tls TLS, attr, detachstate uintptr) int32 {
	panic("TODO")
}

// Documentation: see Xpthread_attr_getdetachstate.
func Xpthread_attr_setdetachstate(tls TLS, attr uintptr, detachstate int32) int32 {
	panic("TODO")
}

//  C: void *pthread_getspecific(pthread_key_t key);
//  C: int pthread_setspecific(pthread_key_t key, const void *value);
//
// The pthread_getspecific() function shall return the value currently bound to
// the specified key on behalf of the calling thread.
//
// The pthread_setspecific() function shall associate a thread-specific value
// with a key obtained via a previous call to pthread_key_create(). Different
// threads may bind different values to the same key. These values are
// typically pointers to blocks of dynamically allocated memory that have been
// reserved for use by the calling thread.
//
// The effect of calling pthread_getspecific() or pthread_setspecific() with a
// key value not obtained from pthread_key_create() or after key has been
// deleted with pthread_key_delete() is undefined.
//
// Both pthread_getspecific() and pthread_setspecific() may be called from a
// thread-specific data destructor function. A call to pthread_getspecific()
// for the thread-specific data key being destroyed shall return the value
// NULL, unless the value is changed (after the destructor starts) by a call to
// pthread_setspecific(). Calling pthread_setspecific() from a thread-specific
// data destructor routine may result either in lost storage (after at least
// PTHREAD_DESTRUCTOR_ITERATIONS attempts at destruction) or in an infinite
// loop.
//
// Both functions may be implemented as macros.
//
// The pthread_getspecific() function shall return the thread-specific data
// value associated with the given key. If no thread-specific data value is
// associated with key, then the value NULL shall be returned.
//
// If successful, the pthread_setspecific() function shall return zero;
// otherwise, an error number shall be returned to indicate the error.
//
// No errors are returned from pthread_getspecific().
//
// The pthread_setspecific() function shall fail if:
//
// [ENOMEM] Insufficient memory exists to associate the non-NULL value with the
// key.
//
// The pthread_setspecific() function shall not return an error code of
// [EINTR].
func Xpthread_getspecific(tls TLS, key pthread_key_t) (r uintptr) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_getspecific.%v(%#x) ", tid, key)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	pthreadDBMu.Lock()
	r = pthreadDB[pthreadDBKey{thread: tid, key: key}]
	pthreadDBMu.Unlock()
	return r
}

// Documentation: see Xpthread_getspecific.
func Xpthread_setspecific(tls TLS, key pthread_key_t, value uintptr) (r int32) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_setspecific.%v(%#x, %#x) ", tid, key, value)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	pthreadDBMu.Lock()
	pthreadDB[pthreadDBKey{thread: tls.getThreadID(), key: key}] = value
	pthreadDBMu.Unlock()
	return 0
}

//  C: int pthread_key_delete(pthread_key_t key);
//
// The pthread_key_delete() function shall delete a thread-specific data key
// previously returned by pthread_key_create(). The thread-specific data values
// associated with key need not be NULL at the time pthread_key_delete() is
// called. It is the responsibility of the application to free any application
// storage or perform any cleanup actions for data structures related to the
// deleted key or associated thread-specific data in any threads; this cleanup
// can be done either before or after pthread_key_delete() is called. Any
// attempt to use key following the call to pthread_key_delete() results in
// undefined behavior.
//
// The pthread_key_delete() function shall be callable from within destructor
// functions. No destructor functions shall be invoked by pthread_key_delete().
// Any destructor function that may have been associated with key shall no
// longer be called upon thread exit.
//
// If successful, the pthread_key_delete() function shall return zero;
// otherwise, an error number shall be returned to indicate the error.
//
// The pthread_key_delete() function shall not return an error code of [EINTR].

//  C: int pthread_key_delete(pthread_key_t key);
func Xpthread_key_delete(tls TLS, key pthread_key_t) (r int32) {
	tid := tls.getThreadID()
	if ptrace {
		fmt.Fprintf(TraceWriter, "pthread_key_delete.%v(%#x) ", tid, key)
		defer func() {
			fmt.Fprintf(TraceWriter, "%v\n", r)
		}()
	}
	pthreadDBMu.Lock()
	delete(pthreadDB, pthreadDBKey{thread: tid, key: key})
	pthreadDBMu.Unlock()
	return 0
}

//  C: int pthread_detach(pthread_t thread);
//
// The pthread_detach() function shall indicate to the implementation that
// storage for the thread thread can be reclaimed when that thread terminates.
// If thread has not terminated, pthread_detach() shall not cause it to
// terminate.
//
// The behavior is undefined if the value specified by the thread argument to
// pthread_detach() does not refer to a joinable thread.
//
// If the call succeeds, pthread_detach() shall return 0; otherwise, an error
// number shall be returned to indicate the error.
//
// The pthread_detach() function shall not return an error code of [EINTR].
func Xpthread_detach(tls TLS, thread pthread_t) int32 {
	panic("TODO")
}

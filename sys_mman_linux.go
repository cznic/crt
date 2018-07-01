// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"syscall"
)

// void *mmap(void *addr, size_t len, int prot, int flags, int fildes, off_t off);
//
// The mmap() function shall establish a mapping between an address space of a
// process and a memory object.
//
// The mmap() function shall be supported for the following memory objects:
//
// Regular files
//
// [SHM] [Option Start] Shared memory objects [Option End]
//
// [TYM] [Option Start] Typed memory objects [Option End]
//
// Support for any other type of file is unspecified.
//
// The format of the call is as follows:
//
// pa=mmap(addr, len, prot, flags, fildes, off);
//
// The mmap() function shall establish a mapping between the address space of
// the process at an address pa for len bytes to the memory object represented
// by the file descriptor fildes at offset off for len bytes. The value of pa
// is an implementation-defined function of the parameter addr and the values
// of flags, further described below. A successful mmap() call shall return pa
// as its result. The address range starting at pa and continuing for len bytes
// shall be legitimate for the possible (not necessarily current) address space
// of the process. The range of bytes starting at off and continuing for len
// bytes shall be legitimate for the possible (not necessarily current) offsets
// in the memory object represented by fildes.
//
// [TYM] [Option Start] If fildes represents a typed memory object opened with
// either the POSIX_TYPED_MEM_ALLOCATE flag or the
// POSIX_TYPED_MEM_ALLOCATE_CONTIG flag, the memory object to be mapped shall
// be that portion of the typed memory object allocated by the implementation
// as specified below. In this case, if off is non-zero, the behavior of mmap()
// is undefined. If fildes refers to a valid typed memory object that is not
// accessible from the calling process, mmap() shall fail. [Option End]
//
// The mapping established by mmap() shall replace any previous mappings for
// those whole pages containing any part of the address space of the process
// starting at pa and continuing for len bytes.
//
// If the size of the mapped file changes after the call to mmap() as a result
// of some other operation on the mapped file, the effect of references to
// portions of the mapped region that correspond to added or removed portions
// of the file is unspecified.
//
// If len is zero, mmap() shall fail and no mapping shall be established.
//
// The parameter prot determines whether read, write, execute, or some
// combination of accesses are permitted to the data being mapped. The prot
// shall be either PROT_NONE or the bitwise-inclusive OR of one or more of the
// other flags in the following table, defined in the <sys/mman.h> header.
//
//  Symbolic Constant	Description
//  PROT_READ		Data can be read.
//  PROT_WRITE		Data can be written.
//  PROT_EXEC		Data can be executed.
//  PROT_NONE		Data cannot be accessed.
//
// If an implementation cannot support the combination of access types
// specified by prot, the call to mmap() shall fail.
//
// An implementation may permit accesses other than those specified by prot;
// however, the implementation shall not permit a write to succeed where
// PROT_WRITE has not been set and shall not permit any access where PROT_NONE
// alone has been set. The implementation shall support at least the following
// values of prot: PROT_NONE, PROT_READ, PROT_WRITE, and the bitwise-inclusive
// OR of PROT_READ and PROT_WRITE. The file descriptor fildes shall have been
// opened with read permission, regardless of the protection options specified.
// If PROT_WRITE is specified, the application shall ensure that it has opened
// the file descriptor fildes with write permission unless MAP_PRIVATE is
// specified in the flags parameter as described below.
//
// The parameter flags provides other information about the handling of the
// mapped data. The value of flags is the bitwise-inclusive OR of these
// options, defined in <sys/mman.h>:
//
//  Symbolic Constant	Description
//  MAP_SHARED		Changes are shared.
//  MAP_PRIVATE		Changes are private.
//  MAP_FIXED		Interpret addr exactly.
//
// It is implementation-defined whether MAP_FIXED shall be supported. [XSI]
// [Option Start]  MAP_FIXED shall be supported on XSI-conformant systems.
// [Option End]
//
// MAP_SHARED and MAP_PRIVATE describe the disposition of write references to
// the memory object. If MAP_SHARED is specified, write references shall change
// the underlying object. If MAP_PRIVATE is specified, modifications to the
// mapped data by the calling process shall be visible only to the calling
// process and shall not change the underlying object. It is unspecified
// whether modifications to the underlying object done after the MAP_PRIVATE
// mapping is established are visible through the MAP_PRIVATE mapping. Either
// MAP_SHARED or MAP_PRIVATE can be specified, but not both. The mapping type
// is retained across fork().
//
// The state of synchronization objects such as mutexes, semaphores, barriers,
// and conditional variables placed in shared memory mapped with MAP_SHARED
// becomes undefined when the last region in any process containing the
// synchronization object is unmapped.
//
// [TYM] [Option Start] When fildes represents a typed memory object opened
// with either the POSIX_TYPED_MEM_ALLOCATE flag or the
// POSIX_TYPED_MEM_ALLOCATE_CONTIG flag, mmap() shall, if there are enough
// resources available, map len bytes allocated from the corresponding typed
// memory object which were not previously allocated to any process in any
// processor that may access that typed memory object. If there are not enough
// resources available, the function shall fail. If fildes represents a typed
// memory object opened with the POSIX_TYPED_MEM_ALLOCATE_CONTIG flag, these
// allocated bytes shall be contiguous within the typed memory object. If
// fildes represents a typed memory object opened with the
// POSIX_TYPED_MEM_ALLOCATE flag, these allocated bytes may be composed of
// non-contiguous fragments within the typed memory object. If fildes
// represents a typed memory object opened with neither the
// POSIX_TYPED_MEM_ALLOCATE_CONTIG flag nor the POSIX_TYPED_MEM_ALLOCATE flag,
// len bytes starting at offset off within the typed memory object are mapped,
// exactly as when mapping a file or shared memory object. In this case, if two
// processes map an area of typed memory using the same off and len values and
// using file descriptors that refer to the same memory pool (either from the
// same port or from a different port), both processes shall map the same
// region of storage. [Option End]
//
// When MAP_FIXED is set in the flags argument, the implementation is informed
// that the value of pa shall be addr, exactly. If MAP_FIXED is set, mmap() may
// return MAP_FAILED and set errno to [EINVAL]. If a MAP_FIXED request is
// successful, then any previous mappings [ML|MLR] [Option Start]  or memory
// locks [Option End]  for those whole pages containing any part of the address
// range [pa,pa+len) shall be removed, as if by an appropriate call to
// munmap(), before the new mapping is established.
//
// When MAP_FIXED is not set, the implementation uses addr in an
// implementation-defined manner to arrive at pa. The pa so chosen shall be an
// area of the address space that the implementation deems suitable for a
// mapping of len bytes to the file. All implementations interpret an addr
// value of 0 as granting the implementation complete freedom in selecting pa,
// subject to constraints described below. A non-zero value of addr is taken to
// be a suggestion of a process address near which the mapping should be
// placed. When the implementation selects a value for pa, it never places a
// mapping at address 0, nor does it replace any extant mapping.
//
// If MAP_FIXED is specified and addr is non-zero, it shall have the same
// remainder as the off parameter, modulo the page size as returned by
// sysconf() when passed _SC_PAGESIZE or _SC_PAGE_SIZE. The implementation may
// require that off is a multiple of the page size. If MAP_FIXED is specified,
// the implementation may require that addr is a multiple of the page size. The
// system performs mapping operations over whole pages. Thus, while the
// parameter len need not meet a size or alignment constraint, the system shall
// include, in any mapping operation, any partial page specified by the address
// range starting at pa and continuing for len bytes.
//
// The system shall always zero-fill any partial page at the end of an object.
// Further, the system shall never write out any modified portions of the last
// page of an object which are beyond its end. References within the address
// range starting at pa and continuing for len bytes to whole pages following
// the end of an object shall result in delivery of a SIGBUS signal.
//
// An implementation may generate SIGBUS signals when a reference would cause
// an error in the mapped object, such as out-of-space condition.
//
// The mmap() function shall add an extra reference to the file associated with
// the file descriptor fildes which is not removed by a subsequent close() on
// that file descriptor. This reference shall be removed when there are no more
// mappings to the file.
//
// The last data access timestamp of the mapped file may be marked for update
// at any time between the mmap() call and the corresponding munmap() call. The
// initial read or write reference to a mapped region shall cause the file's
// last data access timestamp to be marked for update if it has not already
// been marked for update.
//
// The last data modification and last file status change timestamps of a file
// that is mapped with MAP_SHARED and PROT_WRITE shall be marked for update at
// some point in the interval between a write reference to the mapped region
// and the next call to msync() with MS_ASYNC or MS_SYNC for that portion of
// the file by any process. If there is no such call and if the underlying file
// is modified as a result of a write reference, then these timestamps shall be
// marked for update at some time after the write reference.
//
// There may be implementation-defined limits on the number of memory regions
// that can be mapped (per process or per system).
//
// [XSI] [Option Start] If such a limit is imposed, whether the number of
// memory regions that can be mapped by a process is decreased by the use of
// shmat() is implementation-defined. [Option End]
//
// If mmap() fails for reasons other than [EBADF], [EINVAL], or [ENOTSUP], some
// of the mappings in the address range starting at addr and continuing for len
// bytes may have been unmapped.
//
// Upon successful completion, the mmap() function shall return the address at
// which the mapping was placed (pa); otherwise, it shall return a value of
// MAP_FAILED and set errno to indicate the error. The symbol MAP_FAILED is
// defined in the <sys/mman.h> header. No successful return from mmap() shall
// return the value MAP_FAILED.
func Xmmap(tls TLS, addr uintptr, len size_t, prot, flags, fildes int32, off int64) uintptr {
	r, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, uintptr(len), uintptr(prot), uintptr(flags), uintptr(fildes), uintptr(off))
	if strace {
		fmt.Fprintf(os.Stderr, "mmap(%#x, %#x, %#x, %#x, %#x, %#x) (%#x, %v)\n", addr, len, prot, flags, fildes, off, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return r
}

func Xmmap64(tls TLS, addr uintptr, len size_t, prot, flags, fildes int32, off int64) uintptr {
	return Xmmap(tls, addr, len, prot, flags, fildes, off)
}

// int munmap(void *addr, size_t len);
func Xmunmap(tls TLS, addr uintptr, len size_t) int32 {
	r, _, err := syscall.Syscall(syscall.SYS_MUNMAP, addr, uintptr(len), 0)
	if strace {
		fmt.Fprintf(os.Stderr, "munmap(%#x, %#x) (%#x, %v)\n", addr, len, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return int32(r)
}

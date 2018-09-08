#include "pthread_impl.h"
#include "atomic.h"
#include "syscall.h"
/* cheat and reuse CRTJMP macro from dynlink code */
#include "dynlink.h"
#include <assert.h>

static volatile int lock;
static void *unmap_base;
static size_t unmap_size;
static char shared_stack[256];

static void do_unmap()
{
	__syscall(SYS_munmap, unmap_base, unmap_size);
	__syscall(SYS_exit);
}

void __unmapself(void *base, size_t size)
{
	int tid=__pthread_self()->tid;
	char *stack = shared_stack + sizeof shared_stack;
	stack -= (uintptr_t)stack % 16;
	while (lock || a_cas(&lock, 0, tid))
		a_spin();
	__syscall(SYS_set_tid_address, &lock);
	unmap_base = base;
	unmap_size = size;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		CRTJMP(do_unmap, stack);
}

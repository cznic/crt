#include <sys/mman.h>
#include "syscall.h"
#include <assert.h>

int msync(void *start, size_t len, int flags)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_msync, start, len, flags);
}

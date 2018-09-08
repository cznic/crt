#include <sys/random.h>
#include "syscall.h"
#include <assert.h>

ssize_t getrandom(void *buf, size_t buflen, unsigned flags)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_getrandom, buf, buflen, flags);
}

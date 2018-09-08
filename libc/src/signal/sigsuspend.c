#include <signal.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int sigsuspend(const sigset_t *mask)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_rt_sigsuspend, mask, _NSIG/8);
}

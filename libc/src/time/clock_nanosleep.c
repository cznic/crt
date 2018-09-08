#include <time.h>
#include <errno.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int clock_nanosleep(clockid_t clk, int flags, const struct timespec *req, struct timespec *rem)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		int r = -__syscall_cp(SYS_clock_nanosleep, clk, flags, req, rem);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return clk == CLOCK_THREAD_CPUTIME_ID ? EINVAL : r;
}

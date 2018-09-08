#include <signal.h>
#include <errno.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int sigtimedwait(const sigset_t *restrict mask, siginfo_t *restrict si, const struct timespec *restrict timeout)
{
	int ret;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		do ret = syscall_cp(SYS_rt_sigtimedwait, mask,
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			si, timeout, _NSIG/8);
	while (ret<0 && errno==EINTR);
	return ret;
}

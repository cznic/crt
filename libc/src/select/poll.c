#include <poll.h>
#include <time.h>
#include <signal.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int poll(struct pollfd *fds, nfds_t n, int timeout)
{
#ifdef SYS_poll
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_poll, fds, n, timeout);
#else
	return syscall_cp(SYS_ppoll, fds, n, timeout>=0 ?
		&((struct timespec){ .tv_sec = timeout/1000,
		.tv_nsec = timeout%1000*1000000 }) : 0, 0, _NSIG/8);
#endif
}

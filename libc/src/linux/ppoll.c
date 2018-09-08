#define _GNU_SOURCE
#include <poll.h>
#include <signal.h>
#include "syscall.h"
#include <assert.h>

int ppoll(struct pollfd *fds, nfds_t n, const struct timespec *to, const sigset_t *mask)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_ppoll, fds, n,
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			to ? (struct timespec []){*to} : 0, mask, _NSIG/8);
}

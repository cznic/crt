#define _GNU_SOURCE
#include <sys/socket.h>
#include <errno.h>
#include <fcntl.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int accept4(int fd, struct sockaddr *restrict addr, socklen_t *restrict len, int flg)
{
	if (!flg) return accept(fd, addr, len);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		int ret = socketcall_cp(accept4, fd, addr, len, flg, 0, 0);
//TODO(ccgo)		if (ret>=0 || (errno != ENOSYS && errno != EINVAL)) return ret;
//TODO(ccgo)		ret = accept(fd, addr, len);
//TODO(ccgo)		if (ret<0) return ret;
//TODO(ccgo)		if (flg & SOCK_CLOEXEC)
//TODO(ccgo)			__syscall(SYS_fcntl, ret, F_SETFD, FD_CLOEXEC);
//TODO(ccgo)		if (flg & SOCK_NONBLOCK)
//TODO(ccgo)			__syscall(SYS_fcntl, ret, F_SETFL, O_NONBLOCK);
//TODO(ccgo)		return ret;
}

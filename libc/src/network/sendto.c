#include <sys/socket.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

ssize_t sendto(int fd, const void *buf, size_t len, int flags, const struct sockaddr *addr, socklen_t alen)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return socketcall_cp(sendto, fd, buf, len, flags, addr, alen);
}

#include <sys/socket.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

ssize_t recvfrom(int fd, void *restrict buf, size_t len, int flags, struct sockaddr *restrict addr, socklen_t *restrict alen)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return socketcall_cp(recvfrom, fd, buf, len, flags, addr, alen);
}

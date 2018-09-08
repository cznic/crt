#include <sys/socket.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int connect(int fd, const struct sockaddr *addr, socklen_t len)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return socketcall_cp(connect, fd, addr, len, 0, 0, 0);
}

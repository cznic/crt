#include <sys/socket.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int accept(int fd, struct sockaddr *restrict addr, socklen_t *restrict len)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return socketcall_cp(accept, fd, addr, len, 0, 0, 0);
}

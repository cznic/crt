#include <sys/uio.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

ssize_t readv(int fd, const struct iovec *iov, int count)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_readv, fd, iov, count);
}

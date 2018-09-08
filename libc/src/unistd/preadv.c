#define _BSD_SOURCE
#include <sys/uio.h>
#include <unistd.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

ssize_t preadv(int fd, const struct iovec *iov, int count, off_t ofs)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_preadv, fd, iov, count,
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			(long)(ofs), (long)(ofs>>32));
}

LFS64(preadv);

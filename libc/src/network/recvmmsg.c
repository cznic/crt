#define _GNU_SOURCE
#include <sys/socket.h>
#include <limits.h>
#include "syscall.h"
#include <assert.h>

int recvmmsg(int fd, struct mmsghdr *msgvec, unsigned int vlen, unsigned int flags, struct timespec *timeout)
{
#if LONG_MAX > INT_MAX
	struct mmsghdr *mh = msgvec;
	unsigned int i;
	for (i = vlen; i; i--, mh++)
		mh->msg_hdr.__pad1 = mh->msg_hdr.__pad2 = 0;
#endif
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_recvmmsg, fd, msgvec, vlen, flags, timeout);
}

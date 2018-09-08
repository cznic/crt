#include <mqueue.h>
#include "syscall.h"
#include <assert.h>

ssize_t mq_timedreceive(mqd_t mqd, char *restrict msg, size_t len, unsigned *restrict prio, const struct timespec *restrict at)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_mq_timedreceive, mqd, msg, len, prio, at);
}

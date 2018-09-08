#include <mqueue.h>
#include "syscall.h"
#include <assert.h>

int mq_timedsend(mqd_t mqd, const char *msg, size_t len, unsigned prio, const struct timespec *at)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_mq_timedsend, mqd, msg, len, prio, at);
}

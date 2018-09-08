#include <unistd.h>
#include <sys/time.h>
#include "syscall.h"
#include <assert.h>

unsigned alarm(unsigned seconds)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		struct itimerval it = { .it_value.tv_sec = seconds };
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		__syscall(SYS_setitimer, ITIMER_REAL, &it, &it);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return it.it_value.tv_sec + !!it.it_value.tv_usec;
}

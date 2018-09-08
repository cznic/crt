#define _GNU_SOURCE
#include <unistd.h>
#include <sys/time.h>
#include <assert.h>

unsigned ualarm(unsigned value, unsigned interval)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		struct itimerval it = {
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			.it_interval.tv_usec = interval,
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			.it_value.tv_usec = value
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		};
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		setitimer(ITIMER_REAL, &it, &it);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return it.it_value.tv_sec*1000000 + it.it_value.tv_usec;
}

#define _GNU_SOURCE
#include <fcntl.h>
#include <string.h>
#include <unistd.h>
#include <sys/prctl.h>
#include <assert.h>

#include "pthread_impl.h"

int pthread_setname_np(pthread_t thread, const char *name)
{
	int fd, cs, status = 0;
	char f[sizeof "/proc/self/task//comm" + 3*sizeof(int)];
	size_t len;

	if ((len = strnlen(name, 16)) > 15) return ERANGE;

	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if (thread == pthread_self())
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			return prctl(PR_SET_NAME, (unsigned long)name, 0UL, 0UL, 0UL) ? errno : 0;

	snprintf(f, sizeof f, "/proc/self/task/%d/comm", thread->tid);
	pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, &cs);
	if ((fd = open(f, O_WRONLY)) < 0 || write(fd, name, len) < 0) status = errno;
	if (fd >= 0) close(fd);
	pthread_setcancelstate(cs, 0);
	return status;
}

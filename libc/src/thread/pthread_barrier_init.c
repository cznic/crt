#include "pthread_impl.h"
#include <assert.h>

int pthread_barrier_init(pthread_barrier_t *restrict b, const pthread_barrierattr_t *restrict a, unsigned count)
{
	if (count-1 > INT_MAX-1) return EINVAL;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		*b = (pthread_barrier_t){ ._b_limit = count-1 | (a?a->__attr:0) };
	return 0;
}

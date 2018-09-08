#include "pthread_impl.h"
#include <string.h>

extern size_t __default_stacksize;
extern size_t __default_guardsize;

int pthread_setattr_default_np(const pthread_attr_t *attrp)
{
	/* Reject anything in the attr object other than stack/guard size. */
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		pthread_attr_t tmp = *attrp, zero = { 0 };
//TODO(ccgo)		tmp._a_stacksize = 0;
//TODO(ccgo)		tmp._a_guardsize = 0;
//TODO(ccgo)		if (memcmp(&tmp, &zero, sizeof tmp))
//TODO(ccgo)			return EINVAL;
//TODO(ccgo)	
//TODO(ccgo)		__inhibit_ptc();
//TODO(ccgo)		if (attrp->_a_stacksize >= __default_stacksize)
//TODO(ccgo)			__default_stacksize = attrp->_a_stacksize;
//TODO(ccgo)		if (attrp->_a_guardsize >= __default_guardsize)
//TODO(ccgo)			__default_guardsize = attrp->_a_guardsize;
//TODO(ccgo)		__release_ptc();
//TODO(ccgo)	
//TODO(ccgo)		return 0;
}

int pthread_getattr_default_np(pthread_attr_t *attrp)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		__acquire_ptc();
//TODO(ccgo)		*attrp = (pthread_attr_t) {
//TODO(ccgo)			._a_stacksize = __default_stacksize,
//TODO(ccgo)			._a_guardsize = __default_guardsize,
//TODO(ccgo)		};
//TODO(ccgo)		__release_ptc();
//TODO(ccgo)		return 0;
}

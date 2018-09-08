#include <stddef.h>
#include "pthread_impl.h"
#include "libc.h"

__attribute__((__visibility__("hidden")))
void *__tls_get_new(tls_mod_off_t *);

void *__tls_get_addr(tls_mod_off_t *v)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		pthread_t self = __pthread_self();
//TODO(ccgo)		if (v[0]<=(size_t)self->dtv[0])
//TODO(ccgo)			return (char *)self->dtv[v[0]]+v[1]+DTP_OFFSET;
//TODO(ccgo)		return __tls_get_new(v);
}

weak_alias(__tls_get_addr, __tls_get_new);

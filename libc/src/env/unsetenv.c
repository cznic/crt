#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include "libc.h"
#include <assert.h>

char *__strchrnul(const char *, int);

static void dummy(char *old, char *new) {}
weak_alias(dummy, __env_rm_add);

int unsetenv(const char *name)
{
	size_t l = __strchrnul(name, '=') - name;
	if (!l || name[l]) {
		errno = EINVAL;
		return -1;
	}
	if (__environ) {
		char **e = __environ, **eo = e;
		for (; *e; e++)
			__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)				if (!strncmp(name, *e, l) && l[*e] == '=')
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)					__env_rm_add(*e, 0);
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)				else if (eo != e)
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)					*eo++ = *e;
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)				else
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)					eo++;
		if (eo != e) *eo = 0;
	}
	return 0;
}

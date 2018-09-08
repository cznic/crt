#define _GNU_SOURCE
#include "stdio_impl.h"
#include <string.h>
#include <assert.h>

char *fgetln(FILE *f, size_t *plen)
{
	char *ret = 0, *z;
	ssize_t l;
	FLOCK(f);
	ungetc(getc_unlocked(f), f);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if ((z=memchr(f->rpos, '\n', f->rend - f->rpos))) {
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			ret = (char *)f->rpos;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			*plen = ++z - ret;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			f->rpos = (void *)z;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		} else if ((l = getline(&f->getln_buf, (size_t[]){0}, f)) > 0) {
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			*plen = l;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			ret = f->getln_buf;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		}
	FUNLOCK(f);
	return ret;
}

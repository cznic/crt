#include "stdio_impl.h"
#include "locale_impl.h"
#include <assert.h>

int fwide(FILE *f, int mode)
{
	FLOCK(f);
	if (mode) {
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)			if (!f->locale) f->locale = MB_CUR_MAX==1
		__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)				? C_LOCALE : UTF8_LOCALE;
		if (!f->mode) f->mode = mode>0 ? 1 : -1;
	}
	mode = f->mode;
	FUNLOCK(f);
	return mode;
}

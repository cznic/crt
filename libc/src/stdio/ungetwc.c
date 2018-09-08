#include "stdio_impl.h"
#include "locale_impl.h"
#include <wchar.h>
#include <limits.h>
#include <ctype.h>
#include <string.h>
#include <assert.h>

wint_t ungetwc(wint_t c, FILE *f)
{
	unsigned char mbc[MB_LEN_MAX];
	int l;
	locale_t *ploc = &CURRENT_LOCALE, loc = *ploc;

	FLOCK(f);

	if (f->mode <= 0) fwide(f, 1);
	*ploc = f->locale;

	if (!f->rpos) __toread(f);
	if (!f->rpos || c == WEOF || (l = wcrtomb((void *)mbc, c, 0)) < 0 ||
	    f->rpos < f->buf - UNGET + l) {
		FUNLOCK(f);
		*ploc = loc;
		return WEOF;
	}

	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if (isascii(c)) *--f->rpos = c;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		else memcpy(f->rpos -= l, mbc, l);

	f->flags &= ~F_EOF;

	FUNLOCK(f);
	*ploc = loc;
	return c;
}

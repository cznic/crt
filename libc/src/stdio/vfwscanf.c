#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <ctype.h>
#include <wchar.h>
#include <wctype.h>
#include <limits.h>
#include <string.h>
#include <assert.h>

#include "stdio_impl.h"
#include "shgetc.h"
#include "intscan.h"
#include "floatscan.h"
#include "libc.h"

#define SIZE_hh -2
#define SIZE_h  -1
#define SIZE_def 0
#define SIZE_l   1
#define SIZE_L   2
#define SIZE_ll  3

static void store_int(void *dest, int size, unsigned long long i)
{
	if (!dest) return;
	switch (size) {
	case SIZE_hh:
		*(char *)dest = i;
		break;
	case SIZE_h:
		*(short *)dest = i;
		break;
	case SIZE_def:
		*(int *)dest = i;
		break;
	case SIZE_l:
		*(long *)dest = i;
		break;
	case SIZE_ll:
		*(long long *)dest = i;
		break;
	}
}

static void *arg_n(va_list ap, unsigned int n)
{
	void *p;
	unsigned int i;
	va_list ap2;
	va_copy(ap2, ap);
	for (i=n; i>1; i--) va_arg(ap2, void *);
	p = va_arg(ap2, void *);
	va_end(ap2);
	return p;
}

static int in_set(const wchar_t *set, int c)
{
	int j;
	const wchar_t *p = set;
	if (*p == '-') {
		if (c=='-') return 1;
		p++;
	} else if (*p == ']') {
		if (c==']') return 1;
		p++;
	}
	for (; *p && *p != ']'; p++) {
		if (*p=='-' && p[1] && p[1] != ']')
			for (j=p++[-1]; j<*p; j++)
				if (c==j) return 1;
		if (c==*p) return 1;
	}
	return 0;
}

#if 1
#undef getwc
#define getwc(f) \
	((f)->rpos < (f)->rend && *(f)->rpos < 128 ? *(f)->rpos++ : (getwc)(f))

#undef ungetwc
#define ungetwc(c,f) \
	((f)->rend && (c)<128U ? *--(f)->rpos : ungetwc((c),(f)))
#endif

int vfwscanf(FILE *restrict f, const wchar_t *restrict fmt, va_list ap)
{
	int width;
	int size;
	int alloc;
	const wchar_t *p;
	int c, t;
	char *s;
	wchar_t *wcs;
	void *dest=NULL;
	int invert;
	int matches=0;
	off_t pos = 0, cnt;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		static const char size_pfx[][3] = { "hh", "h", "", "l", "L", "ll" };
//TODO(ccgo)		char tmp[3*sizeof(int)+10];
//TODO(ccgo)		const wchar_t *set;
//TODO(ccgo)		size_t i, k;
//TODO(ccgo)	
//TODO(ccgo)		FLOCK(f);
//TODO(ccgo)	
//TODO(ccgo)		fwide(f, 1);
//TODO(ccgo)	
//TODO(ccgo)		for (p=fmt; *p; p++) {
//TODO(ccgo)	
//TODO(ccgo)			alloc = 0;
//TODO(ccgo)	
//TODO(ccgo)			if (iswspace(*p)) {
//TODO(ccgo)				while (iswspace(p[1])) p++;
//TODO(ccgo)				while (iswspace((c=getwc(f)))) pos++;
//TODO(ccgo)				ungetwc(c, f);
//TODO(ccgo)				continue;
//TODO(ccgo)			}
//TODO(ccgo)			if (*p != '%' || p[1] == '%') {
//TODO(ccgo)				if (*p == '%') {
//TODO(ccgo)					p++;
//TODO(ccgo)					while (iswspace((c=getwc(f)))) pos++;
//TODO(ccgo)				} else {
//TODO(ccgo)					c = getwc(f);
//TODO(ccgo)				}
//TODO(ccgo)				if (c!=*p) {
//TODO(ccgo)					ungetwc(c, f);
//TODO(ccgo)					if (c<0) goto input_fail;
//TODO(ccgo)					goto match_fail;
//TODO(ccgo)				}
//TODO(ccgo)				pos++;
//TODO(ccgo)				continue;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			p++;
//TODO(ccgo)			if (*p=='*') {
//TODO(ccgo)				dest = 0; p++;
//TODO(ccgo)			} else if (iswdigit(*p) && p[1]=='$') {
//TODO(ccgo)				dest = arg_n(ap, *p-'0'); p+=2;
//TODO(ccgo)			} else {
//TODO(ccgo)				dest = va_arg(ap, void *);
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			for (width=0; iswdigit(*p); p++) {
//TODO(ccgo)				width = 10*width + *p - '0';
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			if (*p=='m') {
//TODO(ccgo)				wcs = 0;
//TODO(ccgo)				s = 0;
//TODO(ccgo)				alloc = !!dest;
//TODO(ccgo)				p++;
//TODO(ccgo)			} else {
//TODO(ccgo)				alloc = 0;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			size = SIZE_def;
//TODO(ccgo)			switch (*p++) {
//TODO(ccgo)			case 'h':
//TODO(ccgo)				if (*p == 'h') p++, size = SIZE_hh;
//TODO(ccgo)				else size = SIZE_h;
//TODO(ccgo)				break;
//TODO(ccgo)			case 'l':
//TODO(ccgo)				if (*p == 'l') p++, size = SIZE_ll;
//TODO(ccgo)				else size = SIZE_l;
//TODO(ccgo)				break;
//TODO(ccgo)			case 'j':
//TODO(ccgo)				size = SIZE_ll;
//TODO(ccgo)				break;
//TODO(ccgo)			case 'z':
//TODO(ccgo)			case 't':
//TODO(ccgo)				size = SIZE_l;
//TODO(ccgo)				break;
//TODO(ccgo)			case 'L':
//TODO(ccgo)				size = SIZE_L;
//TODO(ccgo)				break;
//TODO(ccgo)			case 'd': case 'i': case 'o': case 'u': case 'x':
//TODO(ccgo)			case 'a': case 'e': case 'f': case 'g':
//TODO(ccgo)			case 'A': case 'E': case 'F': case 'G': case 'X':
//TODO(ccgo)			case 's': case 'c': case '[':
//TODO(ccgo)			case 'S': case 'C':
//TODO(ccgo)			case 'p': case 'n':
//TODO(ccgo)				p--;
//TODO(ccgo)				break;
//TODO(ccgo)			default:
//TODO(ccgo)				goto fmt_fail;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			t = *p;
//TODO(ccgo)	
//TODO(ccgo)			/* Transform S,C -> ls,lc */
//TODO(ccgo)			if ((t&0x2f)==3) {
//TODO(ccgo)				size = SIZE_l;
//TODO(ccgo)				t |= 32;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			if (t != 'n') {
//TODO(ccgo)				if (t != '[' && (t|32) != 'c')
//TODO(ccgo)					while (iswspace((c=getwc(f)))) pos++;
//TODO(ccgo)				else
//TODO(ccgo)					c=getwc(f);
//TODO(ccgo)				if (c < 0) goto input_fail;
//TODO(ccgo)				ungetwc(c, f);
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			switch (t) {
//TODO(ccgo)			case 'n':
//TODO(ccgo)				store_int(dest, size, pos);
//TODO(ccgo)				/* do not increment match count, etc! */
//TODO(ccgo)				continue;
//TODO(ccgo)	
//TODO(ccgo)			case 's':
//TODO(ccgo)			case 'c':
//TODO(ccgo)			case '[':
//TODO(ccgo)				if (t == 'c') {
//TODO(ccgo)					if (width<1) width = 1;
//TODO(ccgo)					invert = 1;
//TODO(ccgo)					set = L"";
//TODO(ccgo)				} else if (t == 's') {
//TODO(ccgo)					invert = 1;
//TODO(ccgo)					static const wchar_t spaces[] = {
//TODO(ccgo)						' ', '\t', '\n', '\r', 11, 12,  0x0085,
//TODO(ccgo)						0x2000, 0x2001, 0x2002, 0x2003, 0x2004, 0x2005,
//TODO(ccgo)						0x2006, 0x2008, 0x2009, 0x200a,
//TODO(ccgo)						0x2028, 0x2029, 0x205f, 0x3000, 0 };
//TODO(ccgo)					set = spaces;
//TODO(ccgo)				} else {
//TODO(ccgo)					if (*++p == '^') p++, invert = 1;
//TODO(ccgo)					else invert = 0;
//TODO(ccgo)					set = p;
//TODO(ccgo)					if (*p==']') p++;
//TODO(ccgo)					while (*p!=']') {
//TODO(ccgo)						if (!*p) goto fmt_fail;
//TODO(ccgo)						p++;
//TODO(ccgo)					}
//TODO(ccgo)				}
//TODO(ccgo)	
//TODO(ccgo)				s = (size == SIZE_def) ? dest : 0;
//TODO(ccgo)				wcs = (size == SIZE_l) ? dest : 0;
//TODO(ccgo)	
//TODO(ccgo)				int gotmatch = 0;
//TODO(ccgo)	
//TODO(ccgo)				if (width < 1) width = -1;
//TODO(ccgo)	
//TODO(ccgo)				i = 0;
//TODO(ccgo)				if (alloc) {
//TODO(ccgo)					k = t=='c' ? width+1U : 31;
//TODO(ccgo)					if (size == SIZE_l) {
//TODO(ccgo)						wcs = malloc(k*sizeof(wchar_t));
//TODO(ccgo)						if (!wcs) goto alloc_fail;
//TODO(ccgo)					} else {
//TODO(ccgo)						s = malloc(k);
//TODO(ccgo)						if (!s) goto alloc_fail;
//TODO(ccgo)					}
//TODO(ccgo)				}
//TODO(ccgo)				while (width) {
//TODO(ccgo)					if ((c=getwc(f))<0) break;
//TODO(ccgo)					if (in_set(set, c) == invert)
//TODO(ccgo)						break;
//TODO(ccgo)					if (wcs) {
//TODO(ccgo)						wcs[i++] = c;
//TODO(ccgo)						if (alloc && i==k) {
//TODO(ccgo)							k += k+1;
//TODO(ccgo)							wchar_t *tmp = realloc(wcs, k*sizeof(wchar_t));
//TODO(ccgo)							if (!tmp) goto alloc_fail;
//TODO(ccgo)							wcs = tmp;
//TODO(ccgo)						}
//TODO(ccgo)					} else if (size != SIZE_l) {
//TODO(ccgo)						int l = wctomb(s?s+i:tmp, c);
//TODO(ccgo)						if (l<0) goto input_fail;
//TODO(ccgo)						i += l;
//TODO(ccgo)						if (alloc && i > k-4) {
//TODO(ccgo)							k += k+1;
//TODO(ccgo)							char *tmp = realloc(s, k);
//TODO(ccgo)							if (!tmp) goto alloc_fail;
//TODO(ccgo)							s = tmp;
//TODO(ccgo)						}
//TODO(ccgo)					}
//TODO(ccgo)					pos++;
//TODO(ccgo)					width-=(width>0);
//TODO(ccgo)					gotmatch=1;
//TODO(ccgo)				}
//TODO(ccgo)				if (width) {
//TODO(ccgo)					ungetwc(c, f);
//TODO(ccgo)					if (t == 'c' || !gotmatch) goto match_fail;
//TODO(ccgo)				}
//TODO(ccgo)	
//TODO(ccgo)				if (alloc) {
//TODO(ccgo)					if (size == SIZE_l) *(wchar_t **)dest = wcs;
//TODO(ccgo)					else *(char **)dest = s;
//TODO(ccgo)				}
//TODO(ccgo)				if (t != 'c') {
//TODO(ccgo)					if (wcs) wcs[i] = 0;
//TODO(ccgo)					if (s) s[i] = 0;
//TODO(ccgo)				}
//TODO(ccgo)				break;
//TODO(ccgo)	
//TODO(ccgo)			case 'd': case 'i': case 'o': case 'u': case 'x':
//TODO(ccgo)			case 'a': case 'e': case 'f': case 'g':
//TODO(ccgo)			case 'A': case 'E': case 'F': case 'G': case 'X':
//TODO(ccgo)			case 'p':
//TODO(ccgo)				if (width < 1) width = 0;
//TODO(ccgo)				snprintf(tmp, sizeof tmp, "%.*s%.0d%s%c%%lln",
//TODO(ccgo)					1+!dest, "%*", width, size_pfx[size+2], t);
//TODO(ccgo)				cnt = 0;
//TODO(ccgo)				if (fscanf(f, tmp, dest?dest:&cnt, &cnt) == -1)
//TODO(ccgo)					goto input_fail;
//TODO(ccgo)				else if (!cnt)
//TODO(ccgo)					goto match_fail;
//TODO(ccgo)				pos += cnt;
//TODO(ccgo)				break;
//TODO(ccgo)			default:
//TODO(ccgo)				goto fmt_fail;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			if (dest) matches++;
//TODO(ccgo)		}
//TODO(ccgo)		if (0) {
//TODO(ccgo)	fmt_fail:
//TODO(ccgo)	alloc_fail:
//TODO(ccgo)	input_fail:
//TODO(ccgo)			if (!matches) matches--;
//TODO(ccgo)	match_fail:
//TODO(ccgo)			if (alloc) {
//TODO(ccgo)				free(s);
//TODO(ccgo)				free(wcs);
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)		FUNLOCK(f);
//TODO(ccgo)		return matches;
}

weak_alias(vfwscanf,__isoc99_vfwscanf);

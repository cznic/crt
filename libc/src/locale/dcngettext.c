#include <libintl.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include <sys/stat.h>
#include <ctype.h>
#include "locale_impl.h"
#include "libc.h"
#include "atomic.h"
#include <assert.h>

struct binding {
	struct binding *next;
	int dirlen;
	volatile int active;
	char *domainname;
	char *dirname;
	char buf[];
};

static void *volatile bindings;

static char *gettextdir(const char *domainname, size_t *dirlen)
{
	struct binding *p;
	for (p=bindings; p; p=p->next) {
		if (!strcmp(p->domainname, domainname) && p->active) {
			*dirlen = p->dirlen;
			return (char *)p->dirname;
		}
	}
	return 0;
}

char *bindtextdomain(const char *domainname, const char *dirname)
{
	static volatile int lock[1];
	struct binding *p, *q;

	if (!domainname) return 0;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if (!dirname) return gettextdir(domainname, &(size_t){0});

	size_t domlen = strnlen(domainname, NAME_MAX+1);
	size_t dirlen = strnlen(dirname, PATH_MAX);
	if (domlen > NAME_MAX || dirlen >= PATH_MAX) {
		errno = EINVAL;
		return 0;
	}

	LOCK(lock);

	for (p=bindings; p; p=p->next) {
		if (!strcmp(p->domainname, domainname) &&
		    !strcmp(p->dirname, dirname)) {
			break;
		}
	}

	if (!p) {
		p = calloc(sizeof *p + domlen + dirlen + 2, 1);
		if (!p) {
			UNLOCK(lock);
			return 0;
		}
		p->next = bindings;
		p->dirlen = dirlen;
		p->domainname = p->buf;
		p->dirname = p->buf + domlen + 1;
		memcpy(p->domainname, domainname, domlen+1);
		memcpy(p->dirname, dirname, dirlen+1);
		a_cas_p(&bindings, bindings, p);
	}

	a_store(&p->active, 1);

	for (q=bindings; q; q=q->next) {
		if (!strcmp(q->domainname, domainname) && q != p)
			a_store(&q->active, 0);
	}

	UNLOCK(lock);
	
	return (char *)p->dirname;
}

//TODO(ccgo)	static const char catnames[][12] = {
//TODO(ccgo)		"LC_CTYPE",
//TODO(ccgo)		"LC_NUMERIC",
//TODO(ccgo)		"LC_TIME",
//TODO(ccgo)		"LC_COLLATE",
//TODO(ccgo)		"LC_MONETARY",
//TODO(ccgo)		"LC_MESSAGES",
//TODO(ccgo)	};

static const char catlens[] = { 8, 10, 7, 10, 11, 11 };

struct msgcat {
	struct msgcat *next;
	const void *map;
	size_t map_size;
	void *volatile plural_rule;
	volatile int nplurals;
	struct binding *binding;
	const struct __locale_map *lm;
	int cat;
};

static char *dummy_gettextdomain()
{
	return "messages";
}

weak_alias(dummy_gettextdomain, __gettextdomain);

const unsigned char *__map_file(const char *, size_t *);
int __munmap(void *, size_t);
unsigned long __pleval(const char *, unsigned long);

char *dcngettext(const char *domainname, const char *msgid1, const char *msgid2, unsigned long int n, int category)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		static struct msgcat *volatile cats;
//TODO(ccgo)		struct msgcat *p;
//TODO(ccgo)		struct __locale_struct *loc = CURRENT_LOCALE;
//TODO(ccgo)		const struct __locale_map *lm;
//TODO(ccgo)		size_t domlen;
//TODO(ccgo)		struct binding *q;
//TODO(ccgo)	
//TODO(ccgo)		if ((unsigned)category >= LC_ALL) goto notrans;
//TODO(ccgo)	
//TODO(ccgo)		if (!domainname) domainname = __gettextdomain();
//TODO(ccgo)	
//TODO(ccgo)		domlen = strnlen(domainname, NAME_MAX+1);
//TODO(ccgo)		if (domlen > NAME_MAX) goto notrans;
//TODO(ccgo)	
//TODO(ccgo)		for (q=bindings; q; q=q->next)
//TODO(ccgo)			if (!strcmp(q->domainname, domainname) && q->active)
//TODO(ccgo)				break;
//TODO(ccgo)		if (!q) goto notrans;
//TODO(ccgo)	
//TODO(ccgo)		lm = loc->cat[category];
//TODO(ccgo)		if (!lm) {
//TODO(ccgo)	notrans:
//TODO(ccgo)			return (char *) ((n == 1) ? msgid1 : msgid2);
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		for (p=cats; p; p=p->next)
//TODO(ccgo)			if (p->binding == q && p->lm == lm && p->cat == category)
//TODO(ccgo)				break;
//TODO(ccgo)	
//TODO(ccgo)		if (!p) {
//TODO(ccgo)			const char *dirname, *locname, *catname, *modname, *locp;
//TODO(ccgo)			size_t dirlen, loclen, catlen, modlen, alt_modlen;
//TODO(ccgo)			void *old_cats;
//TODO(ccgo)			size_t map_size;
//TODO(ccgo)	
//TODO(ccgo)			dirname = q->dirname;
//TODO(ccgo)			locname = lm->name;
//TODO(ccgo)			__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)	//TODO(ccgo)			catname = catnames[category];
//TODO(ccgo)	
//TODO(ccgo)			dirlen = q->dirlen;
//TODO(ccgo)			loclen = strlen(locname);
//TODO(ccgo)			catlen = catlens[category];
//TODO(ccgo)	
//TODO(ccgo)			/* Logically split @mod suffix from locale name. */
//TODO(ccgo)			modname = memchr(locname, '@', loclen);
//TODO(ccgo)			if (!modname) modname = locname + loclen;
//TODO(ccgo)			alt_modlen = modlen = loclen - (modname-locname);
//TODO(ccgo)			loclen = modname-locname;
//TODO(ccgo)	
//TODO(ccgo)			/* Drop .charset identifier; it is not used. */
//TODO(ccgo)			const char *csp = memchr(locname, '.', loclen);
//TODO(ccgo)			if (csp) loclen = csp-locname;
//TODO(ccgo)	
//TODO(ccgo)			char name[dirlen+1 + loclen+modlen+1 + catlen+1 + domlen+3 + 1];
//TODO(ccgo)			const void *map;
//TODO(ccgo)	
//TODO(ccgo)			for (;;) {
//TODO(ccgo)				__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)	//TODO(ccgo)				snprintf(name, sizeof name, "%s/%.*s%.*s/%s/%s.mo\0",
//TODO(ccgo)				__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)	//TODO(ccgo)					dirname, (int)loclen, locname,
//TODO(ccgo)				__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)	//TODO(ccgo)					(int)alt_modlen, modname, catname, domainname);
//TODO(ccgo)				if (map = __map_file(name, &map_size)) break;
//TODO(ccgo)	
//TODO(ccgo)				/* Try dropping @mod, _YY, then both. */
//TODO(ccgo)				if (alt_modlen) {
//TODO(ccgo)					alt_modlen = 0;
//TODO(ccgo)				} else if ((locp = memchr(locname, '_', loclen))) {
//TODO(ccgo)					loclen = locp-locname;
//TODO(ccgo)					alt_modlen = modlen;
//TODO(ccgo)				} else {
//TODO(ccgo)					break;
//TODO(ccgo)				}
//TODO(ccgo)			}
//TODO(ccgo)			if (!map) goto notrans;
//TODO(ccgo)	
//TODO(ccgo)			p = calloc(sizeof *p, 1);
//TODO(ccgo)			if (!p) {
//TODO(ccgo)				__munmap((void *)map, map_size);
//TODO(ccgo)				goto notrans;
//TODO(ccgo)			}
//TODO(ccgo)			p->cat = category;
//TODO(ccgo)			p->binding = q;
//TODO(ccgo)			p->lm = lm;
//TODO(ccgo)			p->map = map;
//TODO(ccgo)			p->map_size = map_size;
//TODO(ccgo)			do {
//TODO(ccgo)				old_cats = cats;
//TODO(ccgo)				p->next = old_cats;
//TODO(ccgo)			} while (a_cas_p(&cats, old_cats, p) != old_cats);
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		const char *trans = __mo_lookup(p->map, p->map_size, msgid1);
//TODO(ccgo)		if (!trans) goto notrans;
//TODO(ccgo)	
//TODO(ccgo)		/* Non-plural-processing gettext forms pass a null pointer as
//TODO(ccgo)		 * msgid2 to request that dcngettext suppress plural processing. */
//TODO(ccgo)		if (!msgid2) return (char *)trans;
//TODO(ccgo)	
//TODO(ccgo)		if (!p->plural_rule) {
//TODO(ccgo)			const char *rule = "n!=1;";
//TODO(ccgo)			unsigned long np = 2;
//TODO(ccgo)			const char *r = __mo_lookup(p->map, p->map_size, "");
//TODO(ccgo)			char *z;
//TODO(ccgo)			while (r && strncmp(r, "Plural-Forms:", 13)) {
//TODO(ccgo)				z = strchr(r, '\n');
//TODO(ccgo)				r = z ? z+1 : 0;
//TODO(ccgo)			}
//TODO(ccgo)			if (r) {
//TODO(ccgo)				r += 13;
//TODO(ccgo)				while (isspace(*r)) r++;
//TODO(ccgo)				if (!strncmp(r, "nplurals=", 9)) {
//TODO(ccgo)					np = strtoul(r+9, &z, 10);
//TODO(ccgo)					r = z;
//TODO(ccgo)				}
//TODO(ccgo)				while (*r && *r != ';') r++;
//TODO(ccgo)				if (*r) {
//TODO(ccgo)					r++;
//TODO(ccgo)					while (isspace(*r)) r++;
//TODO(ccgo)					if (!strncmp(r, "plural=", 7))
//TODO(ccgo)						rule = r+7;
//TODO(ccgo)				}
//TODO(ccgo)			}
//TODO(ccgo)			a_store(&p->nplurals, np);
//TODO(ccgo)			a_cas_p(&p->plural_rule, 0, (void *)rule);
//TODO(ccgo)		}
//TODO(ccgo)		if (p->nplurals) {
//TODO(ccgo)			unsigned long plural = __pleval(p->plural_rule, n);
//TODO(ccgo)			if (plural > p->nplurals) goto notrans;
//TODO(ccgo)			while (plural--) {
//TODO(ccgo)				size_t rem = p->map_size - (trans - (char *)p->map);
//TODO(ccgo)				size_t l = strnlen(trans, rem);
//TODO(ccgo)				if (l+1 >= rem)
//TODO(ccgo)					goto notrans;
//TODO(ccgo)				trans += l+1;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)		return (char *)trans;
}

char *dcgettext(const char *domainname, const char *msgid, int category)
{
	return dcngettext(domainname, msgid, 0, 1, category);
}

char *dngettext(const char *domainname, const char *msgid1, const char *msgid2, unsigned long int n)
{
	return dcngettext(domainname, msgid1, msgid2, n, LC_MESSAGES);
}

char *dgettext(const char *domainname, const char *msgid)
{
	return dcngettext(domainname, msgid, 0, 1, LC_MESSAGES);
}

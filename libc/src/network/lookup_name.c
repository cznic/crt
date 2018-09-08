#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <net/if.h>
#include <arpa/inet.h>
#include <ctype.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <pthread.h>
#include <errno.h>
#include "lookup.h"
#include "stdio_impl.h"
#include "syscall.h"

static int is_valid_hostname(const char *host)
{
	const unsigned char *s;
	if (strnlen(host, 255)-1 >= 254 || mbstowcs(0, host, 0) == -1) return 0;
	for (s=(void *)host; *s>=0x80 || *s=='.' || *s=='-' || isalnum(*s); s++);
	return !*s;
}

//TODO(ccgo)	static int name_from_null(struct address buf[static 2], const char *name, int family, int flags)
//TODO(ccgo)	{
//TODO(ccgo)		int cnt = 0;
//TODO(ccgo)		if (name) return 0;
//TODO(ccgo)		if (flags & AI_PASSIVE) {
//TODO(ccgo)			if (family != AF_INET6)
//TODO(ccgo)				buf[cnt++] = (struct address){ .family = AF_INET };
//TODO(ccgo)			if (family != AF_INET)
//TODO(ccgo)				buf[cnt++] = (struct address){ .family = AF_INET6 };
//TODO(ccgo)		} else {
//TODO(ccgo)			if (family != AF_INET6)
//TODO(ccgo)				buf[cnt++] = (struct address){ .family = AF_INET, .addr = { 127,0,0,1 } };
//TODO(ccgo)			if (family != AF_INET)
//TODO(ccgo)				buf[cnt++] = (struct address){ .family = AF_INET6, .addr = { [15] = 1 } };
//TODO(ccgo)		}
//TODO(ccgo)		return cnt;
//TODO(ccgo)	}

//TODO(ccgo)	static int name_from_numeric(struct address buf[static 1], const char *name, int family)
//TODO(ccgo)	{
//TODO(ccgo)		return __lookup_ipliteral(buf, name, family);
//TODO(ccgo)	}

//TODO(ccgo)	static int name_from_hosts(struct address buf[static MAXADDRS], char canon[static 256], const char *name, int family)
//TODO(ccgo)	{
//TODO(ccgo)		char line[512];
//TODO(ccgo)		size_t l = strlen(name);
//TODO(ccgo)		int cnt = 0, badfam = 0;
//TODO(ccgo)		unsigned char _buf[1032];
//TODO(ccgo)		FILE _f, *f = __fopen_rb_ca("/etc/hosts", &_f, _buf, sizeof _buf);
//TODO(ccgo)		if (!f) switch (errno) {
//TODO(ccgo)		case ENOENT:
//TODO(ccgo)		case ENOTDIR:
//TODO(ccgo)		case EACCES:
//TODO(ccgo)			return 0;
//TODO(ccgo)		default:
//TODO(ccgo)			return EAI_SYSTEM;
//TODO(ccgo)		}
//TODO(ccgo)		while (fgets(line, sizeof line, f) && cnt < MAXADDRS) {
//TODO(ccgo)			char *p, *z;
//TODO(ccgo)	
//TODO(ccgo)			if ((p=strchr(line, '#'))) *p++='\n', *p=0;
//TODO(ccgo)			for(p=line+1; (p=strstr(p, name)) &&
//TODO(ccgo)				(!isspace(p[-1]) || !isspace(p[l])); p++);
//TODO(ccgo)			if (!p) continue;
//TODO(ccgo)	
//TODO(ccgo)			/* Isolate IP address to parse */
//TODO(ccgo)			for (p=line; *p && !isspace(*p); p++);
//TODO(ccgo)			*p++ = 0;
//TODO(ccgo)			switch (name_from_numeric(buf+cnt, line, family)) {
//TODO(ccgo)			case 1:
//TODO(ccgo)				cnt++;
//TODO(ccgo)				break;
//TODO(ccgo)			case 0:
//TODO(ccgo)				continue;
//TODO(ccgo)			default:
//TODO(ccgo)				badfam = EAI_NONAME;
//TODO(ccgo)				continue;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			/* Extract first name as canonical name */
//TODO(ccgo)			for (; *p && isspace(*p); p++);
//TODO(ccgo)			for (z=p; *z && !isspace(*z); z++);
//TODO(ccgo)			*z = 0;
//TODO(ccgo)			if (is_valid_hostname(p)) memcpy(canon, p, z-p+1);
//TODO(ccgo)		}
//TODO(ccgo)		__fclose_ca(f);
//TODO(ccgo)		return cnt ? cnt : badfam;
//TODO(ccgo)	}

struct dpc_ctx {
	struct address *addrs;
	char *canon;
	int cnt;
};

int __dns_parse(const unsigned char *, int, int (*)(void *, int, const void *, int, const void *), void *);
int __dn_expand(const unsigned char *, const unsigned char *, const unsigned char *, char *, int);
int __res_mkquery(int, const char *, int, int, const unsigned char *, int, const unsigned char*, unsigned char *, int);
int __res_msend_rc(int, const unsigned char *const *, const int *, unsigned char *const *, int *, int, const struct resolvconf *);

#define RR_A 1
#define RR_CNAME 5
#define RR_AAAA 28

static int dns_parse_callback(void *c, int rr, const void *data, int len, const void *packet)
{
	char tmp[256];
	struct dpc_ctx *ctx = c;
	if (ctx->cnt >= MAXADDRS) return -1;
	switch (rr) {
	case RR_A:
		if (len != 4) return -1;
		ctx->addrs[ctx->cnt].family = AF_INET;
		ctx->addrs[ctx->cnt].scopeid = 0;
		memcpy(ctx->addrs[ctx->cnt++].addr, data, 4);
		break;
	case RR_AAAA:
		if (len != 16) return -1;
		ctx->addrs[ctx->cnt].family = AF_INET6;
		ctx->addrs[ctx->cnt].scopeid = 0;
		memcpy(ctx->addrs[ctx->cnt++].addr, data, 16);
		break;
	case RR_CNAME:
		if (__dn_expand(packet, (const unsigned char *)packet + 512,
		    data, tmp, sizeof tmp) > 0 && is_valid_hostname(tmp))
			strcpy(ctx->canon, tmp);
		break;
	}
	return 0;
}

//TODO(ccgo)	static int name_from_dns(struct address buf[static MAXADDRS], char canon[static 256], const char *name, int family, const struct resolvconf *conf)
//TODO(ccgo)	{
//TODO(ccgo)		unsigned char qbuf[2][280], abuf[2][512];
//TODO(ccgo)		const unsigned char *qp[2] = { qbuf[0], qbuf[1] };
//TODO(ccgo)		unsigned char *ap[2] = { abuf[0], abuf[1] };
//TODO(ccgo)		int qlens[2], alens[2];
//TODO(ccgo)		int i, nq = 0;
//TODO(ccgo)		struct dpc_ctx ctx = { .addrs = buf, .canon = canon };
//TODO(ccgo)		static const struct { int af; int rr; } afrr[2] = {
//TODO(ccgo)			{ .af = AF_INET6, .rr = RR_A },
//TODO(ccgo)			{ .af = AF_INET, .rr = RR_AAAA },
//TODO(ccgo)		};
//TODO(ccgo)	
//TODO(ccgo)		for (i=0; i<2; i++) {
//TODO(ccgo)			if (family != afrr[i].af) {
//TODO(ccgo)				qlens[nq] = __res_mkquery(0, name, 1, afrr[i].rr,
//TODO(ccgo)					0, 0, 0, qbuf[nq], sizeof *qbuf);
//TODO(ccgo)				if (qlens[nq] == -1)
//TODO(ccgo)					return EAI_NONAME;
//TODO(ccgo)				nq++;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (__res_msend_rc(nq, qp, qlens, ap, alens, sizeof *abuf, conf) < 0)
//TODO(ccgo)			return EAI_SYSTEM;
//TODO(ccgo)	
//TODO(ccgo)		for (i=0; i<nq; i++)
//TODO(ccgo)			__dns_parse(abuf[i], alens[i], dns_parse_callback, &ctx);
//TODO(ccgo)	
//TODO(ccgo)		if (ctx.cnt) return ctx.cnt;
//TODO(ccgo)		if (alens[0] < 4 || (abuf[0][3] & 15) == 2) return EAI_AGAIN;
//TODO(ccgo)		if ((abuf[0][3] & 15) == 0) return EAI_NONAME;
//TODO(ccgo)		if ((abuf[0][3] & 15) == 3) return 0;
//TODO(ccgo)		return EAI_FAIL;
//TODO(ccgo)	}

//TODO(ccgo)	static int name_from_dns_search(struct address buf[static MAXADDRS], char canon[static 256], const char *name, int family)
//TODO(ccgo)	{
//TODO(ccgo)		char search[256];
//TODO(ccgo)		struct resolvconf conf;
//TODO(ccgo)		size_t l, dots;
//TODO(ccgo)		char *p, *z;
//TODO(ccgo)	
//TODO(ccgo)		if (__get_resolv_conf(&conf, search, sizeof search) < 0) return -1;
//TODO(ccgo)	
//TODO(ccgo)		/* Count dots, suppress search when >=ndots or name ends in
//TODO(ccgo)		 * a dot, which is an explicit request for global scope. */
//TODO(ccgo)		for (dots=l=0; name[l]; l++) if (name[l]=='.') dots++;
//TODO(ccgo)		if (dots >= conf.ndots || name[l-1]=='.') *search = 0;
//TODO(ccgo)	
//TODO(ccgo)		/* Strip final dot for canon, fail if multiple trailing dots. */
//TODO(ccgo)		if (name[l-1]=='.') l--;
//TODO(ccgo)		if (!l || name[l-1]=='.') return EAI_NONAME;
//TODO(ccgo)	
//TODO(ccgo)		/* This can never happen; the caller already checked length. */
//TODO(ccgo)		if (l >= 256) return EAI_NONAME;
//TODO(ccgo)	
//TODO(ccgo)		/* Name with search domain appended is setup in canon[]. This both
//TODO(ccgo)		 * provides the desired default canonical name (if the requested
//TODO(ccgo)		 * name is not a CNAME record) and serves as a buffer for passing
//TODO(ccgo)		 * the full requested name to name_from_dns. */
//TODO(ccgo)		memcpy(canon, name, l);
//TODO(ccgo)		canon[l] = '.';
//TODO(ccgo)	
//TODO(ccgo)		for (p=search; *p; p=z) {
//TODO(ccgo)			for (; isspace(*p); p++);
//TODO(ccgo)			for (z=p; *z && !isspace(*z); z++);
//TODO(ccgo)			if (z==p) break;
//TODO(ccgo)			if (z-p < 256 - l - 1) {
//TODO(ccgo)				memcpy(canon+l+1, p, z-p);
//TODO(ccgo)				canon[z-p+1+l] = 0;
//TODO(ccgo)				int cnt = name_from_dns(buf, canon, canon, family, &conf);
//TODO(ccgo)				if (cnt) return cnt;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		canon[l] = 0;
//TODO(ccgo)		return name_from_dns(buf, canon, name, family, &conf);
//TODO(ccgo)	}

static const struct policy {
	unsigned char addr[16];
	unsigned char len, mask;
	unsigned char prec, label;
} defpolicy[] = {
	{ "\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\1", 15, 0xff, 50, 0 },
	{ "\0\0\0\0\0\0\0\0\0\0\xff\xff", 11, 0xff, 35, 4 },
	{ "\x20\2", 1, 0xff, 30, 2 },
	{ "\x20\1", 3, 0xff, 5, 5 },
	{ "\xfc", 0, 0xfe, 3, 13 },
#if 0
	/* These are deprecated and/or returned to the address
	 * pool, so despite the RFC, treating them as special
	 * is probably wrong. */
	{ "", 11, 0xff, 1, 3 },
	{ "\xfe\xc0", 1, 0xc0, 1, 11 },
	{ "\x3f\xfe", 1, 0xff, 1, 12 },
#endif
	/* Last rule must match all addresses to stop loop. */
	{ "", 0, 0, 40, 1 },
};

static const struct policy *policyof(const struct in6_addr *a)
{
	int i;
	for (i=0; ; i++) {
		if (memcmp(a->s6_addr, defpolicy[i].addr, defpolicy[i].len))
			continue;
		if ((a->s6_addr[defpolicy[i].len] & defpolicy[i].mask)
		    != defpolicy[i].addr[defpolicy[i].len])
			continue;
		return defpolicy+i;
	}
}

static int labelof(const struct in6_addr *a)
{
	return policyof(a)->label;
}

static int scopeof(const struct in6_addr *a)
{
	if (IN6_IS_ADDR_MULTICAST(a)) return a->s6_addr[1] & 15;
	if (IN6_IS_ADDR_LINKLOCAL(a)) return 2;
	if (IN6_IS_ADDR_LOOPBACK(a)) return 2;
	if (IN6_IS_ADDR_SITELOCAL(a)) return 5;
	return 14;
}

static int prefixmatch(const struct in6_addr *s, const struct in6_addr *d)
{
	/* FIXME: The common prefix length should be limited to no greater
	 * than the nominal length of the prefix portion of the source
	 * address. However the definition of the source prefix length is
	 * not clear and thus this limiting is not yet implemented. */
	unsigned i;
	for (i=0; i<128 && !((s->s6_addr[i/8]^d->s6_addr[i/8])&(128>>(i%8))); i++);
	return i;
}

#define DAS_USABLE              0x40000000
#define DAS_MATCHINGSCOPE       0x20000000
#define DAS_MATCHINGLABEL       0x10000000
#define DAS_PREC_SHIFT          20
#define DAS_SCOPE_SHIFT         16
#define DAS_PREFIX_SHIFT        8
#define DAS_ORDER_SHIFT         0

static int addrcmp(const void *_a, const void *_b)
{
	const struct address *a = _a, *b = _b;
	return b->sortkey - a->sortkey;
}

//TODO(ccgo)	int __lookup_name(struct address buf[static MAXADDRS], char canon[static 256], const char *name, int family, int flags)
//TODO(ccgo)	{
//TODO(ccgo)		int cnt = 0, i, j;
//TODO(ccgo)	
//TODO(ccgo)		*canon = 0;
//TODO(ccgo)		if (name) {
//TODO(ccgo)			/* reject empty name and check len so it fits into temp bufs */
//TODO(ccgo)			size_t l = strnlen(name, 255);
//TODO(ccgo)			if (l-1 >= 254)
//TODO(ccgo)				return EAI_NONAME;
//TODO(ccgo)			memcpy(canon, name, l+1);
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* Procedurally, a request for v6 addresses with the v4-mapped
//TODO(ccgo)		 * flag set is like a request for unspecified family, followed
//TODO(ccgo)		 * by filtering of the results. */
//TODO(ccgo)		if (flags & AI_V4MAPPED) {
//TODO(ccgo)			if (family == AF_INET6) family = AF_UNSPEC;
//TODO(ccgo)			else flags -= AI_V4MAPPED;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* Try each backend until there's at least one result. */
//TODO(ccgo)		cnt = name_from_null(buf, name, family, flags);
//TODO(ccgo)		if (!cnt) cnt = name_from_numeric(buf, name, family);
//TODO(ccgo)		if (!cnt && !(flags & AI_NUMERICHOST)) {
//TODO(ccgo)			cnt = name_from_hosts(buf, canon, name, family);
//TODO(ccgo)			if (!cnt) cnt = name_from_dns_search(buf, canon, name, family);
//TODO(ccgo)		}
//TODO(ccgo)		if (cnt<=0) return cnt ? cnt : EAI_NONAME;
//TODO(ccgo)	
//TODO(ccgo)		/* Filter/transform results for v4-mapped lookup, if requested. */
//TODO(ccgo)		if (flags & AI_V4MAPPED) {
//TODO(ccgo)			if (!(flags & AI_ALL)) {
//TODO(ccgo)				/* If any v6 results exist, remove v4 results. */
//TODO(ccgo)				for (i=0; i<cnt && buf[i].family != AF_INET6; i++);
//TODO(ccgo)				if (i<cnt) {
//TODO(ccgo)					for (j=0; i<cnt; i++) {
//TODO(ccgo)						if (buf[i].family == AF_INET6)
//TODO(ccgo)							buf[j++] = buf[i];
//TODO(ccgo)					}
//TODO(ccgo)					cnt = i = j;
//TODO(ccgo)				}
//TODO(ccgo)			}
//TODO(ccgo)			/* Translate any remaining v4 results to v6 */
//TODO(ccgo)			for (i=0; i<cnt; i++) {
//TODO(ccgo)				if (buf[i].family != AF_INET) continue;
//TODO(ccgo)				memcpy(buf[i].addr+12, buf[i].addr, 4);
//TODO(ccgo)				memcpy(buf[i].addr, "\0\0\0\0\0\0\0\0\0\0\xff\xff", 12);
//TODO(ccgo)				buf[i].family = AF_INET6;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* No further processing is needed if there are fewer than 2
//TODO(ccgo)		 * results or if there are only IPv4 results. */
//TODO(ccgo)		if (cnt<2 || family==AF_INET) return cnt;
//TODO(ccgo)		for (i=0; i<cnt; i++) if (buf[i].family != AF_INET) break;
//TODO(ccgo)		if (i==cnt) return cnt;
//TODO(ccgo)	
//TODO(ccgo)		int cs;
//TODO(ccgo)		pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, &cs);
//TODO(ccgo)	
//TODO(ccgo)		/* The following implements a subset of RFC 3484/6724 destination
//TODO(ccgo)		 * address selection by generating a single 31-bit sort key for
//TODO(ccgo)		 * each address. Rules 3, 4, and 7 are omitted for having
//TODO(ccgo)		 * excessive runtime and code size cost and dubious benefit.
//TODO(ccgo)		 * So far the label/precedence table cannot be customized. */
//TODO(ccgo)		for (i=0; i<cnt; i++) {
//TODO(ccgo)			int family = buf[i].family;
//TODO(ccgo)			int key = 0;
//TODO(ccgo)			struct sockaddr_in6 sa6 = { 0 }, da6 = {
//TODO(ccgo)				.sin6_family = AF_INET6,
//TODO(ccgo)				.sin6_scope_id = buf[i].scopeid,
//TODO(ccgo)				.sin6_port = 65535
//TODO(ccgo)			};
//TODO(ccgo)			struct sockaddr_in sa4 = { 0 }, da4 = {
//TODO(ccgo)				.sin_family = AF_INET,
//TODO(ccgo)				.sin_port = 65535
//TODO(ccgo)			};
//TODO(ccgo)			void *sa, *da;
//TODO(ccgo)			socklen_t salen, dalen;
//TODO(ccgo)			if (family == AF_INET6) {
//TODO(ccgo)				memcpy(da6.sin6_addr.s6_addr, buf[i].addr, 16);
//TODO(ccgo)				da = &da6; dalen = sizeof da6;
//TODO(ccgo)				sa = &sa6; salen = sizeof sa6;
//TODO(ccgo)			} else {
//TODO(ccgo)				memcpy(sa6.sin6_addr.s6_addr,
//TODO(ccgo)					"\0\0\0\0\0\0\0\0\0\0\xff\xff", 12);
//TODO(ccgo)				memcpy(da6.sin6_addr.s6_addr+12, buf[i].addr, 4);
//TODO(ccgo)				memcpy(da6.sin6_addr.s6_addr,
//TODO(ccgo)					"\0\0\0\0\0\0\0\0\0\0\xff\xff", 12);
//TODO(ccgo)				memcpy(da6.sin6_addr.s6_addr+12, buf[i].addr, 4);
//TODO(ccgo)				memcpy(&da4.sin_addr, buf[i].addr, 4);
//TODO(ccgo)				da = &da4; dalen = sizeof da4;
//TODO(ccgo)				sa = &sa4; salen = sizeof sa4;
//TODO(ccgo)			}
//TODO(ccgo)			const struct policy *dpolicy = policyof(&da6.sin6_addr);
//TODO(ccgo)			int dscope = scopeof(&da6.sin6_addr);
//TODO(ccgo)			int dlabel = dpolicy->label;
//TODO(ccgo)			int dprec = dpolicy->prec;
//TODO(ccgo)			int prefixlen = 0;
//TODO(ccgo)			int fd = socket(family, SOCK_DGRAM|SOCK_CLOEXEC, IPPROTO_UDP);
//TODO(ccgo)			if (fd >= 0) {
//TODO(ccgo)				if (!connect(fd, da, dalen)) {
//TODO(ccgo)					key |= DAS_USABLE;
//TODO(ccgo)					if (!getsockname(fd, sa, &salen)) {
//TODO(ccgo)						if (family == AF_INET) memcpy(
//TODO(ccgo)							&sa6.sin6_addr.s6_addr+12,
//TODO(ccgo)							&sa4.sin_addr, 4);
//TODO(ccgo)						if (dscope == scopeof(&sa6.sin6_addr))
//TODO(ccgo)							key |= DAS_MATCHINGSCOPE;
//TODO(ccgo)						if (dlabel == labelof(&sa6.sin6_addr))
//TODO(ccgo)							key |= DAS_MATCHINGLABEL;
//TODO(ccgo)						prefixlen = prefixmatch(&sa6.sin6_addr,
//TODO(ccgo)							&da6.sin6_addr);
//TODO(ccgo)					}
//TODO(ccgo)				}
//TODO(ccgo)				close(fd);
//TODO(ccgo)			}
//TODO(ccgo)			key |= dprec << DAS_PREC_SHIFT;
//TODO(ccgo)			key |= (15-dscope) << DAS_SCOPE_SHIFT;
//TODO(ccgo)			key |= prefixlen << DAS_PREFIX_SHIFT;
//TODO(ccgo)			key |= (MAXADDRS-i) << DAS_ORDER_SHIFT;
//TODO(ccgo)			buf[i].sortkey = key;
//TODO(ccgo)		}
//TODO(ccgo)		qsort(buf, cnt, sizeof *buf, addrcmp);
//TODO(ccgo)	
//TODO(ccgo)		pthread_setcancelstate(cs, 0);
//TODO(ccgo)	
//TODO(ccgo)		return cnt;
//TODO(ccgo)	}

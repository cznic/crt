#include <stdlib.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <string.h>
#include <pthread.h>
#include <unistd.h>
#include <endian.h>
#include <errno.h>
#include "lookup.h"
#include <assert.h>

int getaddrinfo(const char *restrict host, const char *restrict serv, const struct addrinfo *restrict hint, struct addrinfo **restrict res)
{
	struct service ports[MAXSERVS];
	struct address addrs[MAXADDRS];
	char canon[256], *outcanon;
	int nservs, naddrs, nais, canon_len, i, j, k;
	int family = AF_UNSPEC, flags = 0, proto = 0, socktype = 0;
	struct aibuf {
		struct addrinfo ai;
		union sa {
			struct sockaddr_in sin;
			struct sockaddr_in6 sin6;
		} sa;
	} *out;

	if (!host && !serv) return EAI_NONAME;

	if (hint) {
		family = hint->ai_family;
		flags = hint->ai_flags;
		proto = hint->ai_protocol;
		socktype = hint->ai_socktype;

		const int mask = AI_PASSIVE | AI_CANONNAME | AI_NUMERICHOST |
			AI_V4MAPPED | AI_ALL | AI_ADDRCONFIG | AI_NUMERICSERV;
		if ((flags & mask) != flags)
			return EAI_BADFLAGS;

		switch (family) {
		case AF_INET:
		case AF_INET6:
		case AF_UNSPEC:
			break;
		default:
			return EAI_FAMILY;
		}
	}

	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if (flags & AI_ADDRCONFIG) {
//TODO(ccgo)			/* Define the "an address is configured" condition for address
//TODO(ccgo)			 * families via ability to create a socket for the family plus
//TODO(ccgo)			 * routability of the loopback address for the family. */
//TODO(ccgo)			static const struct sockaddr_in lo4 = {
//TODO(ccgo)				.sin_family = AF_INET, .sin_port = 65535,
//TODO(ccgo)				.sin_addr.s_addr = __BYTE_ORDER == __BIG_ENDIAN
//TODO(ccgo)					? 0x7f000001 : 0x0100007f
//TODO(ccgo)			};
//TODO(ccgo)			static const struct sockaddr_in6 lo6 = {
//TODO(ccgo)				.sin6_family = AF_INET6, .sin6_port = 65535,
//TODO(ccgo)				.sin6_addr = IN6ADDR_LOOPBACK_INIT
//TODO(ccgo)			};
//TODO(ccgo)			int tf[2] = { AF_INET, AF_INET6 };
//TODO(ccgo)			const void *ta[2] = { &lo4, &lo6 };
//TODO(ccgo)			socklen_t tl[2] = { sizeof lo4, sizeof lo6 };
//TODO(ccgo)			for (i=0; i<2; i++) {
//TODO(ccgo)				if (family==tf[1-i]) continue;
//TODO(ccgo)				int s = socket(tf[i], SOCK_CLOEXEC|SOCK_DGRAM,
//TODO(ccgo)					IPPROTO_UDP);
//TODO(ccgo)				if (s>=0) {
//TODO(ccgo)					int cs;
//TODO(ccgo)					pthread_setcancelstate(
//TODO(ccgo)						PTHREAD_CANCEL_DISABLE, &cs);
//TODO(ccgo)					int r = connect(s, ta[i], tl[i]);
//TODO(ccgo)					pthread_setcancelstate(cs, 0);
//TODO(ccgo)					close(s);
//TODO(ccgo)					if (!r) continue;
//TODO(ccgo)				}
//TODO(ccgo)				if (errno != EAFNOSUPPORT) return EAI_SYSTEM;
//TODO(ccgo)				if (family == tf[i]) return EAI_NONAME;
//TODO(ccgo)				family = tf[1-i];
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		nservs = __lookup_serv(ports, serv, proto, socktype, flags);
//TODO(ccgo)		if (nservs < 0) return nservs;
//TODO(ccgo)	
//TODO(ccgo)		naddrs = __lookup_name(addrs, canon, host, family, flags);
//TODO(ccgo)		if (naddrs < 0) return naddrs;
//TODO(ccgo)	
//TODO(ccgo)		nais = nservs * naddrs;
//TODO(ccgo)		canon_len = strlen(canon);
//TODO(ccgo)		out = calloc(1, nais * sizeof(*out) + canon_len + 1);
//TODO(ccgo)		if (!out) return EAI_MEMORY;
//TODO(ccgo)	
//TODO(ccgo)		if (canon_len) {
//TODO(ccgo)			outcanon = (void *)&out[nais];
//TODO(ccgo)			memcpy(outcanon, canon, canon_len+1);
//TODO(ccgo)		} else {
//TODO(ccgo)			outcanon = 0;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		for (k=i=0; i<naddrs; i++) for (j=0; j<nservs; j++, k++) {
//TODO(ccgo)			out[k].ai = (struct addrinfo){
//TODO(ccgo)				.ai_family = addrs[i].family,
//TODO(ccgo)				.ai_socktype = ports[j].socktype,
//TODO(ccgo)				.ai_protocol = ports[j].proto,
//TODO(ccgo)				.ai_addrlen = addrs[i].family == AF_INET
//TODO(ccgo)					? sizeof(struct sockaddr_in)
//TODO(ccgo)					: sizeof(struct sockaddr_in6),
//TODO(ccgo)				.ai_addr = (void *)&out[k].sa,
//TODO(ccgo)				.ai_canonname = outcanon,
//TODO(ccgo)				.ai_next = &out[k+1].ai };
//TODO(ccgo)			switch (addrs[i].family) {
//TODO(ccgo)			case AF_INET:
//TODO(ccgo)				out[k].sa.sin.sin_family = AF_INET;
//TODO(ccgo)				out[k].sa.sin.sin_port = htons(ports[j].port);
//TODO(ccgo)				memcpy(&out[k].sa.sin.sin_addr, &addrs[i].addr, 4);
//TODO(ccgo)				break;
//TODO(ccgo)			case AF_INET6:
//TODO(ccgo)				out[k].sa.sin6.sin6_family = AF_INET6;
//TODO(ccgo)				out[k].sa.sin6.sin6_port = htons(ports[j].port);
//TODO(ccgo)				out[k].sa.sin6.sin6_scope_id = addrs[i].scopeid;
//TODO(ccgo)				memcpy(&out[k].sa.sin6.sin6_addr, &addrs[i].addr, 16);
//TODO(ccgo)				break;			
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)		out[nais-1].ai.ai_next = 0;
//TODO(ccgo)		*res = &out->ai;
//TODO(ccgo)		return 0;
}

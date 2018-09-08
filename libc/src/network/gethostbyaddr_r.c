#define _GNU_SOURCE

#include <sys/socket.h>
#include <netdb.h>
#include <string.h>
#include <netinet/in.h>
#include <errno.h>
#include <inttypes.h>
#include <assert.h>

int gethostbyaddr_r(const void *a, socklen_t l, int af,
	struct hostent *h, char *buf, size_t buflen,
	struct hostent **res, int *err)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		union {
//TODO(ccgo)			struct sockaddr_in sin;
//TODO(ccgo)			struct sockaddr_in6 sin6;
//TODO(ccgo)		} sa = { .sin.sin_family = af };
//TODO(ccgo)		socklen_t sl = af==AF_INET6 ? sizeof sa.sin6 : sizeof sa.sin;
//TODO(ccgo)		int i;
//TODO(ccgo)	
//TODO(ccgo)		*res = 0;
//TODO(ccgo)	
//TODO(ccgo)		/* Load address argument into sockaddr structure */
//TODO(ccgo)		if (af==AF_INET6 && l==16) memcpy(&sa.sin6.sin6_addr, a, 16);
//TODO(ccgo)		else if (af==AF_INET && l==4) memcpy(&sa.sin.sin_addr, a, 4);
//TODO(ccgo)		else {
//TODO(ccgo)			*err = NO_RECOVERY;
//TODO(ccgo)			return EINVAL;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* Align buffer and check for space for pointers and ip address */
//TODO(ccgo)		i = (uintptr_t)buf & sizeof(char *)-1;
//TODO(ccgo)		if (!i) i = sizeof(char *);
//TODO(ccgo)		if (buflen <= 5*sizeof(char *)-i + l) return ERANGE;
//TODO(ccgo)		buf += sizeof(char *)-i;
//TODO(ccgo)		buflen -= 5*sizeof(char *)-i + l;
//TODO(ccgo)	
//TODO(ccgo)		h->h_addr_list = (void *)buf;
//TODO(ccgo)		buf += 2*sizeof(char *);
//TODO(ccgo)		h->h_aliases = (void *)buf;
//TODO(ccgo)		buf += 2*sizeof(char *);
//TODO(ccgo)	
//TODO(ccgo)		h->h_addr_list[0] = buf;
//TODO(ccgo)		memcpy(h->h_addr_list[0], a, l);
//TODO(ccgo)		buf += l;
//TODO(ccgo)		h->h_addr_list[1] = 0;
//TODO(ccgo)		h->h_aliases[0] = buf;
//TODO(ccgo)		h->h_aliases[1] = 0;
//TODO(ccgo)	
//TODO(ccgo)		switch (getnameinfo((void *)&sa, sl, buf, buflen, 0, 0, 0)) {
//TODO(ccgo)		case EAI_AGAIN:
//TODO(ccgo)			*err = TRY_AGAIN;
//TODO(ccgo)			return EAGAIN;
//TODO(ccgo)		case EAI_OVERFLOW:
//TODO(ccgo)			return ERANGE;
//TODO(ccgo)		default:
//TODO(ccgo)		case EAI_MEMORY:
//TODO(ccgo)		case EAI_SYSTEM:
//TODO(ccgo)		case EAI_FAIL:
//TODO(ccgo)			*err = NO_RECOVERY;
//TODO(ccgo)			return errno;
//TODO(ccgo)		case 0:
//TODO(ccgo)			break;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		h->h_addrtype = af;
//TODO(ccgo)		h->h_length = l;
//TODO(ccgo)		h->h_name = h->h_aliases[0];
//TODO(ccgo)		*res = h;
//TODO(ccgo)		return 0;
}

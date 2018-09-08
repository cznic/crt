#define _GNU_SOURCE

#include <sys/socket.h>
#include <netdb.h>
#include <string.h>
#include <netinet/in.h>
#include <errno.h>
#include <stdint.h>
#include "lookup.h"
#include <assert.h>

int gethostbyname2_r(const char *name, int af,
	struct hostent *h, char *buf, size_t buflen,
	struct hostent **res, int *err)
{
	struct address addrs[MAXADDRS];
	char canon[256];
	int i, cnt;
	size_t align, need;

	*res = 0;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		cnt = __lookup_name(addrs, canon, name, af, AI_CANONNAME);
//TODO(ccgo)		if (cnt<0) switch (cnt) {
//TODO(ccgo)		case EAI_NONAME:
//TODO(ccgo)			*err = HOST_NOT_FOUND;
//TODO(ccgo)			return ENOENT;
//TODO(ccgo)		case EAI_AGAIN:
//TODO(ccgo)			*err = TRY_AGAIN;
//TODO(ccgo)			return EAGAIN;
//TODO(ccgo)		default:
//TODO(ccgo)		case EAI_FAIL:
//TODO(ccgo)			*err = NO_RECOVERY;
//TODO(ccgo)			return EBADMSG;
//TODO(ccgo)		case EAI_MEMORY:
//TODO(ccgo)		case EAI_SYSTEM:
//TODO(ccgo)			*err = NO_RECOVERY;
//TODO(ccgo)			return errno;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		h->h_addrtype = af;
//TODO(ccgo)		h->h_length = af==AF_INET6 ? 16 : 4;
//TODO(ccgo)	
//TODO(ccgo)		/* Align buffer */
//TODO(ccgo)		align = -(uintptr_t)buf & sizeof(char *)-1;
//TODO(ccgo)	
//TODO(ccgo)		need = 4*sizeof(char *);
//TODO(ccgo)		need += (cnt + 1) * (sizeof(char *) + h->h_length);
//TODO(ccgo)		need += strlen(name)+1;
//TODO(ccgo)		need += strlen(canon)+1;
//TODO(ccgo)		need += align;
//TODO(ccgo)	
//TODO(ccgo)		if (need > buflen) return ERANGE;
//TODO(ccgo)	
//TODO(ccgo)		buf += align;
//TODO(ccgo)		h->h_aliases = (void *)buf;
//TODO(ccgo)		buf += 3*sizeof(char *);
//TODO(ccgo)		h->h_addr_list = (void *)buf;
//TODO(ccgo)		buf += (cnt+1)*sizeof(char *);
//TODO(ccgo)	
//TODO(ccgo)		for (i=0; i<cnt; i++) {
//TODO(ccgo)			h->h_addr_list[i] = (void *)buf;
//TODO(ccgo)			buf += h->h_length;
//TODO(ccgo)			memcpy(h->h_addr_list[i], addrs[i].addr, h->h_length);
//TODO(ccgo)		}
//TODO(ccgo)		h->h_addr_list[i] = 0;
//TODO(ccgo)	
//TODO(ccgo)		h->h_name = h->h_aliases[0] = buf;
//TODO(ccgo)		strcpy(h->h_name, canon);
//TODO(ccgo)		buf += strlen(h->h_name)+1;
//TODO(ccgo)	
//TODO(ccgo)		if (strcmp(h->h_name, name)) {
//TODO(ccgo)			h->h_aliases[1] = buf;
//TODO(ccgo)			strcpy(h->h_aliases[1], name);
//TODO(ccgo)			buf += strlen(h->h_aliases[1])+1;
//TODO(ccgo)		} else h->h_aliases[1] = 0;
//TODO(ccgo)	
//TODO(ccgo)		h->h_aliases[2] = 0;
//TODO(ccgo)	
//TODO(ccgo)		*res = h;
//TODO(ccgo)		return 0;
}

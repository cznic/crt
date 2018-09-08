#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <net/if.h>
#include <arpa/inet.h>
#include <limits.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include "lookup.h"

int __inet_aton(const char *, struct in_addr *);

//TODO(ccgo)	int __lookup_ipliteral(struct address buf[static 1], const char *name, int family)
//TODO(ccgo)	{
//TODO(ccgo)		struct in_addr a4;
//TODO(ccgo)		struct in6_addr a6;
//TODO(ccgo)		if (__inet_aton(name, &a4) > 0) {
//TODO(ccgo)			if (family == AF_INET6) /* wrong family */
//TODO(ccgo)				return EAI_NONAME;
//TODO(ccgo)			memcpy(&buf[0].addr, &a4, sizeof a4);
//TODO(ccgo)			buf[0].family = AF_INET;
//TODO(ccgo)			buf[0].scopeid = 0;
//TODO(ccgo)			return 1;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		char tmp[64];
//TODO(ccgo)		char *p = strchr(name, '%'), *z;
//TODO(ccgo)		unsigned long long scopeid = 0;
//TODO(ccgo)		if (p && p-name < 64) {
//TODO(ccgo)			memcpy(tmp, name, p-name);
//TODO(ccgo)			tmp[p-name] = 0;
//TODO(ccgo)			name = tmp;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (inet_pton(AF_INET6, name, &a6) <= 0)
//TODO(ccgo)			return 0;
//TODO(ccgo)		if (family == AF_INET) /* wrong family */
//TODO(ccgo)			return EAI_NONAME;
//TODO(ccgo)	
//TODO(ccgo)		memcpy(&buf[0].addr, &a6, sizeof a6);
//TODO(ccgo)		buf[0].family = AF_INET6;
//TODO(ccgo)		if (p) {
//TODO(ccgo)			if (isdigit(*++p)) scopeid = strtoull(p, &z, 10);
//TODO(ccgo)			else z = p-1;
//TODO(ccgo)			if (*z) {
//TODO(ccgo)				if (!IN6_IS_ADDR_LINKLOCAL(&a6) &&
//TODO(ccgo)				    !IN6_IS_ADDR_MC_LINKLOCAL(&a6))
//TODO(ccgo)					return EAI_NONAME;
//TODO(ccgo)				scopeid = if_nametoindex(p);
//TODO(ccgo)				if (!scopeid) return EAI_NONAME;
//TODO(ccgo)			}
//TODO(ccgo)			if (scopeid > UINT_MAX) return EAI_NONAME;
//TODO(ccgo)		}
//TODO(ccgo)		buf[0].scopeid = scopeid;
//TODO(ccgo)		return 1;
//TODO(ccgo)	}

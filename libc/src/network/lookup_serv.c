#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <ctype.h>
#include <string.h>
#include <fcntl.h>
#include <errno.h>
#include "lookup.h"
#include "stdio_impl.h"

//TODO(ccgo)	int __lookup_serv(struct service buf[static MAXSERVS], const char *name, int proto, int socktype, int flags)
//TODO(ccgo)	{
//TODO(ccgo)		char line[128];
//TODO(ccgo)		int cnt = 0;
//TODO(ccgo)		char *p, *z = "";
//TODO(ccgo)		unsigned long port = 0;
//TODO(ccgo)	
//TODO(ccgo)		switch (socktype) {
//TODO(ccgo)		case SOCK_STREAM:
//TODO(ccgo)			switch (proto) {
//TODO(ccgo)			case 0:
//TODO(ccgo)				proto = IPPROTO_TCP;
//TODO(ccgo)			case IPPROTO_TCP:
//TODO(ccgo)				break;
//TODO(ccgo)			default:
//TODO(ccgo)				return EAI_SERVICE;
//TODO(ccgo)			}
//TODO(ccgo)			break;
//TODO(ccgo)		case SOCK_DGRAM:
//TODO(ccgo)			switch (proto) {
//TODO(ccgo)			case 0:
//TODO(ccgo)				proto = IPPROTO_UDP;
//TODO(ccgo)			case IPPROTO_UDP:
//TODO(ccgo)				break;
//TODO(ccgo)			default:
//TODO(ccgo)				return EAI_SERVICE;
//TODO(ccgo)			}
//TODO(ccgo)		case 0:
//TODO(ccgo)			break;
//TODO(ccgo)		default:
//TODO(ccgo)			if (name) return EAI_SERVICE;
//TODO(ccgo)			buf[0].port = 0;
//TODO(ccgo)			buf[0].proto = proto;
//TODO(ccgo)			buf[0].socktype = socktype;
//TODO(ccgo)			return 1;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (name) {
//TODO(ccgo)			if (!*name) return EAI_SERVICE;
//TODO(ccgo)			port = strtoul(name, &z, 10);
//TODO(ccgo)		}
//TODO(ccgo)		if (!*z) {
//TODO(ccgo)			if (port > 65535) return EAI_SERVICE;
//TODO(ccgo)			if (proto != IPPROTO_UDP) {
//TODO(ccgo)				buf[cnt].port = port;
//TODO(ccgo)				buf[cnt].socktype = SOCK_STREAM;
//TODO(ccgo)				buf[cnt++].proto = IPPROTO_TCP;
//TODO(ccgo)			}
//TODO(ccgo)			if (proto != IPPROTO_TCP) {
//TODO(ccgo)				buf[cnt].port = port;
//TODO(ccgo)				buf[cnt].socktype = SOCK_DGRAM;
//TODO(ccgo)				buf[cnt++].proto = IPPROTO_UDP;
//TODO(ccgo)			}
//TODO(ccgo)			return cnt;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (flags & AI_NUMERICSERV) return EAI_NONAME;
//TODO(ccgo)	
//TODO(ccgo)		size_t l = strlen(name);
//TODO(ccgo)	
//TODO(ccgo)		unsigned char _buf[1032];
//TODO(ccgo)		FILE _f, *f = __fopen_rb_ca("/etc/services", &_f, _buf, sizeof _buf);
//TODO(ccgo)		if (!f) switch (errno) {
//TODO(ccgo)		case ENOENT:
//TODO(ccgo)		case ENOTDIR:
//TODO(ccgo)		case EACCES:
//TODO(ccgo)			return EAI_SERVICE;
//TODO(ccgo)		default:
//TODO(ccgo)			return EAI_SYSTEM;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		while (fgets(line, sizeof line, f) && cnt < MAXSERVS) {
//TODO(ccgo)			if ((p=strchr(line, '#'))) *p++='\n', *p=0;
//TODO(ccgo)	
//TODO(ccgo)			/* Find service name */
//TODO(ccgo)			for(p=line; (p=strstr(p, name)); p++) {
//TODO(ccgo)				if (p>line && !isspace(p[-1])) continue;
//TODO(ccgo)				if (p[l] && !isspace(p[l])) continue;
//TODO(ccgo)				break;
//TODO(ccgo)			}
//TODO(ccgo)			if (!p) continue;
//TODO(ccgo)	
//TODO(ccgo)			/* Skip past canonical name at beginning of line */
//TODO(ccgo)			for (p=line; *p && !isspace(*p); p++);
//TODO(ccgo)	
//TODO(ccgo)			port = strtoul(p, &z, 10);
//TODO(ccgo)			if (port > 65535 || z==p) continue;
//TODO(ccgo)			if (!strncmp(z, "/udp", 4)) {
//TODO(ccgo)				if (proto == IPPROTO_TCP) continue;
//TODO(ccgo)				buf[cnt].port = port;
//TODO(ccgo)				buf[cnt].socktype = SOCK_DGRAM;
//TODO(ccgo)				buf[cnt++].proto = IPPROTO_UDP;
//TODO(ccgo)			}
//TODO(ccgo)			if (!strncmp(z, "/tcp", 4)) {
//TODO(ccgo)				if (proto == IPPROTO_UDP) continue;
//TODO(ccgo)				buf[cnt].port = port;
//TODO(ccgo)				buf[cnt].socktype = SOCK_STREAM;
//TODO(ccgo)				buf[cnt++].proto = IPPROTO_TCP;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)		__fclose_ca(f);
//TODO(ccgo)		return cnt > 0 ? cnt : EAI_SERVICE;
//TODO(ccgo)	}

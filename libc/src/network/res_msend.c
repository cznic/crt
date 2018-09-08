#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <arpa/inet.h>
#include <stdint.h>
#include <string.h>
#include <poll.h>
#include <time.h>
#include <ctype.h>
#include <unistd.h>
#include <errno.h>
#include <pthread.h>
#include "stdio_impl.h"
#include "syscall.h"
#include "lookup.h"
#include <assert.h>

static void cleanup(void *p)
{
	__syscall(SYS_close, (intptr_t)p);
}

static unsigned long mtime()
{
	struct timespec ts;
	clock_gettime(CLOCK_REALTIME, &ts);
	return (unsigned long)ts.tv_sec * 1000
		+ ts.tv_nsec / 1000000;
}

int __res_msend_rc(int nqueries, const unsigned char *const *queries,
	const int *qlens, unsigned char *const *answers, int *alens, int asize,
	const struct resolvconf *conf)
{
	int fd;
	int timeout, attempts, retry_interval, servfail_retry;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		union {
//TODO(ccgo)			struct sockaddr_in sin;
//TODO(ccgo)			struct sockaddr_in6 sin6;
//TODO(ccgo)		} sa = {0}, ns[MAXNS] = {{0}};
//TODO(ccgo)		socklen_t sl = sizeof sa.sin;
//TODO(ccgo)		int nns = 0;
//TODO(ccgo)		int family = AF_INET;
//TODO(ccgo)		int rlen;
//TODO(ccgo)		int next;
//TODO(ccgo)		int i, j;
//TODO(ccgo)		int cs;
//TODO(ccgo)		struct pollfd pfd;
//TODO(ccgo)		unsigned long t0, t1, t2;
//TODO(ccgo)	
//TODO(ccgo)		pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, &cs);
//TODO(ccgo)	
//TODO(ccgo)		timeout = 1000*conf->timeout;
//TODO(ccgo)		attempts = conf->attempts;
//TODO(ccgo)	
//TODO(ccgo)		for (nns=0; nns<conf->nns; nns++) {
//TODO(ccgo)			const struct address *iplit = &conf->ns[nns];
//TODO(ccgo)			if (iplit->family == AF_INET) {
//TODO(ccgo)				memcpy(&ns[nns].sin.sin_addr, iplit->addr, 4);
//TODO(ccgo)				ns[nns].sin.sin_port = htons(53);
//TODO(ccgo)				ns[nns].sin.sin_family = AF_INET;
//TODO(ccgo)			} else {
//TODO(ccgo)				sl = sizeof sa.sin6;
//TODO(ccgo)				memcpy(&ns[nns].sin6.sin6_addr, iplit->addr, 16);
//TODO(ccgo)				ns[nns].sin6.sin6_port = htons(53);
//TODO(ccgo)				ns[nns].sin6.sin6_scope_id = iplit->scopeid;
//TODO(ccgo)				ns[nns].sin6.sin6_family = family = AF_INET6;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* Get local address and open/bind a socket */
//TODO(ccgo)		sa.sin.sin_family = family;
//TODO(ccgo)		fd = socket(family, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, 0);
//TODO(ccgo)	
//TODO(ccgo)		/* Handle case where system lacks IPv6 support */
//TODO(ccgo)		if (fd < 0 && family == AF_INET6 && errno == EAFNOSUPPORT) {
//TODO(ccgo)			fd = socket(AF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, 0);
//TODO(ccgo)			family = AF_INET;
//TODO(ccgo)		}
//TODO(ccgo)		if (fd < 0 || bind(fd, (void *)&sa, sl) < 0) {
//TODO(ccgo)			if (fd >= 0) close(fd);
//TODO(ccgo)			pthread_setcancelstate(cs, 0);
//TODO(ccgo)			return -1;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* Past this point, there are no errors. Each individual query will
//TODO(ccgo)		 * yield either no reply (indicated by zero length) or an answer
//TODO(ccgo)		 * packet which is up to the caller to interpret. */
//TODO(ccgo)	
//TODO(ccgo)		pthread_cleanup_push(cleanup, (void *)(intptr_t)fd);
//TODO(ccgo)		pthread_setcancelstate(cs, 0);
//TODO(ccgo)	
//TODO(ccgo)		/* Convert any IPv4 addresses in a mixed environment to v4-mapped */
//TODO(ccgo)		if (family == AF_INET6) {
//TODO(ccgo)			setsockopt(fd, IPPROTO_IPV6, IPV6_V6ONLY, &(int){0}, sizeof 0);
//TODO(ccgo)			for (i=0; i<nns; i++) {
//TODO(ccgo)				if (ns[i].sin.sin_family != AF_INET) continue;
//TODO(ccgo)				memcpy(ns[i].sin6.sin6_addr.s6_addr+12,
//TODO(ccgo)					&ns[i].sin.sin_addr, 4);
//TODO(ccgo)				memcpy(ns[i].sin6.sin6_addr.s6_addr,
//TODO(ccgo)					"\0\0\0\0\0\0\0\0\0\0\xff\xff", 12);
//TODO(ccgo)				ns[i].sin6.sin6_family = AF_INET6;
//TODO(ccgo)				ns[i].sin6.sin6_flowinfo = 0;
//TODO(ccgo)				ns[i].sin6.sin6_scope_id = 0;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		memset(alens, 0, sizeof *alens * nqueries);
//TODO(ccgo)	
//TODO(ccgo)		pfd.fd = fd;
//TODO(ccgo)		pfd.events = POLLIN;
//TODO(ccgo)		retry_interval = timeout / attempts;
//TODO(ccgo)		next = 0;
//TODO(ccgo)		t0 = t2 = mtime();
//TODO(ccgo)		t1 = t2 - retry_interval;
//TODO(ccgo)	
//TODO(ccgo)		for (; t2-t0 < timeout; t2=mtime()) {
//TODO(ccgo)			if (t2-t1 >= retry_interval) {
//TODO(ccgo)				/* Query all configured namservers in parallel */
//TODO(ccgo)				for (i=0; i<nqueries; i++)
//TODO(ccgo)					if (!alens[i])
//TODO(ccgo)						for (j=0; j<nns; j++)
//TODO(ccgo)							sendto(fd, queries[i],
//TODO(ccgo)								qlens[i], MSG_NOSIGNAL,
//TODO(ccgo)								(void *)&ns[j], sl);
//TODO(ccgo)				t1 = t2;
//TODO(ccgo)				servfail_retry = 2 * nqueries;
//TODO(ccgo)			}
//TODO(ccgo)	
//TODO(ccgo)			/* Wait for a response, or until time to retry */
//TODO(ccgo)			if (poll(&pfd, 1, t1+retry_interval-t2) <= 0) continue;
//TODO(ccgo)	
//TODO(ccgo)			while ((rlen = recvfrom(fd, answers[next], asize, 0,
//TODO(ccgo)			  (void *)&sa, (socklen_t[1]){sl})) >= 0) {
//TODO(ccgo)	
//TODO(ccgo)				/* Ignore non-identifiable packets */
//TODO(ccgo)				if (rlen < 4) continue;
//TODO(ccgo)	
//TODO(ccgo)				/* Ignore replies from addresses we didn't send to */
//TODO(ccgo)				for (j=0; j<nns && memcmp(ns+j, &sa, sl); j++);
//TODO(ccgo)				if (j==nns) continue;
//TODO(ccgo)	
//TODO(ccgo)				/* Find which query this answer goes with, if any */
//TODO(ccgo)				for (i=next; i<nqueries && (
//TODO(ccgo)					answers[next][0] != queries[i][0] ||
//TODO(ccgo)					answers[next][1] != queries[i][1] ); i++);
//TODO(ccgo)				if (i==nqueries) continue;
//TODO(ccgo)				if (alens[i]) continue;
//TODO(ccgo)	
//TODO(ccgo)				/* Only accept positive or negative responses;
//TODO(ccgo)				 * retry immediately on server failure, and ignore
//TODO(ccgo)				 * all other codes such as refusal. */
//TODO(ccgo)				switch (answers[next][3] & 15) {
//TODO(ccgo)				case 0:
//TODO(ccgo)				case 3:
//TODO(ccgo)					break;
//TODO(ccgo)				case 2:
//TODO(ccgo)					if (servfail_retry && servfail_retry--)
//TODO(ccgo)						sendto(fd, queries[i],
//TODO(ccgo)							qlens[i], MSG_NOSIGNAL,
//TODO(ccgo)							(void *)&ns[j], sl);
//TODO(ccgo)				default:
//TODO(ccgo)					continue;
//TODO(ccgo)				}
//TODO(ccgo)	
//TODO(ccgo)				/* Store answer in the right slot, or update next
//TODO(ccgo)				 * available temp slot if it's already in place. */
//TODO(ccgo)				alens[i] = rlen;
//TODO(ccgo)				if (i == next)
//TODO(ccgo)					for (; next<nqueries && alens[next]; next++);
//TODO(ccgo)				else
//TODO(ccgo)					memcpy(answers[i], answers[next], rlen);
//TODO(ccgo)	
//TODO(ccgo)				if (next == nqueries) goto out;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	out:
//TODO(ccgo)		pthread_cleanup_pop(1);
//TODO(ccgo)	
//TODO(ccgo)		return 0;
}

int __res_msend(int nqueries, const unsigned char *const *queries,
	const int *qlens, unsigned char *const *answers, int *alens, int asize)
{
	struct resolvconf conf;
	if (__get_resolv_conf(&conf, 0, 0) < 0) return -1;
	return __res_msend_rc(nqueries, queries, qlens, answers, alens, asize, &conf);
}

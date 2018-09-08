#include <sys/socket.h>
#include <byteswap.h>
#include <unistd.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>
#include <limits.h>
#include "nscd.h"
#include <assert.h>

static const struct {
	short sun_family;
	char sun_path[21];
} addr = {
	AF_UNIX,
	"/var/run/nscd/socket"
};

FILE *__nscd_query(int32_t req, const char *key, int32_t *buf, size_t len, int *swap)
{
	size_t i;
	int fd;
	FILE *f = 0;
	int32_t req_buf[REQ_LEN] = {
		NSCDVERSION,
		req,
		strnlen(key,LOGIN_NAME_MAX)+1
	};
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		struct msghdr msg = {
//TODO(ccgo)			.msg_iov = (struct iovec[]){
//TODO(ccgo)				{&req_buf, sizeof(req_buf)},
//TODO(ccgo)				{(char*)key, strlen(key)+1}
//TODO(ccgo)			},
//TODO(ccgo)			.msg_iovlen = 2
//TODO(ccgo)		};
//TODO(ccgo)		int errno_save = errno;
//TODO(ccgo)	
//TODO(ccgo)		*swap = 0;
//TODO(ccgo)	retry:
//TODO(ccgo)		memset(buf, 0, len);
//TODO(ccgo)		buf[0] = NSCDVERSION;
//TODO(ccgo)	
//TODO(ccgo)		fd = socket(PF_UNIX, SOCK_STREAM | SOCK_CLOEXEC, 0);
//TODO(ccgo)		if (fd < 0) return NULL;
//TODO(ccgo)	
//TODO(ccgo)		if(!(f = fdopen(fd, "r"))) {
//TODO(ccgo)			close(fd);
//TODO(ccgo)			return 0;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (req_buf[2] > LOGIN_NAME_MAX)
//TODO(ccgo)			return f;
//TODO(ccgo)	
//TODO(ccgo)		if (connect(fd, (struct sockaddr*)&addr, sizeof(addr)) < 0) {
//TODO(ccgo)			/* If there isn't a running nscd we simulate a "not found"
//TODO(ccgo)			 * result and the caller is responsible for calling
//TODO(ccgo)			 * fclose on the (unconnected) socket. The value of
//TODO(ccgo)			 * errno must be left unchanged in this case.  */
//TODO(ccgo)			if (errno == EACCES || errno == ECONNREFUSED || errno == ENOENT) {
//TODO(ccgo)				errno = errno_save;
//TODO(ccgo)				return f;
//TODO(ccgo)			}
//TODO(ccgo)			goto error;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (sendmsg(fd, &msg, MSG_NOSIGNAL) < 0)
//TODO(ccgo)			goto error;
//TODO(ccgo)	
//TODO(ccgo)		if (!fread(buf, len, 1, f)) {
//TODO(ccgo)			/* If the VERSION entry mismatches nscd will disconnect. The
//TODO(ccgo)			 * most likely cause is that the endianness mismatched. So, we
//TODO(ccgo)			 * byteswap and try once more. (if we already swapped, just
//TODO(ccgo)			 * fail out)
//TODO(ccgo)			 */
//TODO(ccgo)			if (ferror(f)) goto error;
//TODO(ccgo)			if (!*swap) {
//TODO(ccgo)				fclose(f);
//TODO(ccgo)				for (i = 0; i < sizeof(req_buf)/sizeof(req_buf[0]); i++) {
//TODO(ccgo)					req_buf[i] = bswap_32(req_buf[i]);
//TODO(ccgo)				}
//TODO(ccgo)				*swap = 1;
//TODO(ccgo)				goto retry;
//TODO(ccgo)			} else {
//TODO(ccgo)				errno = EIO;
//TODO(ccgo)				goto error;
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		if (*swap) {
//TODO(ccgo)			for (i = 0; i < len/sizeof(buf[0]); i++) {
//TODO(ccgo)				buf[i] = bswap_32(buf[i]);
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		/* The first entry in every nscd response is the version number. This
//TODO(ccgo)		 * really shouldn't happen, and is evidence of some form of malformed
//TODO(ccgo)		 * response.
//TODO(ccgo)		 */
//TODO(ccgo)		if(buf[0] != NSCDVERSION) {
//TODO(ccgo)			errno = EIO;
//TODO(ccgo)			goto error;
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		return f;
//TODO(ccgo)	error:
//TODO(ccgo)		fclose(f);
//TODO(ccgo)		return 0;
}

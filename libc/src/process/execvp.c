#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <errno.h>
#include <limits.h>
#include "libc.h"
#include <assert.h>

extern char **__environ;

int __execvpe(const char *file, char *const argv[], char *const envp[])
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		const char *p, *z, *path = getenv("PATH");
//TODO(ccgo)		size_t l, k;
//TODO(ccgo)		int seen_eacces = 0;
//TODO(ccgo)	
//TODO(ccgo)		errno = ENOENT;
//TODO(ccgo)		if (!*file) return -1;
//TODO(ccgo)	
//TODO(ccgo)		if (strchr(file, '/'))
//TODO(ccgo)			return execve(file, argv, envp);
//TODO(ccgo)	
//TODO(ccgo)		if (!path) path = "/usr/local/bin:/bin:/usr/bin";
//TODO(ccgo)		k = strnlen(file, NAME_MAX+1);
//TODO(ccgo)		if (k > NAME_MAX) {
//TODO(ccgo)			errno = ENAMETOOLONG;
//TODO(ccgo)			return -1;
//TODO(ccgo)		}
//TODO(ccgo)		l = strnlen(path, PATH_MAX-1)+1;
//TODO(ccgo)	
//TODO(ccgo)		for(p=path; ; p=z) {
//TODO(ccgo)			char b[l+k+1];
//TODO(ccgo)			z = strchr(p, ':');
//TODO(ccgo)			if (!z) z = p+strlen(p);
//TODO(ccgo)			if (z-p >= l) {
//TODO(ccgo)				if (!*z++) break;
//TODO(ccgo)				continue;
//TODO(ccgo)			}
//TODO(ccgo)			memcpy(b, p, z-p);
//TODO(ccgo)			b[z-p] = '/';
//TODO(ccgo)			memcpy(b+(z-p)+(z>p), file, k+1);
//TODO(ccgo)			execve(b, argv, envp);
//TODO(ccgo)			switch (errno) {
//TODO(ccgo)			case EACCES:
//TODO(ccgo)				seen_eacces = 1;
//TODO(ccgo)			case ENOENT:
//TODO(ccgo)			case ENOTDIR:
//TODO(ccgo)				break;
//TODO(ccgo)			default:
//TODO(ccgo)				return -1;
//TODO(ccgo)			}
//TODO(ccgo)			if (!*z++) break;
//TODO(ccgo)		}
//TODO(ccgo)		if (seen_eacces) errno = EACCES;
//TODO(ccgo)		return -1;
}

int execvp(const char *file, char *const argv[])
{
	return __execvpe(file, argv, __environ);
}

weak_alias(__execvpe, execvpe);

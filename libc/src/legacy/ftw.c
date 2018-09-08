#include <ftw.h>
#include "libc.h"
#include <assert.h>

int ftw(const char *path, int (*fn)(const char *, const struct stat *, int), int fd_limit)
{
	/* The following cast assumes that calling a function with one
	 * argument more than it needs behaves as expected. This is
	 * actually undefined, but works on all real-world machines. */
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return nftw(path, (int (*)())fn, fd_limit, FTW_PHYS);
}

LFS64(ftw);

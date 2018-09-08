#include <unistd.h>
#include "syscall.h"
#include <assert.h>

int fdatasync(int fd)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_fdatasync, fd);
}

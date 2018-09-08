#include <sys/wait.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int waitid(idtype_t type, id_t id, siginfo_t *info, int options)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_waitid, type, id, info, options, 0);
}

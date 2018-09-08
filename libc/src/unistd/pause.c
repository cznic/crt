#include <unistd.h>
#include <signal.h>
#include "syscall.h"
#include "libc.h"
#include <assert.h>

int pause(void)
{
#ifdef SYS_pause
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_pause);
#else
	return syscall_cp(SYS_ppoll, 0, 0, 0, 0);
#endif
}

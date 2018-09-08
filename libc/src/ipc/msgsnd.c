#include <sys/msg.h>
#include "syscall.h"
#include "ipc.h"
#include "libc.h"
#include <assert.h>

int msgsnd(int q, const void *m, size_t len, int flag)
{
#ifdef SYS_msgsnd
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_msgsnd, q, m, len, flag);
#else
	return syscall_cp(SYS_ipc, IPCOP_msgsnd, q, len, flag, m);
#endif
}

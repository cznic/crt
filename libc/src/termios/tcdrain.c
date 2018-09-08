#include <termios.h>
#include <sys/ioctl.h>
#include "libc.h"
#include "syscall.h"
#include <assert.h>

int tcdrain(int fd)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return syscall_cp(SYS_ioctl, fd, TCSBRK, 1);
}

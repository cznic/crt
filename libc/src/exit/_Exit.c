#include <stdlib.h>
#include "syscall.h"

_Noreturn void _Exit(int ec)
{
	__GO__("Log(`==== exit: %v`, _ec)\n");
	__syscall(SYS_exit_group, ec);
	for (;;) __syscall(SYS_exit, ec);
}

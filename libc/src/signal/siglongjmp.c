#include <setjmp.h>
#include <signal.h>
#include "syscall.h"
#include "pthread_impl.h"
#include <assert.h>

_Noreturn void siglongjmp(sigjmp_buf buf, int ret)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		longjmp(buf, ret);
}

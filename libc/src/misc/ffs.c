#include <strings.h>
#include "atomic.h"
#include <libc.h>

int ffs(int i)
{
	return i ? a_ctz_l(i)+1 : 0;
}

weak_alias(ffs, __builtin_ffs);

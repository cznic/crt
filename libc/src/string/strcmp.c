#include <string.h>
#include <libc.h>

int strcmp(const char *l, const char *r)
{
	for (; *l==*r && *l; l++, r++);
	return *(unsigned char *)l - *(unsigned char *)r;
}

weak_alias(strcmp, __builtin_strcmp);

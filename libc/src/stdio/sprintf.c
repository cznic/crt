#include <stdio.h>
#include <stdarg.h>
#include <libc.h>

int sprintf(char *restrict s, const char *restrict fmt, ...)
{
	int ret;
	va_list ap;
	va_start(ap, fmt);
	ret = vsprintf(s, fmt, ap);
	va_end(ap);
	return ret;
}

weak_alias(sprintf, __builtin_sprintf);

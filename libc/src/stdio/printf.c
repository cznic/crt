#include <stdio.h>
#include <stdarg.h>
#include <libc.h>

int printf(const char *restrict fmt, ...)
{
	int ret;
	va_list ap;
	va_start(ap, fmt);
	ret = vfprintf(stdout, fmt, ap);
	va_end(ap);
	return ret;
}

weak_alias(printf, __builtin_printf);

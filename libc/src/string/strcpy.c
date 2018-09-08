#include <string.h>
#include <libc.h>

char *__stpcpy(char *, const char *);

char *strcpy(char *restrict dest, const char *restrict src)
{
#if 1
	__stpcpy(dest, src);
	return dest;
#else
	const unsigned char *s = src;
	unsigned char *d = dest;
	while ((*d++ = *s++));
	return dest;
#endif
}

weak_alias(strcpy, __builtin_strcpy);

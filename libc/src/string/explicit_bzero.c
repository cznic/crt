#define _BSD_SOURCE
#include <string.h>
#include <assert.h>

void explicit_bzero(void *d, size_t n)
{
	d = memset(d, 0, n);
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		__asm__ __volatile__ ("" : : "r"(d) : "memory");
}

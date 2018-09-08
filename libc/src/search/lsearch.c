#include <search.h>
#include <string.h>
#include <assert.h>

void *lsearch(const void *key, void *base, size_t *nelp, size_t width,
	int (*compar)(const void *, const void *))
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		char (*p)[width] = base;
//TODO(ccgo)		size_t n = *nelp;
//TODO(ccgo)		size_t i;
//TODO(ccgo)	
//TODO(ccgo)		for (i = 0; i < n; i++)
//TODO(ccgo)			if (compar(key, p[i]) == 0)
//TODO(ccgo)				return p[i];
//TODO(ccgo)		*nelp = n+1;
//TODO(ccgo)		return memcpy(p[n], key, width);
}

void *lfind(const void *key, const void *base, size_t *nelp,
	size_t width, int (*compar)(const void *, const void *))
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		char (*p)[width] = (void *)base;
//TODO(ccgo)		size_t n = *nelp;
//TODO(ccgo)		size_t i;
//TODO(ccgo)	
//TODO(ccgo)		for (i = 0; i < n; i++)
//TODO(ccgo)			if (compar(key, p[i]) == 0)
//TODO(ccgo)				return p[i];
//TODO(ccgo)		return 0;
}



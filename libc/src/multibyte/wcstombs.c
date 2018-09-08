#include <stdlib.h>
#include <wchar.h>
#include <assert.h>

size_t wcstombs(char *restrict s, const wchar_t *restrict ws, size_t n)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return wcsrtombs(s, &(const wchar_t *){ws}, n, 0);
}

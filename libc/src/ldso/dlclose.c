#include <dlfcn.h>
#include <assert.h>

__attribute__((__visibility__("hidden")))
int __dl_invalid_handle(void *);

int dlclose(void *p)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		return __dl_invalid_handle(p);
}

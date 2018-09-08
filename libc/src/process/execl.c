#include <unistd.h>
#include <stdarg.h>
#include <assert.h>

int execl(const char *path, const char *argv0, ...)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		int argc;
//TODO(ccgo)		va_list ap;
//TODO(ccgo)		va_start(ap, argv0);
//TODO(ccgo)		for (argc=1; va_arg(ap, const char *); argc++);
//TODO(ccgo)		va_end(ap);
//TODO(ccgo)		{
//TODO(ccgo)			int i;
//TODO(ccgo)			char *argv[argc+1];
//TODO(ccgo)			va_start(ap, argv0);
//TODO(ccgo)			argv[0] = (char *)argv0;
//TODO(ccgo)			for (i=1; i<argc; i++)
//TODO(ccgo)				argv[i] = va_arg(ap, char *);
//TODO(ccgo)			argv[i] = NULL;
//TODO(ccgo)			va_end(ap);
//TODO(ccgo)			return execv(path, argv);
//TODO(ccgo)		}
}

#include <utime.h>
#include <sys/stat.h>
#include <time.h>
#include <fcntl.h>
#include <assert.h>

int utime(const char *path, const struct utimbuf *times)
{
//TODO(ccgo)		return utimensat(AT_FDCWD, path, times ? ((struct timespec [2]){
//TODO(ccgo)			{ .tv_sec = times->actime }, { .tv_sec = times->modtime }})
//TODO(ccgo)			: 0, 0);
	struct timespec ts[2] = {0, 0};
	if (times) {
		ts[0].tv_sec = times->actime;
		ts[1].tv_sec = times->modtime;
	}
	return utimensat(AT_FDCWD, path, ts, 0);
}

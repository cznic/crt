package crt

import (
	"fmt"

	"github.com/cznic/crt/errno"
	"github.com/cznic/crt/unistd"
	"golang.org/x/crypto/ssh/terminal"
)

// int access(const char *path, int amode);
func X_access(tls TLS, path uintptr, amode int32) int32 {
	mode := 0
	if amode != unistd.XF_OK {
		panic("access mode not supported")
	}

	f := openFile(tls, GoString(path), mode)
	if f != nil {
		if err := f.Close(); err != nil {
			return -1
		}
		return 0
	}
	// TODO: potentially support more
	return -1
}

// int isatty(int fd);
func X_isatty(tls TLS, fd int32) int32 {
	if terminal.IsTerminal(int(fd)) {
		return 1
	}

	tls.setErrno(errno.XENOTTY)
	return 0
}

func X_wunlink(tls TLS, path uintptr) int32 {
	panic(fmt.Sprintf("TBI: wunlink: %d", path))
}

func X_setmode(tls TLS, fd int32, mode int32) int32 {
	// TODO
	return 0
}

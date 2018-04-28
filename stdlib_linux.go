package crt

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"github.com/cznic/crt/errno"
)

// char *realpath(const char *path, char *resolved_path);
func Xrealpath(tls TLS, path, resolved_path uintptr) uintptr {
	// realpath()  expands  all  symbolic links and resolves references to
	// /./, /../ and extra '/' characters in the null-terminated string
	// named by path to produce a canonicalized absolute path       name.
	// The resulting pathname is stored as a null-terminated string, up to
	// a maximum of PATH_MAX bytes, in the buffer pointed to by
	// resolved_path.  The resulting path will have no symbolic link, /./
	// or /../ components.

	pth := GoString(path)
	resolved, err := filepath.EvalSymlinks(pth)
	// If there is no error, realpath() returns a pointer to the
	// resolved_path.

	// Otherwise, it returns NULL, the contents of the array resolved_path
	// are undefined, and errno is set to indicate the error.
	if err != nil {
		switch {
		case os.IsNotExist(err):
			tls.setErrno(errno.XENOENT)
		default:
			panic(fmt.Errorf("%q: %v", pth, err))
		}
		return 0
	}

	// If resolved_path is specified as NULL, then realpath() uses
	// malloc(3) to allocate a buffer of up to PATH_MAX bytes to hold the
	// resolved pathname, and returns a pointer to this buffer.  The caller
	// should deallocate this buffer using free(3).
	if resolved_path == 0 {
		resolved_path = CString(resolved)
		return 0
	}

	if len(resolved) >= limits.XPATH_MAX {
		panic("TODO")
	}

	b := make([]byte, len(resolved)+1)
	copy(b, resolved)
	Copy(resolved_path, uintptr(unsafe.Pointer(&b[0])), len(b))
	return resolved_path
}

// int system(const char *command);
func Xsystem(tls TLS, command uintptr) int32 {
	if command == 0 || *(*int8)(unsafe.Pointer(command)) == 0 {
		return 1
	}

	cmd := exec.Command("sh", "-c", GoString(command))
	if err := cmd.Run(); err != nil {
		return int32(cmd.ProcessState.Sys().(syscall.WaitStatus))
	}

	return 0
}

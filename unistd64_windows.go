package crt

import (
	"fmt"
	"math"
	"os"
	"syscall"
	"unsafe"
)

// ssize_t read(int fd, void *buf, size_t count);
func Xread(tls *TLS, fd int32, buf unsafe.Pointer, count uint64) int64 { //TODO stdin
	slice := (*[math.MaxInt32]byte)(unsafe.Pointer(buf))[:count]
	r, err := syscall.Read(syscall.Handle(uintptr(fd)), slice)
	if strace {
		fmt.Fprintf(os.Stderr, "read(%v, %#x, %v) %v %v\n", fd, buf, count, r, err)
	}
	if err != nil {
		tls.setErrno(err)
	}
	return int64(r)
}

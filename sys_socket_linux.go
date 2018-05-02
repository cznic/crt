// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"syscall"

	"golang.org/x/sys/unix"
)

// int getsockopt(int sockfd, int level, int optname, void *optval, socklen_t *optlen);
func Xgetsockopt(tls TLS, sockfd, level, optname int32, optval, optlen uintptr) int32 {
	panic("TODO")
}

// int setsockopt(int sockfd, int level, int optname, const void *optval, socklen_t optlen);
func Xsetsockopt(tls TLS, sockfd, level, optname int32, optval uintptr, optlen uint32) int32 {
	panic("TODO")
}

// int getsockname(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
func Xgetsockname(tls TLS, sockfd int32, addr, addrlen uintptr) int32 {
	// The getsockname() function shall retrieve the locally-bound name of
	// the specified socket, store this address in the sockaddr structure
	// pointed to by the address argument, and store the length of this
	// address in the object pointed to by the address_len argument.

	// The address_len argument points to a socklen_t object which on input
	// specifies the length of the supplied sockaddr structure, and on
	// output specifies the length of the stored address. If the actual
	// length of the address is greater than the length of the supplied
	// sockaddr structure, the stored address shall be truncated.

	// If the socket has not been bound to a local name, the value stored
	// in the object pointed to by address is unspecified.
	r, _, err := syscall.Syscall(unix.SYS_GETSOCKNAME, uintptr(sockfd), addr, addrlen)
	if strace {
		fmt.Fprintf(TraceWriter, "getsockname(%#x, %#x, %#x) %v %v\n", sockfd, addr, addrlen, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
		return -1
	}

	return 0

}

// int gethostbyaddr_r(const void *addr, socklen_t len, int type, struct hostent *ret, char *buf, size_t buflen, struct hostent **result, int *h_errnop);
func Xgethostbyaddr_r(tls TLS, addr uintptr, len uint32, type_ int32, ret, buf uintptr, buflen size_t, result, h_errnop uintptr) int32 {
	panic("TODO")
}

// int gethostbyname_r(const char *name, struct hostent *ret, char *buf, size_t buflen, struct hostent **result, int *h_errnop);
func Xgethostbyname_r(tls TLS, name, ret, buf uintptr, buflen size_t, result, h_errnop uintptr) int32 {
	panic("TODO")
}

// ssize_t recv(int sockfd, void *buf, size_t len, int flags);
func Xrecv(tls TLS, sockfd int32, buf uintptr, len size_t, flags int32) ssize_t {
	panic("TODO")
}

// ssize_t send(int sockfd, const void *buf, size_t len, int flags);
func Xsend(tls TLS, sockfd int32, buf uintptr, len size_t, flags int32) ssize_t {
	panic("TODO")
}

// int getpeername(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
func Xgetpeername(tls TLS, sockfd int32, addr, addrlen uintptr) int32 {
	panic("TODO")
}

// int shutdown(int sockfd, int how);
func Xshutdown(tls TLS, sockfd, how int32) int32 {
	panic("TODO")
}

// int getnameinfo(const struct sockaddr *sa, socklen_t salen, char *host, socklen_t hostlen, char *serv, socklen_t servlen, int flags);
func Xgetnameinfo(tls TLS, sa uintptr, salen uint32, host uintptr, hostlen uint32, serv uintptr, servlen uint32, flags int32) int32 {
	panic("TODO")
}

// int socket(int domain, int type, int protocol);
func Xsocket(tls TLS, domain, type_, protocol int32) int32 {
	panic("TODO")
}

// int bind(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
func Xbind(tls TLS, sockfd int32, addr uintptr, addrlen uint32) int32 {
	panic("TODO")
}

// int connect(int sockfd, const struct sockaddr *addr, socklen_t addrlen);
func Xconnect(tls TLS, sockfd int32, addr uintptr, addrlen uint32) int32 {
	panic("TODO")
}

// int accept(int sockfd, struct sockaddr *addr, socklen_t *addrlen);
func Xaccept(tls TLS, sockfd int32, addr, addrlen uintptr) int32 {
	panic("TODO")
}

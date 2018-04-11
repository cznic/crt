// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"net"
	"unsafe"

	"github.com/cznic/crt/netdb"
)

// int getaddrinfo(const char *node, const char *service, const struct addrinfo *hints, struct addrinfo **res);
func Xgetaddrinfo(tls TLS, node, service, hints, res uintptr) int32 {
	panic("TODO")
}

// void freeaddrinfo(struct addrinfo *res);
func Xfreeaddrinfo(tls TLS, res uintptr) {
	panic("TODO")
}

// const char *gai_strerror(int errcode);
func Xgai_strerror(tls TLS, errcode int32) uintptr {
	panic("TODO")
}

// struct servent *getservbyname(const char *name, const char *proto);
func Xgetservbyname(tls TLS, name, proto uintptr) uintptr {
	panic("TODO")
}

// struct hostent *gethostbyaddr(const void *addr, socklen_t len, int type);
func Xgethostbyaddr(tls TLS, addr uintptr, len uint32, type_ int32) uintptr {
	panic("TODO")
}

// struct hostent *gethostbyname(const char *name);
func Xgethostbyname(tls TLS, name uintptr) uintptr {
	// The gethostbyname() function returns a structure of type hostent for
	// the given host name.  Here name is either a hostname or an IPv4
	// address in standard dot notation (as for inet_addr(3)).  If  name
	// is  an  IPv4  address,  no lookup is performed and gethostbyname()
	// simply copies name into the h_name field and its struct in_addr
	// equivalent into the h_addr_list[0] field of the returned hostent
	// structure.  If name doesn't end in a dot and the environment
	// variable HOSTALIASES is set, the alias file pointed to by
	// HOSTALIASES will first be  searched  for  name  (see hostname(7) for
	// the file format).  The current domain and its parents are searched
	// unless name ends in a dot.

	// The  gethostbyname() and gethostbyaddr() functions return the
	// hostent structure or a null pointer if an error occurs.  On error,
	// the h_errno variable holds an error number.  When non-NULL, the
	// return value may point at static data,
	host := GoString(name)
	ip := net.ParseIP(host)
	if ip == nil {
		addrs, err := net.LookupHost(host)
		if err != nil {
			tls.setErrno(err)
			return 0
		}

		r := MustMalloc(int(unsafe.Sizeof(Shostent{})))
		var h Shostent
		h.Xh_name = name
		h.Xh_addrtype = netdb.XAF_INET
		var a []int32
		for _, v := range addrs {
			ip := net.ParseIP(v)
			if ip == nil {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			var n int32
			for _, v := range ip {
				n = n<<8 | int32(v)
			}
			a = append(a, n)
		}
		h.Xh_length = int32(len(a))
		p := MustMalloc(4 * len(a))
		h.Xh_addr_list = p
		for _, v := range a {
			q := MustMalloc(4)
			*(*uintptr)(unsafe.Pointer(p)) = q
			*(*int32)(unsafe.Pointer(q)) = v
		}
		*(*Shostent)(unsafe.Pointer(r)) = h
		return r
	}

	panic(host)
}

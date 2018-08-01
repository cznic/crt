// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// uint32_t ntohl(uint32_t netlong);
func Xntohl(tls TLS, netlong uint32) uint32 {
	panic("TODO")
}

// uint16_t ntohs(uint16_t netshort);
func Xntohs(tls TLS, netshort uint16) uint16 {
	panic("TODO")
}

// uint32_t htonl(uint32_t hostlong);
func Xhtonl(tls TLS, hostlong uint32) uint32 {
	panic("TODO")
}

// uint16_t htons(uint16_t hostshort);
func Xhtons(tls TLS, hostshort uint16) uint16 {
	panic("TODO")
}

// int listen(int sockfd, int backlog);
func Xlisten(tls TLS, sockfd, backlog int32) int32 {
	panic("TODO")
}

// char *inet_ntoa(struct in_addr in);
//
// The inet_ntoa() function shall convert the Internet host address specified
// by in to a string in the Internet standard dot notation.
//
// The inet_ntoa() function need not be thread-safe.
//
// All Internet addresses shall be returned in network order (bytes ordered
// from left to right).
//
// Values specified using IPv4 dotted decimal notation take one of the
// following forms:
//
// a.b.c.d When four parts are specified, each shall be interpreted as a byte
// of data and assigned, from left to right, to the four bytes of an Internet
// address.  a.b.c When a three-part address is specified, the last part shall
// be interpreted as a 16-bit quantity and placed in the rightmost two bytes of
// the network address. This makes the three-part address format convenient for
// specifying Class B network addresses as "128.net.host".  a.b When a two-part
// address is supplied, the last part shall be interpreted as a 24-bit quantity
// and placed in the rightmost three bytes of the network address. This makes
// the two-part address format convenient for specifying Class A network
// addresses as "net.host".  a When only one part is given, the value shall be
// stored directly in the network address without any byte rearrangement.  All
// numbers supplied as parts in IPv4 dotted decimal notation may be decimal,
// octal, or hexadecimal, as specified in the ISO C standard (that is, a
// leading 0x or 0X implies hexadecimal; otherwise, a leading '0' implies
// octal; otherwise, the number is interpreted as decimal).
//
// The inet_ntoa() function shall return a pointer to the network address in
// Internet standard dot notation.
func Xinet_ntoa(tls TLS, in struct{ Xs_addr uint32 }) uintptr {
	panic("TODO")
}

// in_addr_t inet_addr(const char *cp);
//
// The inet_addr() function shall convert the string pointed to by cp, in the
// standard IPv4 dotted decimal notation, to an integer value suitable for use
// as an Internet address.
//
// The inet_ntoa() function shall convert the Internet host address specified
// by in to a string in the Internet standard dot notation.
//
// The inet_ntoa() function need not be thread-safe.
//
// All Internet addresses shall be returned in network order (bytes ordered
// from left to right).
//
// Values specified using IPv4 dotted decimal notation take one of the
// following forms:
//
// a.b.c.d When four parts are specified, each shall be interpreted as a byte
// of data and assigned, from left to right, to the four bytes of an Internet
// address.  a.b.c When a three-part address is specified, the last part shall
// be interpreted as a 16-bit quantity and placed in the rightmost two bytes of
// the network address. This makes the three-part address format convenient for
// specifying Class B network addresses as "128.net.host".  a.b When a two-part
// address is supplied, the last part shall be interpreted as a 24-bit quantity
// and placed in the rightmost three bytes of the network address. This makes
// the two-part address format convenient for specifying Class A network
// addresses as "net.host".  a When only one part is given, the value shall be
// stored directly in the network address without any byte rearrangement.  All
// numbers supplied as parts in IPv4 dotted decimal notation may be decimal,
// octal, or hexadecimal, as specified in the ISO C standard (that is, a
// leading 0x or 0X implies hexadecimal; otherwise, a leading '0' implies
// octal; otherwise, the number is interpreted as decimal).
//
// Upon successful completion, inet_addr() shall return the Internet address.
// Otherwise, it shall return ( in_addr_t)(-1).
//
// The inet_ntoa() function shall return a pointer to the network address in
// Internet standard dot notation.
//
// No errors are defined.
func Xinet_addr(tls TLS, cp uintptr) in_addr_t {
	panic("TODO")
}

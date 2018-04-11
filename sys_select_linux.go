// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int select(int nfds, fd_set *readfds, fd_set *writefds, fd_set *exceptfds, struct timeval *timeout);
func Xselect(tls TLS, nfds int32, readfds, writefds, exceptfds, timeout uintptr) int32 {
	panic("TODO")
}

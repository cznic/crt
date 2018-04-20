// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int fstatfs(int fd, struct statfs *buf);
func Xfstatfs(tls TLS, fd int32, buf uintptr) int32 {
	panic("TODO")
}

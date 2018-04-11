// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// struct group *getgrgid(gid_t gid);
func Xgetgrgid(tls TLS, gid uint32) uintptr {
	panic("TODO")
}

// int getgrgid_r(gid_t gid, struct group *grp, char *buf, size_t buflen, struct group **result);
func Xgetgrgid_r(tls TLS, gid uint32, grp, buf uintptr, buflen size_t, result uintptr) int32 {
	panic("TODO")
}

// struct group *getgrnam(const char *name);
func Xgetgrnam(tls TLS, name uintptr) uintptr {
	panic("TODO")
}

// int getgrnam_r(const char *name, struct group *grp, char *buf, size_t buflen, struct group **result);
func Xgetgrnam_r(tls TLS, name, grp, buf uintptr, buflen size_t, result uintptr) int32 {
	panic("TODO")
}

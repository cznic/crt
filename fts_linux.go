// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// FTS *fts_open(char * const *path_argv, int options, int (*compar)(const FTSENT **, const FTSENT **));
func Xfts_open(tls TLS, path_argv uintptr, options int32, compar uintptr) uintptr {
	panic("TODO")
}

// FTSENT *fts_read(FTS *ftsp);
func Xfts_read(tls TLS, ftsp uintptr) uintptr {
	panic("TODO")
}

// int fts_close(FTS *ftsp);
func Xfts_close(tls TLS, ftsp uintptr) int32 {
	panic("TODO")
}

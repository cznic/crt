// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

type (
	T__gid_t = uint32
	T__uid_t = uint32
)

type TDIR = S__dirstream

type S__dirstream struct {
	buf  uintptr
	next uintptr
	n    int
}

type Shostent struct {
	Xh_name      uintptr // *int8
	Xh_aliases   uintptr // **int8
	Xh_addrtype  int32
	Xh_length    int32
	Xh_addr_list uintptr // **int8
}

type Sin6_addr struct {
	X__in6_u struct {
		X int32
		_ [12]byte
	}
}

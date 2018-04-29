// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

type (
	long_t    = int32
	pthread_t = uint64
	rawmem    [1<<50 - 1]byte
	size_t    = uint64
	ssize_t   = int64
	ulong_t   = uint32
)

type tm struct {
	Xtm_sec      int32
	Xtm_min      int32
	Xtm_hour     int32
	Xtm_mday     int32
	Xtm_mon      int32
	Xtm_year     int32
	Xtm_wday     int32
	Xtm_yday     int32
	Xtm_isdst    int32
	X__tm_gmtoff int64
	X__tm_zone   uintptr // *int8
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

type (
	clock_t       = int32
	long_t        = int32
	pthread_key_t = uint32
	pthread_t     = uint32
	rawmem        [1<<31 - 1]byte
	size_t        = uint32
	ssize_t       = int32
	ulong_t       = uint32
)

type pthread_mutex_t struct { //TODO this is an approximation
	X       int32
	X__kind int32
	_       [20]byte
}

type pthread_mutexattr_t struct{ X int32 }

type Stimespec struct {
	Xtv_sec  int32
	Xtv_nsec int32
}

type passwd struct {
	Xpw_name   uintptr // *int8
	Xpw_passwd uintptr // *int8
	Xpw_uid    uint32
	Xpw_gid    uint32
	Xpw_gecos  uintptr // *int8
	Xpw_dir    uintptr // *int8
	Xpw_shell  uintptr // *int8
}

type utsname struct {
	Xsysname      [65]int8
	Xnodename     [65]int8
	Xrelease      [65]int8
	Xversion      [65]int8
	Xmachine      [65]int8
	X__domainname [65]int8
}

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
	X__tm_gmtoff int32
	X__tm_zone   uintptr // *int8
}

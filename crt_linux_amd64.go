// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

type (
	long_t        = int64
	pthread_key_t = uint32
	pthread_t     = uint64
	rawmem        [1<<50 - 1]byte
	size_t        = uint64
	ssize_t       = int64
	ulong_t       = uint64
)

type Sdirent struct {
	Xd_ino    uint64
	Xd_off    int64
	Xd_reclen uint16
	Xd_type   uint8
	Xd_name   [256]int8
}

type Spasswd struct {
	Xpw_name   uintptr // *int8
	Xpw_passwd uintptr // *int8
	Xpw_uid    uint32
	Xpw_gid    uint32
	Xpw_gecos  uintptr // *int8
	Xpw_dir    uintptr // *int8
	Xpw_shell  uintptr // *int8
}

type Spthread_cond_t struct {
	X int64
	_ [40]byte
}

type Spthread_mutexattr_t struct{ X int32 }

type Spthread_mutex_t struct {
	X int64
	_ [32]byte
}

type Sstat struct {
	Xst_dev           uint64
	Xst_ino           uint64
	Xst_nlink         uint64
	Xst_mode          uint32
	Xst_uid           uint32
	Xst_gid           uint32
	X__pad0           int32
	Xst_rdev          uint64
	Xst_size          int64
	Xst_blksize       int64
	Xst_blocks        int64
	Xst_atim          Stimespec
	Xst_mtim          Stimespec
	Xst_ctim          Stimespec
	X__glibc_reserved [3]int64
}

type Sstat64 struct {
	Xst_dev           uint64
	Xst_ino           uint64
	Xst_nlink         uint64
	Xst_mode          uint32
	Xst_uid           uint32
	Xst_gid           uint32
	X__pad0           int32
	Xst_rdev          uint64
	Xst_size          int64
	Xst_blksize       int64
	Xst_blocks        int64
	Xst_atim          Stimespec
	Xst_mtim          Stimespec
	Xst_ctim          Stimespec
	X__glibc_reserved [3]int64
}

type Stimespec struct {
	Xtv_sec  int64
	Xtv_nsec int64
}

type Stm struct {
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

type Sutimbuf struct {
	Xactime  int64
	Xmodtime int64
}

type Sutsname struct {
	Xsysname      [65]int8
	Xnodename     [65]int8
	Xrelease      [65]int8
	Xversion      [65]int8
	Xmachine      [65]int8
	X__domainname [65]int8
}

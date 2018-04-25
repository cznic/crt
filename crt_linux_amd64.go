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

//TODO export types as aliases.

type dirent struct {
	Xd_ino    uint64
	Xd_off    int64
	Xd_reclen uint16
	Xd_type   uint8
	Xd_name   [256]int8
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

// type Tpthread_cond_t = struct {
// 	X__data [0]struct {
// 		X__lock          int32
// 		X__futex         uint32
// 		X__total_seq     uint64
// 		X__wakeup_seq    uint64
// 		X__woken_seq     uint64
// 		X__mutex         uintptr
// 		X__nwaiters      uint32
// 		X__broadcast_seq uint32
// 	}
// 	X__size  [0][48]int8
// 	X__align [0]int64
// 	X        int64
// 	_        [40]byte
// }

// type Tpthread_mutexattr_t = struct {
// 	X__size  [0][4]int8
// 	X__align [0]int32
// 	X        int32
// }

type pthread_mutexattr_t struct{ X int32 }

// type Tpthread_mutex_t = struct {
// 	X__data [0]struct {
// 		X__lock    int32
// 		X__count   uint32
// 		X__owner   int32
// 		X__nusers  uint32
// 		X__kind    int32
// 		X__spins   int16
// 		X__elision int16
// 		X__list    struct {
// 			X__prev uintptr
// 			X__next uintptr
// 		}
// 	}
// 	X__size  [0][40]int8
// 	X__align [0]int64
// 	X        int64
// 	_        [32]byte
// }

type pthread_mutex_t struct {
	X__lock    int32
	X__count   uint32
	X__owner   int32
	X__nusers  uint32
	X__kind    int32 // pthread.CPTHREAD_MUTEX_NORMAL, ...
	X__spins   int16
	X__elision int16
	X__list    struct {
		X__prev uintptr
		X__next uintptr
	}
}

type stat struct {
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

type stat64 struct {
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

// type Upthread_attr_t struct {
// 	X__size  [0][56]int8
// 	X__align [0]int64
// 	X        int64
// 	_        [48]byte
// }

type utimbuf struct {
	Xactime  int64
	Xmodtime int64
}

type utsname struct {
	Xsysname      [65]int8
	Xnodename     [65]int8
	Xrelease      [65]int8
	Xversion      [65]int8
	Xmachine      [65]int8
	X__domainname [65]int8
}

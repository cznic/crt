// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package crt

type XOSVERSIONINFOA struct {
	X0 uint32
	X1 uint32
	X2 uint32
	X3 uint32
	X4 uint32
	X5 [128]int8
}

type XHMODULE struct {
	X0 int32
}

type XCRITICAL_SECTION struct {
	X0 uintptr
	X1 int32
	X2 int32
	X3 uintptr
	X4 uintptr
	X5 uint64
}

type XFILETIME struct {
	X0 uint32
	X1 uint32
}

type XLARGE_INTEGER struct {
	X [0]struct {
		X0 struct {
			X0 uint32
			X1 int32
		}
		X1 struct {
			X0 uint32
			X1 int32
		}
		X2 int64
	}
	U [8]byte
}

type XSYSTEM_INFO struct {
	X0 struct {
		X [0]struct {
			X0 uint32
			X1 struct {
				X0 uint16
				X1 uint16
			}
		}
		U [4]byte
	}
	X1 uint32
	X2 uintptr
	X3 uintptr
	X4 uint64
	X5 uint32
	X6 uint32
	X7 uint32
	X8 uint16
	X9 uint16
}

type XSYSTEMTIME struct {
	X0 uint16
	X1 uint16
	X2 uint16
	X3 uint16
	X4 uint16
	X5 uint16
	X6 uint16
	X7 uint16
}

type XOSVERSIONINFOW struct {
	X0 uint32
	X1 uint32
	X2 uint32
	X3 uint32
	X4 uint32
	X5 [128]uint16
}

type XOVERLAPPED struct {
	X0 uint64
	X1 uint64
	X2 struct {
		X [0]struct {
			X0 struct {
				X0 uint32
				X1 uint32
			}
			X1 uintptr
		}
		U [8]byte
	}
	X3 uintptr
}

type XSECURITY_ATTRIBUTES struct {
	X0 uint32
	X1 uintptr
	X2 int32
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64
// +build !windows

package crt

const (
	longBits = 64
)

func (r *argsReader) readLong() int64 {
	s := *r
	v := s[0].(int64)
	*r = s[1:]
	return v
}

func (r *argsReader) readULong() uint64 {
	s := *r
	v := s[0].(uint64)
	*r = s[1:]
	return v
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64 amd64p32 arm64 mips64 mips64le mips64p32 mips64p32le ppc64 sparc64

package crt

const (
	longBits = 32
)

func (r *varargReader) readLong() int64 {
	s := *r
	v := s[0].(int32)
	*r = s[1:]
	return int64(v)
}

func (r *varargReader) readULong() uint64 {
	s := *r
	v := s[0].(uint32)
	*r = s[1:]
	return int64(v)
}

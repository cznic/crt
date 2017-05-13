// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int tolower(int c);
func Xtolower(c int32) int32 {
	if c >= 'A' && c <= 'Z' {
		c |= ' '
	}
	return c
}

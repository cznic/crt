// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

type tm *struct {
	X0  int32
	X1  int32
	X2  int32
	X3  int32
	X4  int32
	X5  int32
	X6  int32
	X7  int32
	X8  int32
	X9  int64
	X10 *int8
}

// struct tm *localtime(const time_t *timep);
func Xlocaltime(timep *int64) tm {
	TODO("")
	panic("TODO")
}

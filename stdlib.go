// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os"
)

// void exit(int);
func X__builtin_exit(n int32) {
	os.Exit(int(n))
}

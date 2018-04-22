// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
)

// pid_t waitpid(pid_t pid, int *status, int options);
func Xwaitpid(tls TLS, pid int32, status uintptr, options int32) int32 {
	panic(fmt.Sprintf("%v %#x %#x", pid, status, options))
}

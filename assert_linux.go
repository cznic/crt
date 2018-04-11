// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
)

// void __assert_fail(const char *__assertion, const char *__file, unsigned int __line, const char *__function)
func X__assert_fail(tls TLS, msg, file uintptr, line uint32, fn uintptr) {
	panic(fmt.Errorf("%s.%s:%d: assertion failure: %s", GoString(file), GoString(fn), line, GoString(msg)))
}

// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"runtime"
)

// int sched_yield(void);
func Xsched_yield(tls TLS) int32 {
	// sched_yield() causes the calling thread to relinquish the CPU.  The
	// thread is moved to the end of the queue for its static priority and
	// a new thread gets to run.

	// On success, sched_yield() returns 0.  On error, -1 is returned, and
	// errno is set appropriately.

	runtime.Gosched()
	return 0
}

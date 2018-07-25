// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"syscall"
)

// clock_t times(struct tms *buffer);
//
// The times() function shall fill the tms structure pointed to by buffer with
// time-accounting information. The tms structure is defined in <sys/times.h>.
//
// All times are measured in terms of the number of clock ticks used.
//
// The times of a terminated child process shall be included in the tms_cutime
// and tms_cstime elements of the parent when wait(), waitid(), or waitpid()
// returns the process ID of this terminated child. If a child process has not
// waited for its children, their times shall not be included in its times.
//
// The tms_utime structure member is the CPU time charged for the execution of
// user instructions of the calling process.
//
// The tms_stime structure member is the CPU time charged for execution by the
// system on behalf of the calling process.
//
// The tms_cutime structure member is the sum of the tms_utime and tms_cutime
// times of the child processes.
//
// The tms_cstime structure member is the sum of the tms_stime and tms_cstime
// times of the child processes.
//
// Upon successful completion, times() shall return the elapsed real time, in
// clock ticks, since an arbitrary point in the past (for example, system
// start-up time). This point does not change from one invocation of times()
// within the process to another. The return value may overflow the possible
// range of type clock_t. If times() fails, (clock_t)-1 shall be returned and
// errno set to indicate the error.
//
// The times() function shall fail if:
//
// [EOVERFLOW] The return value would overflow the range of clock_t.
func Xtimes(tls TLS, buffer uintptr) clock_t {
	r, _, err := syscall.Syscall(syscall.SYS_TIMES, buffer, 0, 0)
	if strace {
		fmt.Fprintf(TraceWriter, "times(%#x) %v %v\n", buffer, r, err)
	}
	if err != 0 {
		tls.setErrno(err)
	}
	return clock_t(r)
}

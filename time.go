// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"time"
	"unsafe"
)

var (
	localtime = MustCalloc(int(unsafe.Sizeof(tm{})))
)

// struct tm *localtime(const time_t *timep);
func Xlocaltime(tls TLS, timep uintptr) uintptr { return Xlocaltime_r(tls, timep, localtime) }

// time_t time(time_t *tloc);
func Xtime(tls TLS, tloc uintptr) int64 {
	return time.Now().Unix()
}

// struct tm *localtime_r(const time_t *timep, struct tm *result);
func Xlocaltime_r(tls TLS, timep, stm uintptr) uintptr {
	ut := *(*int64)(unsafe.Pointer(timep))
	t := time.Unix(ut, 0)
	p := (*tm)(unsafe.Pointer(stm))
	p.Xtm_sec = int32(t.Second())
	p.Xtm_min = int32(t.Minute())
	p.Xtm_hour = int32(t.Hour())
	p.Xtm_mday = int32(t.Day())
	p.Xtm_mon = int32(t.Month())
	p.Xtm_year = int32(t.Year())
	p.Xtm_wday = int32(t.Weekday())
	p.Xtm_yday = int32(t.YearDay())
	p.Xtm_isdst = -1 //TODO
	if strace {
		fmt.Fprintf(os.Stderr, "localtime_r(%v, %#x) %+v\n", ut, stm, p)
	}
	return stm
}

// void tzset(void);
//
// The tzset() function shall use the value of the environment variable TZ to
// set time conversion information used by ctime, localtime, mktime, and
// strftime. If TZ is absent from the environment, implementation-defined
// default timezone information shall be used.
//
// The tzset() function shall set the external variable tzname as follows:
//
//	tzname[0] = "std";
//	tzname[1] = "dst";
//
// where std and dst are as described in XBD Environment Variables.
//
// The tzset() function also shall set the external variable daylight
// to 0 if Daylight Savings Time conversions should never be applied for the
// timezone in use; otherwise, non-zero. The external variable timezone shall
// be set to the difference, in seconds, between Coordinated Universal Time
// (UTC) and local standard time.
//
// If a thread accesses tzname, daylight, or timezone ￼  directly while
// another thread is in a call to tzset(), or to any function that is required
// or allowed to set timezone information as if by calling tzset(), the
// behavior is undefined.
//
// The tzset() function shall not return a value.
//
// No errors are defined.
func Xtzset(tls TLS) {
	//TODO so far nothing reads the daylight, timezone or tzname variables (to be defined in crt0.c)
}

// time_t mktime(struct tm *tm);
func Xmktime(tls TLS, tm uintptr) int64 {
	panic("TODO")
}

// struct tm *gmtime_r(const time_t *timep, struct tm *result);
func Xgmtime_r(tls TLS, timep, result uintptr) uintptr {
	panic("TODO")
}

// struct tm *gmtime(const time_t *timer);
//
// For gmtime(): [CX] ￼  The functionality described on this reference page is
// aligned with the ISO C standard. Any conflict between the requirements
// described here and the ISO C standard is unintentional. This volume of
// POSIX.1-2017 defers to the ISO C standard. ￼
//
// The gmtime() function shall convert the time in seconds since the Epoch
// pointed to by timer into a broken-down time, expressed as Coordinated
// Universal Time (UTC).
//
// [CX] ￼ The relationship between a time in seconds since the Epoch used as an
// argument to gmtime() and the tm structure (defined in the <time.h> header)
// is that the result shall be as specified in the expression given in the
// definition of seconds since the Epoch (see XBD Seconds Since the Epoch),
// where the names in the structure and in the expression correspond.
//
// The gmtime() function need not be thread-safe.
//
// The asctime(), ctime(), gmtime(), and localtime() functions shall return
// values in one of two static objects: a broken-down time structure and an
// array of type char. Execution of any of the functions may overwrite the
// information returned in either of these objects by any of the other
// functions.
//
// The gmtime_r() function shall convert the time in seconds since the Epoch
// pointed to by timer into a broken-down time expressed as Coordinated
// Universal Time (UTC). The broken-down time is stored in the structure
// referred to by result. The gmtime_r() function shall also return the address
// of the same structure. ￼
//
// Upon successful completion, the gmtime() function shall return a pointer to
// a struct tm. If an error is detected, gmtime() shall return a null pointer
// [CX] ￼  and set errno to indicate the error.
func Xgmtime(tls TLS, timer uintptr) uintptr {
	panic("TODO")
}

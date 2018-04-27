// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"
)

// /* Fake signal functions.  */
// #define SIG_ERR	((__sighandler_t) -1)	/* Error return.  */
// #define SIG_DFL	((__sighandler_t) 0)	/* Default action.  */
// #define SIG_IGN	((__sighandler_t) 1)	/* Ignore signal.  */
//
// #ifdef __USE_UNIX98
// # define SIG_HOLD	((__sighandler_t) 2)	/* Add signal to hold mask.  */
// #endif
const (
	sigErr  = -1
	sigDfl  = 0
	sigIgn  = 1
	sigHold = 2
)

var (
	sigMap   = map[int32]*sigHandler{}
	sigMapMu sync.Mutex
)

type sigHandler struct {
	ch      chan os.Signal
	handler uintptr
}

// sighandler_t signal(int signum, sighandler_t handler);
func Xsignal(tls TLS, signum int32, handler uintptr) uintptr {
	return X__sysv_signal(tls, signum, handler)
}

// sighandler_t sysv_signal(int signum, sighandler_t handler);
func X__sysv_signal(tls TLS, signum int32, handler uintptr) (r uintptr) {
	sigMapMu.Lock()

	defer sigMapMu.Unlock()

	switch long_t(handler) {
	case sigErr:
		panic("TODO")
	case sigDfl:
		panic("TODO")
	case sigIgn:
		if sh := sigMap[signum]; sh != nil {
			sh.ch <- nil
			r = sh.handler
			delete(sigMap, signum)
		}
		signal.Ignore(syscall.Signal(signum))
	case sigHold:
		panic("TODO")
	default:
		if sh := sigMap[signum]; sh != nil {
			sh.ch <- nil
			r = sh.handler
			delete(sigMap, signum)
		}
		ch := make(chan os.Signal)
		sigMap[signum] = &sigHandler{ch, handler}
		go func() {
			if <-ch == nil {
				return
			}

			(*(*func(TLS, int32))(unsafe.Pointer(&handler)))(tls, signum)
		}()
		signal.Notify(ch, syscall.Signal(signum))
	}
	return r
}

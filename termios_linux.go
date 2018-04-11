// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// int tcgetattr(int fd, struct termios *termios_p);
func Xtcgetattr(tls TLS, fd int32, termios_p uintptr) int32 {
	panic("TODO")
}

// int tcsetattr(int fd, int optional_actions, const struct termios *termios_p);
func Xtcsetattr(tls TLS, fd int32, optional_actions int32, termios_p uintptr) int32 {
	panic("TODO")
}

// int cfsetospeed(struct termios *termios_p, speed_t speed);
func Xcfsetospeed(tls TLS, termios_p uintptr, speed uint32) int32 {
	panic("TODO")
}

// speed_t cfgetospeed(const struct termios *termios_p);
func Xcfgetospeed(tls TLS, termios_p uintptr) uint32 {
	panic("TODO")
}

// int cfsetispeed(struct termios *termios_p, speed_t speed);
func Xcfsetispeed(tls TLS, termios_p uintptr, speed uint32) int32 {
	panic("TODO")
}

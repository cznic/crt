// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"github.com/cznic/crt/ctype"
)

// int toupper(int c);
func Xtoupper(tls TLS, c int32) int32 {
	// The toupper() [CX] ￼  and toupper_l() ￼  functions have as a domain
	// a type int, the value of which is representable as an unsigned char
	// or the value of EOF. If the argument has any other value, the
	// behavior is undefined.
	//
	// If the argument of toupper() [CX] ￼  or toupper_l() ￼  represents a
	// lowercase letter, and there exists a corresponding uppercase letter
	// as defined by character type information in the current locale [CX]
	// ￼  or in the locale represented by locale, ￼  respectively (category
	// LC_CTYPE), the result shall be the corresponding uppercase letter.
	//
	// All other arguments in the domain are returned unchanged.
	//
	// [CX] ￼ The behavior is undefined if the locale argument to
	// toupper_l() is the special locale object LC_GLOBAL_LOCALE or is not
	// a valid locale object handle. ￼
	//
	// Upon successful completion, toupper() [CX] ￼  and toupper_l() ￼
	// shall return the uppercase letter corresponding to the argument
	// passed; otherwise, they shall return the argument unchanged.
	if c >= 'a' && c <= 'z' {
		c &^= ' '
	}
	return c
}

// int tolower(int c);
func Xtolower(tls TLS, c int32) int32 {
	// See the documentation of toupper.
	if c >= 'A' && c <= 'Z' {
		c |= ' '
	}
	return c
}

// int isprint(int c);
func X__builtin_isprint(tls TLS, c int32) int32 {
	if c >= ' ' && c <= '~' {
		return 1
	}

	return 0
}

// int isprint(int c);
func Xisprint(tls TLS, c int32) int32 { return X__builtin_isprint(tls, c) }

// int isspace(int c);
//
// For isspace(): [CX] ￼  The functionality described on this reference page is
// aligned with the ISO C standard. Any conflict between the requirements
// described here and the ISO C standard is unintentional. This volume of
// POSIX.1-2017 defers to the ISO C standard. ￼
//
// The isspace() [CX] ￼  and isspace_l() ￼  functions shall test whether c is a
// character of class space in the current locale, [CX] ￼  or in the locale
// represented by locale, ￼ respectively; see XBD Locale.
//
// The c argument is an int, the value of which the application shall ensure is
// a character representable as an unsigned char or equal to the value of the
// macro EOF. If the argument has any other value, the behavior is undefined.
//
// [CX] ￼ The behavior is undefined if the locale argument to isspace_l() is
// the special locale object LC_GLOBAL_LOCALE or is not a valid locale object
// handle. ￼
//
// The isspace() [CX] ￼  and isspace_l() ￼  functions shall return non-zero if
// c is a white-space character; otherwise, they shall return 0.
func Xisspace(tls TLS, c int32) int32 {
	if c < 0 || c > 255 || _table[c+128]&ctype.C_ISspace == 0 {
		return 0
	}

	return 1
}

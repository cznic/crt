// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"os"

	"github.com/cznic/crt/locale"
)

var locales = map[int32]string{}

// char *setlocale(int category, const char *locale);
func Xsetlocale(tls TLS, category int32, locale_ uintptr) (r uintptr) {
	switch {
	case locale_ != 0:
		// If  locale  is  not NULL, the program's current locale is
		// modified according to the arguments.  The argument category
		// determines which parts of the program's current locale
		// should be modified.
		switch s := GoString(locale_); {
		case s == "":
			// If  locale  is  an  empty  string, "", each part of
			// the locale that should be modified is set according
			// to the environment variables.  The details are
			// implementation-dependent.  For glibc, first
			// (regardless of category), the environment variable
			// LC_ALL is inspected, next the environment variable
			// with the same name as the category (see the table
			// above), and finally the environment variable LANG.
			// The first existing environment variable is used.  If
			// its value is not a valid locale specification, the
			// locale is unchanged, and setlocale() returns NULL.
			switch category {
			case locale.XLC_CTYPE:
				s = os.Getenv("LC_CTYPE")
				locales[category] = s
				r = CString(s)
			default:
				panic(category)
			}
		default:
			switch category {
			case locale.XLC_NUMERIC:
				locales[category] = s
				r = locale_
			default:
				panic(category)
			}
		}
	default:
		// If locale is NULL, the current locale is only queried, not
		// modified.
		panic(category)
	}
	return r
}

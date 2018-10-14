// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"sync"
)

var (
	vaLists  = map[uintptr][]interface{}{}
	vaListMu sync.Mutex
	vaListID uintptr
)

func X__builtin_va_start(tls TLS, ap []interface{}) uintptr {
	vaListMu.Lock()
	vaListID++
	r := vaListID
	vaLists[r] = ap
	vaListMu.Unlock()
	return r
}

func X__builtin_va_copy(tls TLS, src uintptr) uintptr {
	vaListMu.Lock()
	vaListID++
	r := vaListID
	vaLists[r] = vaLists[src]
	vaListMu.Unlock()
	return r
}

func X__builtin_va_end(tls TLS, ap uintptr) {
	vaListMu.Lock()
	delete(vaLists, ap)
	vaListMu.Unlock()
}

func VAfloat32(ap uintptr) float32 { return float32(VAfloat64(ap)) }

func VAfloat64(ap uintptr) (v float64) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case float64:
		v = x
	default:
		panic(fmt.Errorf("crt.VAfloat64 %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAint32(ap uintptr) (v int32) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case int32:
		v = x
	case uint32:
		v = int32(x)
	case int64:
		v = int32(x)
	case uint64:
		v = int32(x)
	case uintptr:
		v = int32(x)
	default:
		panic(fmt.Errorf("crt.VAint32 %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAint64(ap uintptr) (v int64) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case int32:
		v = int64(x)
	case uint32:
		v = int64(x)
	case int64:
		v = x
	case uint64:
		v = int64(x)
	case uintptr:
		v = int64(x)
	default:
		panic(fmt.Errorf("crt.VAint64 %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAuint32(ap uintptr) (v uint32) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case int32:
		v = uint32(x)
	case uint32:
		v = x
	case int64:
		v = uint32(x)
	case uint64:
		v = uint32(x)
	case uintptr:
		v = uint32(x)
	default:
		panic(fmt.Errorf("crt.VAuint32 %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAuint64(ap uintptr) (v uint64) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case int32:
		v = uint64(x)
	case uint32:
		v = uint64(x)
	case int64:
		v = uint64(x)
	case uint64:
		v = x
	case uintptr:
		v = uint64(x)
	default:
		panic(fmt.Errorf("crt.VAuint64 %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAuintptr(ap uintptr) (v uintptr) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	switch x := s[0].(type) {
	case int32:
		v = uintptr(x)
	case uintptr:
		v = x
	default:
		panic(fmt.Errorf("crt.VAuintptr %T", x))
	}
	vaLists[ap] = s[1:]
	return v
}

func VAother(ap uintptr) (v interface{}) {
	vaListMu.Lock()

	defer vaListMu.Unlock()

	s := vaLists[ap]
	if len(s) == 0 {
		return 0
	}

	v = s[0]
	vaLists[ap] = s[1:]
	return v
}

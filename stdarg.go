// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"reflect"
	"unsafe"
)

func roundup(n, to int) int {
	if r := n % to; r != 0 {
		return n + to - r
	}

	return n
}

const (
	vaListDataOff = unsafe.Offsetof(T__builtin_va_list_item{}.Fdata)
)

func X__builtin_va_copy(tls TLS, src uintptr) uintptr {
	sz := int((*T__builtin_va_list_header)(unsafe.Pointer(src)).Fsize)
	dest := MustMalloc(sz)
	Copy(dest, src, sz)
	return dest
}

func X__builtin_va_start(tls TLS, ap []interface{}) uintptr {
	var hdr T__builtin_va_list_header
	b := make([]byte, unsafe.Sizeof(T__builtin_va_list_header{}))
	off := len(b)
	var rq, off0 int
	var sz uintptr
	for i, v := range ap {
		off = roundup(len(b), 8)
		if i == 0 {
			off0 = off
		}
		if off > len(b) {
			b = append(b, make([]byte, off-len(b))...)
		}
		switch x := v.(type) {
		case float64:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*float64)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		case int32:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*int32)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		case int64:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*int64)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		case uint32:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*uint32)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		case uint64:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*uint64)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		case uintptr:
			sz = unsafe.Sizeof(x)
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			*(*uintptr)(unsafe.Pointer(&b[off+int(vaListDataOff)])) = x
		default:
			sz = reflect.TypeOf(x).Size()
			rq = int(vaListDataOff + sz)
			b = append(b, make([]byte, rq)...)
			(*T__builtin_va_list_item)(unsafe.Pointer(&b[off])).Fsize = int32(sz)
			copy(b[off+int(vaListDataOff):], (*rawmem)(unsafe.Pointer(reflect.ValueOf(&x).Elem().InterfaceData()[1]))[:sz])
		}
		off += rq
	}
	p := MustMalloc(off)
	hdr.Fcount = int32(len(ap))
	hdr.Fsize = int32(off)
	hdr.Fitem = p + uintptr(off0)
	*(*T__builtin_va_list_header)(unsafe.Pointer(&b[0])) = hdr
	copy((*rawmem)(unsafe.Pointer(p))[:off], b)
	return p
}

func vaArg(ap uintptr) {
	(*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount--
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	sz := int(vaListDataOff)
	sz += int((*T__builtin_va_list_item)(unsafe.Pointer(p)).Fsize)
	sz = roundup(sz, 8)
	(*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem = p + uintptr(sz)
}

func VAuintptr(ap uintptr) (r uintptr) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	r = *(*uintptr)(unsafe.Pointer(p + vaListDataOff))
	vaArg(ap)
	return r
}

func VAfloat32(ap uintptr) float32 { return float32(VAfloat64(ap)) }

func VAfloat64(ap uintptr) (r float64) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	r = *(*float64)(unsafe.Pointer(p + vaListDataOff))
	vaArg(ap)
	return r
}

func VAint32(ap uintptr) (r int32) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	r = *(*int32)(unsafe.Pointer(p + vaListDataOff))
	vaArg(ap)
	return r
}

func VAuint32(ap uintptr) (r uint32) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	r = *(*uint32)(unsafe.Pointer(p + vaListDataOff))
	vaArg(ap)
	return r
}

func VAint64(ap uintptr) (r int64) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	switch (*T__builtin_va_list_item)(unsafe.Pointer(p)).Fsize {
	case 4:
		r = int64(*(*int32)(unsafe.Pointer(p + vaListDataOff)))
	case 8:
		r = *(*int64)(unsafe.Pointer(p + vaListDataOff))
	default:
		panic("internal error")
	}
	vaArg(ap)
	return r
}

func VAuint64(ap uintptr) (r uint64) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	switch (*T__builtin_va_list_item)(unsafe.Pointer(p)).Fsize {
	case 4:
		r = uint64(*(*uint32)(unsafe.Pointer(p + vaListDataOff)))
	case 8:
		r = *(*uint64)(unsafe.Pointer(p + vaListDataOff))
	default:
		panic("internal error")
	}
	vaArg(ap)
	return r
}

func VAother(ap uintptr) (r uintptr) {
	if (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fcount == 0 {
		return 0
	}

	p := (*T__builtin_va_list_header)(unsafe.Pointer(ap)).Fitem
	vaArg(ap)
	return p + vaListDataOff
}

// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"unsafe"
)

// struct passwd *getpwuid(uid_t uid);
func Xgetpwuid(tls TLS, uid uint32) uintptr {
	u, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		tls.setErrno(err)
		return 0
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		tls.setErrno(err) //TODO Exxx
		return 0
	}

	p, err := Malloc(int(unsafe.Sizeof(Spasswd{})))
	if err != nil {
		tls.setErrno(err) //TODO Exxx
		return 0
	}

	*(*Spasswd)(unsafe.Pointer(p)) = Spasswd{
		Xpw_name:   CString(u.Username),
		Xpw_passwd: CString("x"),
		Xpw_uid:    uid,
		Xpw_gid:    T__gid_t(gid),
		Xpw_gecos:  CString(u.Name),
		Xpw_dir:    CString(u.HomeDir),
		Xpw_shell:  CString(os.Getenv("SHELL")),
	}
	return p
}

// struct passwd *getpwnam(const char *name);
func Xgetpwnam(tls TLS, name uintptr) uintptr {
	panic("TODO")
}

// int getpwnam_r(const char *name, struct passwd *pwd, char *buf, size_t buflen, struct passwd **result);
func Xgetpwnam_r(tls TLS, name, pwd, buf uintptr, buflen size_t, result uintptr) int32 {
	panic("TODO")
}

// int getpwuid_r(uid_t uid, struct passwd *pwd, char *buf, size_t buflen, struct passwd **result);
func Xgetpwuid_r(tls TLS, uid uint32, pwd, buf uintptr, buflen size_t, result uintptr) int32 {
	// The getpwuid_r() function shall update the passwd structure pointed
	// to by pwd and store a pointer to that structure at the location
	// pointed to by result. The structure shall contain an entry from the
	// user database with a matching uid. Storage referenced by the
	// structure is allocated from the memory provided with the buffer
	// parameter, which is bufsize bytes in size. A call to
	// sysconf(_SC_GETPW_R_SIZE_MAX) returns either -1 without changing
	// errno or an initial value suggested for the size of this buffer. A
	// null pointer shall be returned at the location pointed to by result
	// on error or if the requested entry is not found.
	u, err := user.LookupId(fmt.Sprint(uid))
	if err != nil {
		tls.setErrno(err)
		return tls.err()
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		tls.setErrno(err) //TODO Exxx
		return tls.err()
	}

	*(*Spasswd)(unsafe.Pointer(pwd)) = Spasswd{
		Xpw_name:   CString(u.Username),
		Xpw_passwd: CString("x"),
		Xpw_uid:    uid,
		Xpw_gid:    T__gid_t(gid),
		Xpw_gecos:  CString(u.Name),
		Xpw_dir:    CString(u.HomeDir),
		Xpw_shell:  CString(os.Getenv("SHELL")),
	}
	*(*uintptr)(unsafe.Pointer(result)) = pwd

	// If successful, the getpwuid_r() function shall return zero;
	// otherwise, an error number shall be returned to indicate the error.
	return 0
}

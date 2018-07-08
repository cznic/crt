// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

// void rewinddir(DIR *dirp);
//
// The rewinddir() function shall reset the position of the directory stream to
// which dirp refers to the beginning of the directory. It shall also cause the
// directory stream to refer to the current state of the corresponding
// directory, as a call to opendir() would have done. If dirp does not refer to
// a directory stream, the effect is undefined.
//
// After a call to the fork() function, either the parent or child (but not
// both) may continue processing the directory stream using readdir(),
// rewinddir(), or [XSI] ￼ seekdir(). ￼ If both the parent and child processes
// use these functions, the result is undefined.
//
// The rewinddir() function shall not return a value.
//
// No errors are defined.
func Xrewinddir(tls TLS, dirp uintptr) {
	panic("TODO")
}

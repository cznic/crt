// Copyright 2018 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	goos := flag.String("goos", runtime.GOOS, "")
	goarch := flag.String("goarch", runtime.GOARCH, "")
	flag.Parse()

	switch *goos {
	case "linux", "windows":
		// ok
	default:
		panic(fmt.Errorf("unknown/unsupported goos: %s", *goos))
	}

	var libcArch string
	switch *goarch {
	case "amd64":
		libcArch = "x86_64"
	default:
		panic(fmt.Errorf("unknown/unsupported goarch: %s", *goarch))
	}

	if err := os.Chdir("libc"); err != nil {
		panic(err)
	}

	a := []string{
		"src/string/memchr.c",
		"src/string/strnlen.c",
		"src/string/memset.c",
		"src/string/memcpy.c",
		"src/stdio/__stdio_exit.c",
		"src/stdio/__towrite.c",
		"src/stdio/fwrite.c",
		"src/exit/abort.c",
		"src/stdio/__stdout_write.c",
		"src/stdio/__stdio_read.c",
		"src/stdio/__lockfile.c",
		"src/thread/__lock.c",
		"src/stdio/ofl.c",
		"src/stdio/fflush.c",
		"src/aio/aio.c",
		"src/stdio/__stdio_close.c",
		"src/stdio/__stdio_seek.c",
		"src/errno/__errno_location.c",
		"src/internal/syscall_ret.c",
		"src/stdio/__stdio_write.c",
		"src/stdio/stdin.c",
		"src/stdio/stdout.c",
		"src/stdio/stderr.c",
		"src/stdio/fprintf.c",
		"src/exit/assert.c",
		"src/stdio/vfprintf.c",
		"src/stdio/printf.c",
	}
	for i, v := range a {
		a[i] = filepath.FromSlash(v)
	}

	cmd := exec.Command("ccgo", append([]string{
		"--ccgo-full-paths", //TODO-
		"--ccgo-pkg-name", "crt",
		"--ccgo-struct-checks", //TODO-
		"-D_XOPEN_SOURCE=700",
		"-D__typeof=typeof",
		"-I", "include",
		"-I", filepath.Join("arch", "generic"),
		"-I", filepath.Join("arch", libcArch),
		"-I", filepath.Join("obj", "include"),
		"-I", filepath.Join("obj", "src", "internal"),
		"-I", filepath.Join("src", "internal"),
		"-ffreestanding",
		"-o", filepath.Join("..", fmt.Sprintf("libc_%s_%s.go", *goos, *goarch)),
	}, a...)...)
	for _, v := range os.Environ() {
		if v != "GOOS" && v != "GOARCH" {
			cmd.Env = append(cmd.Env, v)
		}
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOOS=%s", *goos), fmt.Sprintf("GOARCH=%s", *goarch))
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Errorf("%s\n%s", err, out))
	}
}

#include <libc.h>
#include <pthread.h>
#include <bits/syscall.h>
#include <assert.h>
#include <byteswap.h>
#include <errno.h>

__attribute__((__visibility__("hidden")))
void (*const __init_array_start)(void)=0, (*const __fini_array_start)(void)=0;

__attribute__((__visibility__("hidden")))
extern void (*const __init_array_end)(void), (*const __fini_array_end)(void);

weak_alias(__init_array_start, __init_array_end);
weak_alias(__fini_array_start, __fini_array_end);

void _init() {}

void _fini() {}

long __syscall(long n, long a1, long a2, long a3, long a4, long a5, long a6) {
	int err;
	__GO__(
		"if !checkSyscall(_n) {\n"
		"\treturn 0 // ignore\n"
		"}\n"
		"x, y, err := syscall.Syscall6(uintptr(_n), uintptr(_a1), uintptr(_a2), uintptr(_a3), uintptr(_a4), uintptr(_a5), uintptr(_a6));\n"
		"switch _n {\n"
		"case DSYS_open:"
		"\tlog(`%s(%q, %#x, %#o, %#x, %#x, %#x) -> (%#x, %#x, %v)`, syscalls[int(_n)], GoString(uintptr(_a1)), _a2, _a3, _a4, _a5, _a6, x, y, err)\n"
		"case DSYS_unlink:"
		"\tlog(`%s(%q, %#x, %#x, %#x, %#x, %#x) -> (%#x, %#x, %v)`, syscalls[int(_n)], GoString(uintptr(_a1)), _a2, _a3, _a4, _a5, _a6, x, y, err)\n"
		"case DSYS_fcntl:"
		"\tlog(`%s(%#x, %s, %#x, %#x, %#x, %#x) -> (%#x, %#x, %v)`, syscalls[int(_n)], _a1, fcntls[int(_a2)], _a3, _a4, _a5, _a6, x, y, err)\n"
		"default:"
		"\tlog(`%s(%#x, %#x, %#x, %#x, %#x, %#x) -> (%#x, %#x, %v)`, syscalls[int(_n)], _a1, _a2, _a3, _a4, _a5, _a6, x, y, err)\n"
		"}\n"
		"_err = int32(err)\n"
		"if _err == DEBADF { println(`EBADF`); println(string(debug.Stack())) }\n"
		"r = long(x)\n"
		"_ = y\n"
		//"fmt.Printf(`syscall(tls %#x, %v, %#x, %#x, %#x, %#x, %#x, %#x) retval %#x, %#x, err %v(%v)`, tls, _n, _a1, _a2, _a3, _a4, _a5, _a6, x, y, int(err), err)\n"
		//"fmt.Println()\n"
	);
	if (err) {
		errno = err;
		return -err;
	}
}

char *__ccgo_arg(int i) { __GO__("return MustCString(os.Args[_i])\n"); }
char *__ccgo_env(int i) { __GO__("return MustCString(env[_i])\n"); }
int  __ccgo_argc(void)  { __GO__("return int32(len(os.Args))\n"); }
int  __ccgo_envc(void)  { __GO__("return int32(len(env))\n"); }
pthread_t __ccgo_main_tls;
size_t _DYNAMIC[1]; //TODO the size if fake

weak_alias(__bswap_64, __builtin_bswap64);
//TODO weak_alias(trap, __builtin_trap);

void  *dlsym(void *__restrict, const char *__restrict) {
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
}

void __builtin_trap(void) {
	abort();
}

void __restore_rt() {
	__syscall(SYS_rt_sigreturn, 0, 0, 0, 0, 0, 0);
}

#include "syscall.h"

long __cancel();

long __syscall_cp_asm(volatile int *p, syscall_arg_t n,
                      syscall_arg_t a1, syscall_arg_t a2, syscall_arg_t a3,
                      syscall_arg_t a4, syscall_arg_t a5, syscall_arg_t a6) {

	if (*p) {
		return __cancel();
	}

	return __syscall(n, a1, a2, a3, a4, a5, a6);
}

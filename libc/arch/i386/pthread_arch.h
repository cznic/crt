static inline struct pthread *__pthread_self()
{
	// struct pthread *self;
	// __asm__ __volatile__ ("movl %%gs:0,%0" : "=r" (self) );
	// return self;
	__GO__("return uintptr(tls)\n");
}

#define TP_ADJ(p) (p)

#define MC_PC gregs[REG_EIP]

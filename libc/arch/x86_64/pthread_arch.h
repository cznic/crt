static inline struct pthread *__pthread_self()
{
	// struct pthread *self;
	// __asm__ __volatile__ ("mov %%fs:0,%0" : "=r" (self) );
	// return self;
	__GO__("return uintptr(tls)\n");
}

#define TP_ADJ(p) (p)

#define MC_PC gregs[REG_RIP]

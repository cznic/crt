#include "pthread_impl.h"
#include <semaphore.h>
#include <unistd.h>
#include <dirent.h>
#include <string.h>
#include <ctype.h>
#include "futex.h"
#include "atomic.h"
#include "../dirent/__dirent.h"
#include <assert.h>

static struct chain {
	struct chain *next;
	int tid;
	sem_t target_sem, caller_sem;
} *volatile head;

static volatile int synccall_lock[1];
static volatile int target_tid;
static void (*callback)(void *), *context;
static volatile int dummy = 0;
weak_alias(dummy, __block_new_threads);

static void handler(int sig)
{
	struct chain ch;
	int old_errno = errno;

	sem_init(&ch.target_sem, 0, 0);
	sem_init(&ch.caller_sem, 0, 0);

	ch.tid = __syscall(SYS_gettid);

	do ch.next = head;
	while (a_cas_p(&head, ch.next, &ch) != ch.next);

	if (a_cas(&target_tid, ch.tid, 0) == (ch.tid | 0x80000000))
		__syscall(SYS_futex, &target_tid, FUTEX_UNLOCK_PI|FUTEX_PRIVATE);

	sem_wait(&ch.target_sem);
	callback(context);
	sem_post(&ch.caller_sem);
	sem_wait(&ch.target_sem);

	errno = old_errno;
}

void __synccall(void (*func)(void *), void *ctx)
{
	sigset_t oldmask;
	int cs, i, r, pid, self;;
	DIR dir = {0};
	struct dirent *de;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		struct sigaction sa = { .sa_flags = SA_RESTART, .sa_handler = handler };
//TODO(ccgo)		struct chain *cp, *next;
//TODO(ccgo)		struct timespec ts;
//TODO(ccgo)	
//TODO(ccgo)		/* Blocking signals in two steps, first only app-level signals
//TODO(ccgo)		 * before taking the lock, then all signals after taking the lock,
//TODO(ccgo)		 * is necessary to achieve AS-safety. Blocking them all first would
//TODO(ccgo)		 * deadlock if multiple threads called __synccall. Waiting to block
//TODO(ccgo)		 * any until after the lock would allow re-entry in the same thread
//TODO(ccgo)		 * with the lock already held. */
//TODO(ccgo)		__block_app_sigs(&oldmask);
//TODO(ccgo)		LOCK(synccall_lock);
//TODO(ccgo)		__block_all_sigs(0);
//TODO(ccgo)		pthread_setcancelstate(PTHREAD_CANCEL_DISABLE, &cs);
//TODO(ccgo)	
//TODO(ccgo)		head = 0;
//TODO(ccgo)	
//TODO(ccgo)		if (!libc.threaded) goto single_threaded;
//TODO(ccgo)	
//TODO(ccgo)		callback = func;
//TODO(ccgo)		context = ctx;
//TODO(ccgo)	
//TODO(ccgo)		/* This atomic store ensures that any signaled threads will see the
//TODO(ccgo)		 * above stores, and prevents more than a bounded number of threads,
//TODO(ccgo)		 * those already in pthread_create, from creating new threads until
//TODO(ccgo)		 * the value is cleared to zero again. */
//TODO(ccgo)		a_store(&__block_new_threads, 1);
//TODO(ccgo)	
//TODO(ccgo)		/* Block even implementation-internal signals, so that nothing
//TODO(ccgo)		 * interrupts the SIGSYNCCALL handlers. The main possible source
//TODO(ccgo)		 * of trouble is asynchronous cancellation. */
//TODO(ccgo)		memset(&sa.sa_mask, -1, sizeof sa.sa_mask);
//TODO(ccgo)		__libc_sigaction(SIGSYNCCALL, &sa, 0);
//TODO(ccgo)	
//TODO(ccgo)		pid = __syscall(SYS_getpid);
//TODO(ccgo)		self = __syscall(SYS_gettid);
//TODO(ccgo)	
//TODO(ccgo)		/* Since opendir is not AS-safe, the DIR needs to be setup manually
//TODO(ccgo)		 * in automatic storage. Thankfully this is easy. */
//TODO(ccgo)		dir.fd = open("/proc/self/task", O_RDONLY|O_DIRECTORY|O_CLOEXEC);
//TODO(ccgo)		if (dir.fd < 0) goto out;
//TODO(ccgo)	
//TODO(ccgo)		/* Initially send one signal per counted thread. But since we can't
//TODO(ccgo)		 * synchronize with thread creation/exit here, there could be too
//TODO(ccgo)		 * few signals. This initial signaling is just an optimization, not
//TODO(ccgo)		 * part of the logic. */
//TODO(ccgo)		for (i=libc.threads_minus_1; i; i--)
//TODO(ccgo)			__syscall(SYS_kill, pid, SIGSYNCCALL);
//TODO(ccgo)	
//TODO(ccgo)		/* Loop scanning the kernel-provided thread list until it shows no
//TODO(ccgo)		 * threads that have not already replied to the signal. */
//TODO(ccgo)		for (;;) {
//TODO(ccgo)			int miss_cnt = 0;
//TODO(ccgo)			while ((de = readdir(&dir))) {
//TODO(ccgo)				if (!isdigit(de->d_name[0])) continue;
//TODO(ccgo)				int tid = atoi(de->d_name);
//TODO(ccgo)				if (tid == self || !tid) continue;
//TODO(ccgo)	
//TODO(ccgo)				/* Set the target thread as the PI futex owner before
//TODO(ccgo)				 * checking if it's in the list of caught threads. If it
//TODO(ccgo)				 * adds itself to the list after we check for it, then
//TODO(ccgo)				 * it will see its own tid in the PI futex and perform
//TODO(ccgo)				 * the unlock operation. */
//TODO(ccgo)				a_store(&target_tid, tid);
//TODO(ccgo)	
//TODO(ccgo)				/* Thread-already-caught is a success condition. */
//TODO(ccgo)				for (cp = head; cp && cp->tid != tid; cp=cp->next);
//TODO(ccgo)				if (cp) continue;
//TODO(ccgo)	
//TODO(ccgo)				r = -__syscall(SYS_tgkill, pid, tid, SIGSYNCCALL);
//TODO(ccgo)	
//TODO(ccgo)				/* Target thread exit is a success condition. */
//TODO(ccgo)				if (r == ESRCH) continue;
//TODO(ccgo)	
//TODO(ccgo)				/* The FUTEX_LOCK_PI operation is used to loan priority
//TODO(ccgo)				 * to the target thread, which otherwise may be unable
//TODO(ccgo)				 * to run. Timeout is necessary because there is a race
//TODO(ccgo)				 * condition where the tid may be reused by a different
//TODO(ccgo)				 * process. */
//TODO(ccgo)				clock_gettime(CLOCK_REALTIME, &ts);
//TODO(ccgo)				ts.tv_nsec += 10000000;
//TODO(ccgo)				if (ts.tv_nsec >= 1000000000) {
//TODO(ccgo)					ts.tv_sec++;
//TODO(ccgo)					ts.tv_nsec -= 1000000000;
//TODO(ccgo)				}
//TODO(ccgo)				r = -__syscall(SYS_futex, &target_tid,
//TODO(ccgo)					FUTEX_LOCK_PI|FUTEX_PRIVATE, 0, &ts);
//TODO(ccgo)	
//TODO(ccgo)				/* Obtaining the lock means the thread responded. ESRCH
//TODO(ccgo)				 * means the target thread exited, which is okay too. */
//TODO(ccgo)				if (!r || r == ESRCH) continue;
//TODO(ccgo)	
//TODO(ccgo)				miss_cnt++;
//TODO(ccgo)			}
//TODO(ccgo)			if (!miss_cnt) break;
//TODO(ccgo)			rewinddir(&dir);
//TODO(ccgo)		}
//TODO(ccgo)		close(dir.fd);
//TODO(ccgo)	
//TODO(ccgo)		/* Serialize execution of callback in caught threads. */
//TODO(ccgo)		for (cp=head; cp; cp=cp->next) {
//TODO(ccgo)			sem_post(&cp->target_sem);
//TODO(ccgo)			sem_wait(&cp->caller_sem);
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)		sa.sa_handler = SIG_IGN;
//TODO(ccgo)		__libc_sigaction(SIGSYNCCALL, &sa, 0);
//TODO(ccgo)	
//TODO(ccgo)	single_threaded:
//TODO(ccgo)		func(ctx);
//TODO(ccgo)	
//TODO(ccgo)		/* Only release the caught threads once all threads, including the
//TODO(ccgo)		 * caller, have returned from the callback function. */
//TODO(ccgo)		for (cp=head; cp; cp=next) {
//TODO(ccgo)			next = cp->next;
//TODO(ccgo)			sem_post(&cp->target_sem);
//TODO(ccgo)		}
//TODO(ccgo)	
//TODO(ccgo)	out:
//TODO(ccgo)		a_store(&__block_new_threads, 0);
//TODO(ccgo)		__wake(&__block_new_threads, -1, 1);
//TODO(ccgo)	
//TODO(ccgo)		pthread_setcancelstate(cs, 0);
//TODO(ccgo)		UNLOCK(synccall_lock);
//TODO(ccgo)		__restore_sigs(&oldmask);
}

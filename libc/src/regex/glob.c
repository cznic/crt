#include <glob.h>
#include <fnmatch.h>
#include <sys/stat.h>
#include <dirent.h>
#include <limits.h>
#include <string.h>
#include <stdlib.h>
#include <errno.h>
#include <stddef.h>
#include "libc.h"
#include <assert.h>

struct match
{
	struct match *next;
	char name[1];
};

static int is_literal(const char *p, int useesc)
{
	int bracket = 0;
	for (; *p; p++) {
		switch (*p) {
		case '\\':
			if (!useesc) break;
		case '?':
		case '*':
			return 0;
		case '[':
			bracket = 1;
			break;
		case ']':
			if (bracket) return 0;
			break;
		}
	}
	return 1;
}

static int append(struct match **tail, const char *name, size_t len, int mark)
{
	struct match *new = malloc(sizeof(struct match) + len + 1);
	if (!new) return -1;
	(*tail)->next = new;
	new->next = NULL;
	strcpy(new->name, name);
	if (mark) strcat(new->name, "/");
	*tail = new;
	return 0;
}

static int match_in_dir(const char *d, const char *p, int flags, int (*errfunc)(const char *path, int err), struct match **tail)
{
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		DIR *dir;
//TODO(ccgo)		struct dirent de_buf, *de;
//TODO(ccgo)		char pat[strlen(p)+1];
//TODO(ccgo)		char *p2;
//TODO(ccgo)		size_t l = strlen(d);
//TODO(ccgo)		int literal;
//TODO(ccgo)		int fnm_flags= ((flags & GLOB_NOESCAPE) ? FNM_NOESCAPE : 0)
//TODO(ccgo)			| ((!(flags & GLOB_PERIOD)) ? FNM_PERIOD : 0);
//TODO(ccgo)		int error;
//TODO(ccgo)	
//TODO(ccgo)		if ((p2 = strchr(p, '/'))) {
//TODO(ccgo)			strcpy(pat, p);
//TODO(ccgo)			pat[p2-p] = 0;
//TODO(ccgo)			for (; *p2 == '/'; p2++);
//TODO(ccgo)			p = pat;
//TODO(ccgo)		}
//TODO(ccgo)		literal = is_literal(p, !(flags & GLOB_NOESCAPE));
//TODO(ccgo)		if (*d == '/' && !*(d+1)) l = 0;
//TODO(ccgo)	
//TODO(ccgo)		/* rely on opendir failing for nondirectory objects */
//TODO(ccgo)		dir = opendir(*d ? d : ".");
//TODO(ccgo)		error = errno;
//TODO(ccgo)		if (!dir) {
//TODO(ccgo)			/* this is not an error -- we let opendir call stat for us */
//TODO(ccgo)			if (error == ENOTDIR) return 0;
//TODO(ccgo)			if (error == EACCES && !*p) {
//TODO(ccgo)				struct stat st;
//TODO(ccgo)				if (!stat(d, &st) && S_ISDIR(st.st_mode)) {
//TODO(ccgo)					if (append(tail, d, l, l))
//TODO(ccgo)						return GLOB_NOSPACE;
//TODO(ccgo)					return 0;
//TODO(ccgo)				}
//TODO(ccgo)			}
//TODO(ccgo)			if (errfunc(d, error) || (flags & GLOB_ERR))
//TODO(ccgo)				return GLOB_ABORTED;
//TODO(ccgo)			return 0;
//TODO(ccgo)		}
//TODO(ccgo)		if (!*p) {
//TODO(ccgo)			error = append(tail, d, l, l) ? GLOB_NOSPACE : 0;
//TODO(ccgo)			closedir(dir);
//TODO(ccgo)			return error;
//TODO(ccgo)		}
//TODO(ccgo)		while (!(error = readdir_r(dir, &de_buf, &de)) && de) {
//TODO(ccgo)			char namebuf[l+de->d_reclen+2], *name = namebuf;
//TODO(ccgo)			if (!literal && fnmatch(p, de->d_name, fnm_flags))
//TODO(ccgo)				continue;
//TODO(ccgo)			if (literal && strcmp(p, de->d_name))
//TODO(ccgo)				continue;
//TODO(ccgo)			if (p2 && de->d_type && !S_ISDIR(de->d_type<<12) && !S_ISLNK(de->d_type<<12))
//TODO(ccgo)				continue;
//TODO(ccgo)			/* With GLOB_PERIOD, don't allow matching . or .. unless
//TODO(ccgo)			 * fnmatch would match them with FNM_PERIOD rules in effect. */
//TODO(ccgo)			if (p2 && (flags & GLOB_PERIOD) && de->d_name[0]=='.'
//TODO(ccgo)			    && (!de->d_name[1] || de->d_name[1]=='.' && !de->d_name[2])
//TODO(ccgo)			    && fnmatch(p, de->d_name, fnm_flags | FNM_PERIOD))
//TODO(ccgo)				continue;
//TODO(ccgo)			if (*d) {
//TODO(ccgo)				memcpy(name, d, l);
//TODO(ccgo)				name[l] = '/';
//TODO(ccgo)				strcpy(name+l+1, de->d_name);
//TODO(ccgo)			} else {
//TODO(ccgo)				name = de->d_name;
//TODO(ccgo)			}
//TODO(ccgo)			if (p2) {
//TODO(ccgo)				if ((error = match_in_dir(name, p2, flags, errfunc, tail))) {
//TODO(ccgo)					closedir(dir);
//TODO(ccgo)					return error;
//TODO(ccgo)				}
//TODO(ccgo)			} else {
//TODO(ccgo)				int mark = 0;
//TODO(ccgo)				if (flags & GLOB_MARK) {
//TODO(ccgo)					if (de->d_type && !S_ISLNK(de->d_type<<12))
//TODO(ccgo)						mark = S_ISDIR(de->d_type<<12);
//TODO(ccgo)					else {
//TODO(ccgo)						struct stat st;
//TODO(ccgo)						stat(name, &st);
//TODO(ccgo)						mark = S_ISDIR(st.st_mode);
//TODO(ccgo)					}
//TODO(ccgo)				}
//TODO(ccgo)				if (append(tail, name, l+de->d_reclen+1, mark)) {
//TODO(ccgo)					closedir(dir);
//TODO(ccgo)					return GLOB_NOSPACE;
//TODO(ccgo)				}
//TODO(ccgo)			}
//TODO(ccgo)		}
//TODO(ccgo)		closedir(dir);
//TODO(ccgo)		if (error && (errfunc(d, error) || (flags & GLOB_ERR)))
//TODO(ccgo)			return GLOB_ABORTED;
//TODO(ccgo)		return 0;
}

static int ignore_err(const char *path, int err)
{
	return 0;
}

static void freelist(struct match *head)
{
	struct match *match, *next;
	for (match=head->next; match; match=next) {
		next = match->next;
		free(match);
	}
}

static int sort(const void *a, const void *b)
{
	return strcmp(*(const char **)a, *(const char **)b);
}

int glob(const char *restrict pat, int flags, int (*errfunc)(const char *path, int err), glob_t *restrict g)
{
	const char *p=pat, *d;
	struct match head = { .next = NULL }, *tail = &head;
	size_t cnt, i;
	size_t offs = (flags & GLOB_DOOFFS) ? g->gl_offs : 0;
	int error = 0;
	
	if (*p == '/') {
		for (; *p == '/'; p++);
		d = "/";
	} else {
		d = "";
	}

	if (!errfunc) errfunc = ignore_err;

	if (!(flags & GLOB_APPEND)) {
		g->gl_offs = offs;
		g->gl_pathc = 0;
		g->gl_pathv = NULL;
	}

	if (strnlen(p, PATH_MAX+1) > PATH_MAX) return GLOB_NOSPACE;

	if (*pat) error = match_in_dir(d, p, flags, errfunc, &tail);
	if (error == GLOB_NOSPACE) {
		freelist(&head);
		return error;
	}
	
	for (cnt=0, tail=head.next; tail; tail=tail->next, cnt++);
	if (!cnt) {
		if (flags & GLOB_NOCHECK) {
			tail = &head;
			if (append(&tail, pat, strlen(pat), 0))
				return GLOB_NOSPACE;
			cnt++;
		} else
			return GLOB_NOMATCH;
	}

	if (flags & GLOB_APPEND) {
		char **pathv = realloc(g->gl_pathv, (offs + g->gl_pathc + cnt + 1) * sizeof(char *));
		if (!pathv) {
			freelist(&head);
			return GLOB_NOSPACE;
		}
		g->gl_pathv = pathv;
		offs += g->gl_pathc;
	} else {
		g->gl_pathv = malloc((offs + cnt + 1) * sizeof(char *));
		if (!g->gl_pathv) {
			freelist(&head);
			return GLOB_NOSPACE;
		}
		for (i=0; i<offs; i++)
			g->gl_pathv[i] = NULL;
	}
	for (i=0, tail=head.next; i<cnt; tail=tail->next, i++)
		g->gl_pathv[offs + i] = tail->name;
	g->gl_pathv[offs + i] = NULL;
	g->gl_pathc += cnt;

	if (!(flags & GLOB_NOSORT))
		qsort(g->gl_pathv+offs, cnt, sizeof(char *), sort);
	
	return error;
}

void globfree(glob_t *g)
{
	size_t i;
	for (i=0; i<g->gl_pathc; i++)
		free(g->gl_pathv[g->gl_offs + i] - offsetof(struct match, name));
	free(g->gl_pathv);
	g->gl_pathc = 0;
	g->gl_pathv = NULL;
}

LFS64(glob);
LFS64(globfree);

#include <netdb.h>
#include <string.h>
#include <assert.h>

/* do we really need all these?? */

static int idx;
//TODO(ccgo)	static const unsigned char protos[] = {
//TODO(ccgo)		"\000ip\0"
//TODO(ccgo)		"\001icmp\0"
//TODO(ccgo)		"\002igmp\0"
//TODO(ccgo)		"\003ggp\0"
//TODO(ccgo)		"\004ipencap\0"
//TODO(ccgo)		"\005st\0"
//TODO(ccgo)		"\006tcp\0"
//TODO(ccgo)		"\010egp\0"
//TODO(ccgo)		"\014pup\0"
//TODO(ccgo)		"\021udp\0"
//TODO(ccgo)		"\024hmp\0"
//TODO(ccgo)		"\026xns-idp\0"
//TODO(ccgo)		"\033rdp\0"
//TODO(ccgo)		"\035iso-tp4\0"
//TODO(ccgo)		"\044xtp\0"
//TODO(ccgo)		"\045ddp\0"
//TODO(ccgo)		"\046idpr-cmtp\0"
//TODO(ccgo)		"\051ipv6\0"
//TODO(ccgo)		"\053ipv6-route\0"
//TODO(ccgo)		"\054ipv6-frag\0"
//TODO(ccgo)		"\055idrp\0"
//TODO(ccgo)		"\056rsvp\0"
//TODO(ccgo)		"\057gre\0"
//TODO(ccgo)		"\062esp\0"
//TODO(ccgo)		"\063ah\0"
//TODO(ccgo)		"\071skip\0"
//TODO(ccgo)		"\072ipv6-icmp\0"
//TODO(ccgo)		"\073ipv6-nonxt\0"
//TODO(ccgo)		"\074ipv6-opts\0"
//TODO(ccgo)		"\111rspf\0"
//TODO(ccgo)		"\121vmtp\0"
//TODO(ccgo)		"\131ospf\0"
//TODO(ccgo)		"\136ipip\0"
//TODO(ccgo)		"\142encap\0"
//TODO(ccgo)		"\147pim\0"
//TODO(ccgo)		"\377raw"
//TODO(ccgo)	};

void endprotoent(void)
{
	idx = 0;
}

void setprotoent(int stayopen)
{
	idx = 0;
}

struct protoent *getprotoent(void)
{
	static struct protoent p;
	static const char *aliases;
	__assert_fail("TODO(ccgo)", __FILE__, __LINE__, __func__);
//TODO(ccgo)		if (idx >= sizeof protos) return NULL;
//TODO(ccgo)		p.p_proto = protos[idx];
//TODO(ccgo)		p.p_name = (char *)&protos[idx+1];
//TODO(ccgo)		p.p_aliases = (char **)&aliases;
//TODO(ccgo)		idx += strlen(p.p_name) + 2;
//TODO(ccgo)		return &p;
}

struct protoent *getprotobyname(const char *name)
{
	struct protoent *p;
	endprotoent();
	do p = getprotoent();
	while (p && strcmp(name, p->p_name));
	return p;
}

struct protoent *getprotobynumber(int num)
{
	struct protoent *p;
	endprotoent();
	do p = getprotoent();
	while (p && p->p_proto != num);
	return p;
}

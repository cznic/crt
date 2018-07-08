// Code generated by $ go generate - DO NOT EDIT.

package inet

const (
	XAF_ALG                           = 38
	XAF_APPLETALK                     = 5
	XAF_ASH                           = 18
	XAF_ATMPVC                        = 8
	XAF_ATMSVC                        = 20
	XAF_AX25                          = 3
	XAF_BLUETOOTH                     = 31
	XAF_BRIDGE                        = 7
	XAF_CAIF                          = 37
	XAF_CAN                           = 29
	XAF_DECnet                        = 12
	XAF_ECONET                        = 19
	XAF_FILE                          = 1
	XAF_IB                            = 27
	XAF_IEEE802154                    = 36
	XAF_INET                          = 2
	XAF_INET6                         = 10
	XAF_IPX                           = 4
	XAF_IRDA                          = 23
	XAF_ISDN                          = 34
	XAF_IUCV                          = 32
	XAF_KCM                           = 41
	XAF_KEY                           = 15
	XAF_LLC                           = 26
	XAF_LOCAL                         = 1
	XAF_MAX                           = 42
	XAF_MPLS                          = 28
	XAF_NETBEUI                       = 13
	XAF_NETLINK                       = 16
	XAF_NETROM                        = 6
	XAF_NFC                           = 39
	XAF_PACKET                        = 17
	XAF_PHONET                        = 35
	XAF_PPPOX                         = 24
	XAF_RDS                           = 21
	XAF_ROSE                          = 11
	XAF_ROUTE                         = 16
	XAF_RXRPC                         = 33
	XAF_SECURITY                      = 14
	XAF_SNA                           = 22
	XAF_TIPC                          = 30
	XAF_UNIX                          = 1
	XAF_UNSPEC                        = 0
	XAF_VSOCK                         = 40
	XAF_WANPIPE                       = 25
	XAF_X25                           = 9
	XBIG_ENDIAN                       = 4321
	XBYTE_ORDER                       = 1234
	XFD_SETSIZE                       = 1024
	XFIOGETOWN                        = 35075
	XFIOSETOWN                        = 35073
	XINET6_ADDRSTRLEN                 = 46
	XINET_ADDRSTRLEN                  = 16
	XINT16_MAX                        = 32767
	XINT16_MIN                        = -32768
	XINT32_MAX                        = 2147483647
	XINT32_MIN                        = -2147483648
	XINT64_MAX                        = 9223372036854775807
	XINT64_MIN                        = -9223372036854775808
	XINT8_MAX                         = 127
	XINT8_MIN                         = -128
	XINTMAX_MAX                       = 9223372036854775807
	XINTMAX_MIN                       = -9223372036854775808
	XINTPTR_MAX                       = 2147483647
	XINTPTR_MIN                       = -2147483648
	XINT_FAST16_MAX                   = 2147483647
	XINT_FAST16_MIN                   = -2147483648
	XINT_FAST32_MAX                   = 2147483647
	XINT_FAST32_MIN                   = -2147483648
	XINT_FAST64_MAX                   = 9223372036854775807
	XINT_FAST64_MIN                   = -9223372036854775808
	XINT_FAST8_MAX                    = 127
	XINT_FAST8_MIN                    = -128
	XINT_LEAST16_MAX                  = 32767
	XINT_LEAST16_MIN                  = -32768
	XINT_LEAST32_MAX                  = 2147483647
	XINT_LEAST32_MIN                  = -2147483648
	XINT_LEAST64_MAX                  = 9223372036854775807
	XINT_LEAST64_MIN                  = -9223372036854775808
	XINT_LEAST8_MAX                   = 127
	XINT_LEAST8_MIN                   = -128
	XIN_CLASSA_HOST                   = 16777215
	XIN_CLASSA_MAX                    = 128
	XIN_CLASSA_NET                    = 4278190080
	XIN_CLASSA_NSHIFT                 = 24
	XIN_CLASSB_HOST                   = 65535
	XIN_CLASSB_MAX                    = 65536
	XIN_CLASSB_NET                    = 4294901760
	XIN_CLASSB_NSHIFT                 = 16
	XIN_CLASSC_HOST                   = 255
	XIN_CLASSC_NET                    = 4294967040
	XIN_CLASSC_NSHIFT                 = 8
	XIN_LOOPBACKNET                   = 127
	XIPV6_2292DSTOPTS                 = 4
	XIPV6_2292HOPLIMIT                = 8
	XIPV6_2292HOPOPTS                 = 3
	XIPV6_2292PKTINFO                 = 2
	XIPV6_2292PKTOPTIONS              = 6
	XIPV6_2292RTHDR                   = 5
	XIPV6_ADDRFORM                    = 1
	XIPV6_ADD_MEMBERSHIP              = 20
	XIPV6_AUTHHDR                     = 10
	XIPV6_CHECKSUM                    = 7
	XIPV6_DONTFRAG                    = 62
	XIPV6_DROP_MEMBERSHIP             = 21
	XIPV6_DSTOPTS                     = 59
	XIPV6_HDRINCL                     = 36
	XIPV6_HOPLIMIT                    = 52
	XIPV6_HOPOPTS                     = 54
	XIPV6_IPSEC_POLICY                = 34
	XIPV6_JOIN_ANYCAST                = 27
	XIPV6_JOIN_GROUP                  = 20
	XIPV6_LEAVE_ANYCAST               = 28
	XIPV6_LEAVE_GROUP                 = 21
	XIPV6_MTU                         = 24
	XIPV6_MTU_DISCOVER                = 23
	XIPV6_MULTICAST_HOPS              = 18
	XIPV6_MULTICAST_IF                = 17
	XIPV6_MULTICAST_LOOP              = 19
	XIPV6_NEXTHOP                     = 9
	XIPV6_PATHMTU                     = 61
	XIPV6_PKTINFO                     = 50
	XIPV6_PMTUDISC_DO                 = 2
	XIPV6_PMTUDISC_DONT               = 0
	XIPV6_PMTUDISC_INTERFACE          = 4
	XIPV6_PMTUDISC_OMIT               = 5
	XIPV6_PMTUDISC_PROBE              = 3
	XIPV6_PMTUDISC_WANT               = 1
	XIPV6_RECVDSTOPTS                 = 58
	XIPV6_RECVERR                     = 25
	XIPV6_RECVHOPLIMIT                = 51
	XIPV6_RECVHOPOPTS                 = 53
	XIPV6_RECVPATHMTU                 = 60
	XIPV6_RECVPKTINFO                 = 49
	XIPV6_RECVRTHDR                   = 56
	XIPV6_RECVTCLASS                  = 66
	XIPV6_ROUTER_ALERT                = 22
	XIPV6_RTHDR                       = 57
	XIPV6_RTHDRDSTOPTS                = 55
	XIPV6_RTHDR_LOOSE                 = 0
	XIPV6_RTHDR_STRICT                = 1
	XIPV6_RTHDR_TYPE_0                = 0
	XIPV6_RXDSTOPTS                   = 59
	XIPV6_RXHOPOPTS                   = 54
	XIPV6_TCLASS                      = 67
	XIPV6_UNICAST_HOPS                = 16
	XIPV6_V6ONLY                      = 26
	XIPV6_XFRM_POLICY                 = 35
	XIP_ADD_MEMBERSHIP                = 35
	XIP_ADD_SOURCE_MEMBERSHIP         = 39
	XIP_BIND_ADDRESS_NO_PORT          = 24
	XIP_BLOCK_SOURCE                  = 38
	XIP_CHECKSUM                      = 23
	XIP_DEFAULT_MULTICAST_LOOP        = 1
	XIP_DEFAULT_MULTICAST_TTL         = 1
	XIP_DROP_MEMBERSHIP               = 36
	XIP_DROP_SOURCE_MEMBERSHIP        = 40
	XIP_FREEBIND                      = 15
	XIP_HDRINCL                       = 3
	XIP_IPSEC_POLICY                  = 16
	XIP_MAX_MEMBERSHIPS               = 20
	XIP_MINTTL                        = 21
	XIP_MSFILTER                      = 41
	XIP_MTU                           = 14
	XIP_MTU_DISCOVER                  = 10
	XIP_MULTICAST_ALL                 = 49
	XIP_MULTICAST_IF                  = 32
	XIP_MULTICAST_LOOP                = 34
	XIP_MULTICAST_TTL                 = 33
	XIP_NODEFRAG                      = 22
	XIP_OPTIONS                       = 4
	XIP_ORIGDSTADDR                   = 20
	XIP_PASSSEC                       = 18
	XIP_PKTINFO                       = 8
	XIP_PKTOPTIONS                    = 9
	XIP_PMTUDISC                      = 10
	XIP_PMTUDISC_DO                   = 2
	XIP_PMTUDISC_DONT                 = 0
	XIP_PMTUDISC_INTERFACE            = 4
	XIP_PMTUDISC_OMIT                 = 5
	XIP_PMTUDISC_PROBE                = 3
	XIP_PMTUDISC_WANT                 = 1
	XIP_RECVERR                       = 11
	XIP_RECVOPTS                      = 6
	XIP_RECVORIGDSTADDR               = 20
	XIP_RECVRETOPTS                   = 7
	XIP_RECVTOS                       = 13
	XIP_RECVTTL                       = 12
	XIP_RETOPTS                       = 7
	XIP_ROUTER_ALERT                  = 5
	XIP_TOS                           = 1
	XIP_TRANSPARENT                   = 19
	XIP_TTL                           = 2
	XIP_UNBLOCK_SOURCE                = 37
	XIP_UNICAST_IF                    = 50
	XIP_XFRM_POLICY                   = 17
	XLITTLE_ENDIAN                    = 1234
	XMCAST_BLOCK_SOURCE               = 43
	XMCAST_EXCLUDE                    = 0
	XMCAST_INCLUDE                    = 1
	XMCAST_JOIN_GROUP                 = 42
	XMCAST_JOIN_SOURCE_GROUP          = 46
	XMCAST_LEAVE_GROUP                = 45
	XMCAST_LEAVE_SOURCE_GROUP         = 47
	XMCAST_MSFILTER                   = 48
	XMCAST_UNBLOCK_SOURCE             = 44
	XPDP_ENDIAN                       = 3412
	XPF_ALG                           = 38
	XPF_APPLETALK                     = 5
	XPF_ASH                           = 18
	XPF_ATMPVC                        = 8
	XPF_ATMSVC                        = 20
	XPF_AX25                          = 3
	XPF_BLUETOOTH                     = 31
	XPF_BRIDGE                        = 7
	XPF_CAIF                          = 37
	XPF_CAN                           = 29
	XPF_DECnet                        = 12
	XPF_ECONET                        = 19
	XPF_FILE                          = 1
	XPF_IB                            = 27
	XPF_IEEE802154                    = 36
	XPF_INET                          = 2
	XPF_INET6                         = 10
	XPF_IPX                           = 4
	XPF_IRDA                          = 23
	XPF_ISDN                          = 34
	XPF_IUCV                          = 32
	XPF_KCM                           = 41
	XPF_KEY                           = 15
	XPF_LLC                           = 26
	XPF_LOCAL                         = 1
	XPF_MAX                           = 42
	XPF_MPLS                          = 28
	XPF_NETBEUI                       = 13
	XPF_NETLINK                       = 16
	XPF_NETROM                        = 6
	XPF_NFC                           = 39
	XPF_PACKET                        = 17
	XPF_PHONET                        = 35
	XPF_PPPOX                         = 24
	XPF_RDS                           = 21
	XPF_ROSE                          = 11
	XPF_ROUTE                         = 16
	XPF_RXRPC                         = 33
	XPF_SECURITY                      = 14
	XPF_SNA                           = 22
	XPF_TIPC                          = 30
	XPF_UNIX                          = 1
	XPF_UNSPEC                        = 0
	XPF_VSOCK                         = 40
	XPF_WANPIPE                       = 25
	XPF_X25                           = 9
	XPTRDIFF_MAX                      = 2147483647
	XPTRDIFF_MIN                      = -2147483648
	XSCM_SRCRT                        = 0
	XSCM_TIMESTAMP                    = 29
	XSCM_TIMESTAMPING                 = 37
	XSCM_TIMESTAMPNS                  = 35
	XSCM_WIFI_STATUS                  = 41
	XSIG_ATOMIC_MAX                   = 2147483647
	XSIG_ATOMIC_MIN                   = -2147483648
	XSIOCATMARK                       = 35077
	XSIOCGPGRP                        = 35076
	XSIOCGSTAMP                       = 35078
	XSIOCGSTAMPNS                     = 35079
	XSIOCSPGRP                        = 35074
	XSIZE_MAX                         = 4294967295
	XSOL_AAL                          = 265
	XSOL_ALG                          = 279
	XSOL_ATM                          = 264
	XSOL_BLUETOOTH                    = 274
	XSOL_CAIF                         = 278
	XSOL_DCCP                         = 269
	XSOL_DECNET                       = 261
	XSOL_ICMPV6                       = 58
	XSOL_IP                           = 0
	XSOL_IPV6                         = 41
	XSOL_IRDA                         = 266
	XSOL_IUCV                         = 277
	XSOL_KCM                          = 281
	XSOL_LLC                          = 268
	XSOL_NETBEUI                      = 267
	XSOL_NETLINK                      = 270
	XSOL_NFC                          = 280
	XSOL_PACKET                       = 263
	XSOL_PNPIPE                       = 275
	XSOL_PPPOL2TP                     = 273
	XSOL_RAW                          = 255
	XSOL_RDS                          = 276
	XSOL_RXRPC                        = 272
	XSOL_SOCKET                       = 1
	XSOL_TIPC                         = 271
	XSOL_X25                          = 262
	XSOMAXCONN                        = 128
	XSO_ACCEPTCONN                    = 30
	XSO_ATTACH_BPF                    = 50
	XSO_ATTACH_FILTER                 = 26
	XSO_ATTACH_REUSEPORT_CBPF         = 51
	XSO_ATTACH_REUSEPORT_EBPF         = 52
	XSO_BINDTODEVICE                  = 25
	XSO_BPF_EXTENSIONS                = 48
	XSO_BROADCAST                     = 6
	XSO_BSDCOMPAT                     = 14
	XSO_BUSY_POLL                     = 46
	XSO_CNX_ADVICE                    = 53
	XSO_DEBUG                         = 1
	XSO_DETACH_BPF                    = 27
	XSO_DETACH_FILTER                 = 27
	XSO_DOMAIN                        = 39
	XSO_DONTROUTE                     = 5
	XSO_ERROR                         = 4
	XSO_GET_FILTER                    = 26
	XSO_INCOMING_CPU                  = 49
	XSO_KEEPALIVE                     = 9
	XSO_LINGER                        = 13
	XSO_LOCK_FILTER                   = 44
	XSO_MARK                          = 36
	XSO_MAX_PACING_RATE               = 47
	XSO_NOFCS                         = 43
	XSO_NO_CHECK                      = 11
	XSO_OOBINLINE                     = 10
	XSO_PASSCRED                      = 16
	XSO_PASSSEC                       = 34
	XSO_PEEK_OFF                      = 42
	XSO_PEERCRED                      = 17
	XSO_PEERNAME                      = 28
	XSO_PEERSEC                       = 31
	XSO_PRIORITY                      = 12
	XSO_PROTOCOL                      = 38
	XSO_RCVBUF                        = 8
	XSO_RCVBUFFORCE                   = 33
	XSO_RCVLOWAT                      = 18
	XSO_RCVTIMEO                      = 20
	XSO_REUSEADDR                     = 2
	XSO_REUSEPORT                     = 15
	XSO_RXQ_OVFL                      = 40
	XSO_SECURITY_AUTHENTICATION       = 22
	XSO_SECURITY_ENCRYPTION_NETWORK   = 24
	XSO_SECURITY_ENCRYPTION_TRANSPORT = 23
	XSO_SELECT_ERR_QUEUE              = 45
	XSO_SNDBUF                        = 7
	XSO_SNDBUFFORCE                   = 32
	XSO_SNDLOWAT                      = 19
	XSO_SNDTIMEO                      = 21
	XSO_TIMESTAMP                     = 29
	XSO_TIMESTAMPING                  = 37
	XSO_TIMESTAMPNS                   = 35
	XSO_TYPE                          = 3
	XSO_WIFI_STATUS                   = 41
	XUINT16_MAX                       = 65535
	XUINT32_MAX                       = 4294967295
	XUINT64_MAX                       = 18446744073709551615
	XUINT8_MAX                        = 255
	XUINTMAX_MAX                      = 18446744073709551615
	XUINTPTR_MAX                      = 4294967295
	XUINT_FAST16_MAX                  = 4294967295
	XUINT_FAST32_MAX                  = 4294967295
	XUINT_FAST64_MAX                  = 18446744073709551615
	XUINT_FAST8_MAX                   = 255
	XUINT_LEAST16_MAX                 = 65535
	XUINT_LEAST32_MAX                 = 4294967295
	XUINT_LEAST64_MAX                 = 18446744073709551615
	XUINT_LEAST8_MAX                  = 255
	XUIO_MAXIOV                       = 1024
	XWCHAR_MAX                        = 4294967295
	XWCHAR_MIN                        = 0
	XWINT_MAX                         = 4294967295
	XWINT_MIN                         = 0

	CIPPORT_BIFFUDP      = 512
	CIPPORT_CMDSERVER    = 514
	CIPPORT_DAYTIME      = 13
	CIPPORT_DISCARD      = 9
	CIPPORT_ECHO         = 7
	CIPPORT_EFSSERVER    = 520
	CIPPORT_EXECSERVER   = 512
	CIPPORT_FINGER       = 79
	CIPPORT_FTP          = 21
	CIPPORT_LOGINSERVER  = 513
	CIPPORT_MTP          = 57
	CIPPORT_NAMESERVER   = 42
	CIPPORT_NETSTAT      = 15
	CIPPORT_RESERVED     = 1024
	CIPPORT_RJE          = 77
	CIPPORT_ROUTESERVER  = 520
	CIPPORT_SMTP         = 25
	CIPPORT_SUPDUP       = 95
	CIPPORT_SYSTAT       = 11
	CIPPORT_TELNET       = 23
	CIPPORT_TFTP         = 69
	CIPPORT_TIMESERVER   = 37
	CIPPORT_TTYLINK      = 87
	CIPPORT_USERRESERVED = 5000
	CIPPORT_WHOIS        = 43
	CIPPORT_WHOSERVER    = 513
	CIPPROTO_AH          = 51
	CIPPROTO_BEETPH      = 94
	CIPPROTO_COMP        = 108
	CIPPROTO_DCCP        = 33
	CIPPROTO_DSTOPTS     = 60
	CIPPROTO_EGP         = 8
	CIPPROTO_ENCAP       = 98
	CIPPROTO_ESP         = 50
	CIPPROTO_FRAGMENT    = 44
	CIPPROTO_GRE         = 47
	CIPPROTO_HOPOPTS     = 0
	CIPPROTO_ICMP        = 1
	CIPPROTO_ICMPV6      = 58
	CIPPROTO_IDP         = 22
	CIPPROTO_IGMP        = 2
	CIPPROTO_IP          = 0
	CIPPROTO_IPIP        = 4
	CIPPROTO_IPV6        = 41
	CIPPROTO_MAX         = 256
	CIPPROTO_MH          = 135
	CIPPROTO_MPLS        = 137
	CIPPROTO_MTP         = 92
	CIPPROTO_NONE        = 59
	CIPPROTO_PIM         = 103
	CIPPROTO_PUP         = 12
	CIPPROTO_RAW         = 255
	CIPPROTO_ROUTING     = 43
	CIPPROTO_RSVP        = 46
	CIPPROTO_SCTP        = 132
	CIPPROTO_TCP         = 6
	CIPPROTO_TP          = 29
	CIPPROTO_UDP         = 17
	CIPPROTO_UDPLITE     = 136
	CMSG_BATCH           = 262144
	CMSG_CMSG_CLOEXEC    = 1073741824
	CMSG_CONFIRM         = 2048
	CMSG_CTRUNC          = 8
	CMSG_DONTROUTE       = 4
	CMSG_DONTWAIT        = 64
	CMSG_EOR             = 128
	CMSG_ERRQUEUE        = 8192
	CMSG_FASTOPEN        = 536870912
	CMSG_FIN             = 512
	CMSG_MORE            = 32768
	CMSG_NOSIGNAL        = 16384
	CMSG_OOB             = 1
	CMSG_PEEK            = 2
	CMSG_PROXY           = 16
	CMSG_RST             = 4096
	CMSG_SYN             = 1024
	CMSG_TRUNC           = 32
	CMSG_WAITALL         = 256
	CMSG_WAITFORONE      = 65536
	CSCM_RIGHTS          = 1
	CSHUT_RD             = 0
	CSHUT_RDWR           = 2
	CSHUT_WR             = 1
	CSOCK_CLOEXEC        = 524288
	CSOCK_DCCP           = 6
	CSOCK_DGRAM          = 2
	CSOCK_NONBLOCK       = 2048
	CSOCK_PACKET         = 10
	CSOCK_RAW            = 3
	CSOCK_RDM            = 4
	CSOCK_SEQPACKET      = 5
	CSOCK_STREAM         = 1
)

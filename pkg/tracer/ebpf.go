package tracer

type KProbeName string

const (
	TCPv4Connect       KProbeName = "kprobe/tcp_v4_connect"
	TCPv4ConnectReturn KProbeName = "kretprobe/tcp_v4_connect"
	TCPv6Connect       KProbeName = "kprobe/tcp_v6_connect"
	TCPv6ConnectReturn KProbeName = "kretprobe/tcp_v6_connect"

	TCPSendMsg     KProbeName = "kprobe/tcp_sendmsg"
	TCPCleanupRBuf KProbeName = "kprobe/tcp_cleanup_rbuf"
	TCPClose       KProbeName = "kprobe/tcp_close"

	UDPSendMsg       KProbeName = "kprobe/udp_sendmsg"
	UDPRecvMsg       KProbeName = "kprobe/udp_recvmsg"
	UDPRecvMsgReturn KProbeName = "kretprobe/udp_recvmsg"
)

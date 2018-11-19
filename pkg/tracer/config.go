package tracer

import "time"

type Config struct {
	// CollectTCPConns specifies whether the tracer should collect traffic statistics for TCP connections
	CollectTCPConns bool

	// CollectUDPConns specifies whether the tracer should collect traffic statistics for UDP connections
	CollectUDPConns bool

	// TraceIPv6Connections specifics whether the tracer should capture traffic for IPv6 TCP/UDP connections
	TraceIPv6Connections bool

	// UDPConnTimeout determines the length of traffic inactivity between two (IP, port)-pairs before declaring a UDP
	// connection as inactive.
	// Note: As UDP traffic is technically "connection-less", for tracking, we consider a UDP connection to be traffic
	//       between a source and destination IP and port.
	UDPConnTimeout time.Duration

	// TCPConnTimeout is like UDPConnTimeout, but for TCP connections. TCP connections are cleared when
	// the BPF module receives a tcp_close call, but TCP connections also age out to catch cases where
	// tcp_close is not intercepted for some reason.
	TCPConnTimeout time.Duration
}

// NewDefaultConfig enables traffic collection for all connection types
func NewDefaultConfig() *Config {
	return &Config{
		CollectTCPConns:      true,
		CollectUDPConns:      true,
		TraceIPv6Connections: isIPv6EnabledOnHost(),
		UDPConnTimeout:       30 * time.Second,
		TCPConnTimeout:       10 * time.Minute,
	}
}

// EnabledKProbes returns a map of kprobes that are enabled per config settings
func (c *Config) EnabledKProbes() map[KProbeName]struct{} {
	enabled := make(map[KProbeName]struct{}, 0)

	// Note: TCPv4Connect & TCPv4ConnectReturn are always included as they're needed for initialization
	// and can be disabled after field offset guessing has completed.
	enabled[TCPv4Connect] = struct{}{}
	enabled[TCPv4ConnectReturn] = struct{}{}

	if c.CollectTCPConns {
		enabled[TCPSendMsg] = struct{}{}
		enabled[TCPCleanupRBuf] = struct{}{}
		enabled[TCPClose] = struct{}{}
	}

	if c.CollectUDPConns {
		enabled[UDPRecvMsgReturn] = struct{}{}
		enabled[UDPRecvMsg] = struct{}{}
		enabled[UDPSendMsg] = struct{}{}
	}

	if c.TraceIPv6Connections {
		enabled[TCPv6Connect] = struct{}{}
		enabled[TCPv6ConnectReturn] = struct{}{}
	}

	return enabled
}

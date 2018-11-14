package tracer

import "time"

type Config struct {
	// CollectTCPConns specifies whether the tracer should collect traffic statistics for TCP connections
	CollectTCPConns bool
	// CollectUDPConns specifies whether the tracer should collect traffic statistics for UDP connections
	CollectUDPConns bool

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

// DefaultConfig enables traffic collection for all connection types
var DefaultConfig = &Config{
	CollectTCPConns: true,
	CollectUDPConns: true,
	UDPConnTimeout:  30 * time.Second,
	TCPConnTimeout:  10 * time.Minute,
}

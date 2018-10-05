package rpc

import (
	"net"

	"google.golang.org/grpc"
)

type RegisterFunc func(*grpc.Server)

// Serve listens on addr, registers all services provided by regs, and begins serving in a blocking way.
// It will block until there is a fatal error.
// If actualAddr is not nil, it will receive the actual listener addr if successful.
func Serve(addr string, regs []RegisterFunc, actualAddr chan<- net.Addr) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		if actualAddr != nil {
			close(actualAddr)
		}
		return err
	}
	if actualAddr != nil {
		actualAddr <- lis.Addr()
		close(actualAddr)
	}

	server := grpc.NewServer()
	for _, reg := range regs {
		reg(server)
	}

	return server.Serve(lis)
}

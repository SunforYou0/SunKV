package tcp

import (
	"context"
	"net"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn) //context.Context
	Close() error
}

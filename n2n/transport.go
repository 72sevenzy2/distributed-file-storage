package n2n

import (
	"log/slog"
	"net"
	"sync"
)

// Node represents a remote connection.
type Node interface {
	Close() error
}

type TCPTransport struct {
	ListenAddr string
	Logger     slog.Logger
	listener   net.Listener

	mu    sync.RWMutex // allows concurrent reads without blocking for node-to-node communication.
	Nodes map[string]net.Conn
}

func NewTCPTransport(Addr string) *TCPTransport {
	return &TCPTransport{
		ListenAddr: Addr,
	}
}

func (n *TCPTransport) ListenAndAccept() {
	ln, err := net.Listen("tcp", n.ListenAddr)
	if err != nil {
		n.Logger.Error("ERR", "TCP_ERROR", err)
		return
	}
	defer ln.Close()

	n.listener = ln
}

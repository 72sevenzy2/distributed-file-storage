package n2n

import (
	"log/slog"
	"net"
	"sync"
)

// Node represents a remote node which has established an connection to the server.
type TCPNode struct {
	Addr string
	Conn net.Conn
}

type TCPTransport struct {
	ListenAddr string
	Logger     slog.Logger
	Decoder    Decoder
	listener   net.Listener

	mu    sync.RWMutex // allows concurrent reads without blocking for node-to-node communication.
	Nodes map[string]*TCPNode
}

func NewTCPTransport(Addr string) *TCPTransport {
	return &TCPTransport{
		ListenAddr: Addr,
		Nodes:      make(map[string]*TCPNode),
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

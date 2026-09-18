package n2n

import (
	"log/slog"
	"net"
	"sync"
)

type Transport interface {
	ListenAndAccept()
}

type TCPTransport struct {
	ListenAddr  string
	Logger      slog.Logger
	Decoder     Decoder
	listener    net.Listener
	Handshakefn Handshaker

	mu    sync.RWMutex // allows concurrent reads without blocking for node-to-node communication.
	Nodes map[string]Node
}

func NewTCPTransport(Addr string) Transport {
	return &TCPTransport{
		ListenAddr: Addr,
		Decoder:    &GOBDecoder{},
		Nodes:      make(map[string]Node),
	}
}

func (n *TCPTransport) ListenAndAccept() {
	ln, err := net.Listen("tcp", n.ListenAddr)
	if err != nil {
		n.Logger.Error("ERR", "TCP_HANDSHAKE_ERR", err)
		return
	}
	defer ln.Close()

	n.listener = ln

	go n.acceptLoop()
}

func (n *TCPTransport) acceptLoop() {
	var msg []byte
	for {
		conn, err := n.Handshakefn.HandshakeFn()
		if err != nil {
			n.Logger.Error("ERR", "TCP_HANDSHAKE_ERR", err)
			return
		}

		defer conn.Close()
		n.Logger.Info("INCOMING_CONNECTION", "addr", conn.RemoteAddr().String())

		if err2 := n.Decoder.Decode(conn, msg); err2 != nil {
			conn.Close()
			n.Logger.Error("ERR", "TCP_DECODING_ERR", err2)
			return // drop connection upon unsuccessful payload.
		}
		n.mu.Lock()
		n.Nodes[conn.RemoteAddr().String()] = &TCPNode{
			Payload: msg,
			Addr:    conn.RemoteAddr(),
			Conn:    conn,
		}
		n.mu.Unlock()

		n.Logger.Info("PAYLOAD", "from", conn.RemoteAddr().String(), "payload", msg)
	}
}

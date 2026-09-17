package n2n

import (
	"log/slog"
	"net"
	"sync"
)

// Node represents a remote node which has established an connection to the server.
type TCPNode struct {
	Payload []byte
	Addr    string
	Conn    net.Conn
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
		Decoder:    &NOPDecoder{},
		Nodes:      make(map[string]*TCPNode),
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
		conn, err := n.listener.Accept()
		if err != nil {
			conn.Close()
			n.Logger.Error("ERR", "TCP_ACCEPT_ERR", err)
			return
		}
		defer conn.Close()
		n.Logger.Info("INCOMING_CONNECTION", "addr", conn.RemoteAddr().String())

		if err2 := n.Decoder.Decode(conn, msg); err2 != nil {
			conn.Close()
			n.Logger.Error("ERR", "TCP_DECODING_ERR", err2)
			return // drop connection upon unsuccessful payload.
		}
		n.Logger.Info("PAYLOAD", "from", conn.RemoteAddr().String(), "payload", msg)
	}
}

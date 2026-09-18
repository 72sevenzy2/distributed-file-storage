package n2n

import (
	"log/slog"
	"net"
	"sync"
)

// TCPNode represents the remote nodes connection methodology.
type TCPNode struct {
	Payload []byte
	Addr    string
	Conn    net.Conn
}

// Close() implements the Node interface.
func (n *TCPNode) Close() error {
	return n.Conn.Close()
}

type TCPTransport struct {
	ListenAddr string
	Logger     slog.Logger
	Decoder    Decoder
	listener   net.Listener

	mu    sync.RWMutex // allows concurrent reads without blocking for node-to-node communication.
	Nodes map[string]Node
}

func NewTCPTransport(Addr string) *TCPTransport {
	return &TCPTransport{
		ListenAddr: Addr,
		Decoder:    &NOPDecoder{},
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
		n.mu.Lock()
		n.Nodes[conn.RemoteAddr().String()] = &TCPNode{
			Payload: msg,
			Addr:    conn.RemoteAddr().String(),
			Conn:    conn,
		}
		n.mu.Unlock()

		n.Logger.Info("PAYLOAD", "from", conn.RemoteAddr().String(), "payload", msg)
	}
}

package n2n

import (
	"log/slog"
	"net"
)

type Transport interface {
	ListenAndAccept()
	Consume() <-chan TCPNode
}

type TCPTransport struct {
	ListenAddr  string
	Logger      slog.Logger
	Decoder     Decoder
	listener    net.Listener
	Handshakefn Handshaker
	TCPNodeCh   chan TCPNode
}

func NewTCPTransport(Addr string) Transport {
	return &TCPTransport{
		ListenAddr:  Addr,
		Decoder:     &GOBDecoder{},
		Logger:      *slog.Default(),
		Handshakefn: &TCPHandshake{},
		TCPNodeCh:   make(chan TCPNode),
	}
}

// Consume() implements Transport interface.
func (n *TCPTransport) Consume() <-chan TCPNode {
	return n.TCPNodeCh
}

func (n *TCPTransport) ListenAndAccept() {
	ln, err := net.Listen("tcp", n.ListenAddr)
	if err != nil {
		n.Logger.Error("ERR", "TCP_HANDSHAKE_ERR", err)
		return
	}

	n.listener = ln

	go n.acceptLoop()
}

func (n *TCPTransport) acceptLoop() {
	var msg []byte
	decodingErrCount := 0 // temporary spam prevention

	NodeDetails := TCPNode{}
	for {
		conn, err := n.Handshakefn.HandshakeFn(n.listener)
		if err != nil {
			n.Logger.Error("ERR", "TCP_HANDSHAKE_ERR", err)
			continue
		}

		n.Logger.Info("INCOMING_CONNECTION", "addr", conn.RemoteAddr().String())

		if err2 := n.Decoder.Decode(conn, &NodeDetails); err2 != nil {
			decodingErrCount++
			conn.Close()
			n.Logger.Error("ERR", "TCP_DECODING_ERR", err2)
			if decodingErrCount > 10 {
				return
			}
			continue
		}
		NodeDetails.Payload = msg
		NodeDetails.Addr = conn.RemoteAddr()
		NodeDetails.Conn = conn

		n.TCPNodeCh <- NodeDetails

		n.Logger.Info("PAYLOAD", "from", conn.RemoteAddr().String(), "payload", msg)
	}
}

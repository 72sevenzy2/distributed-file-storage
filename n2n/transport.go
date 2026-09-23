package n2n

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
)

type Transport interface {
	ListenAndAccept() error
	Consume() <-chan TCPNode
}

type TCPTransport struct {
	ListenAddr  string
	Logger      slog.Logger
	Decoder     Decoder
	listener    net.Listener
	Handshakefn Handshaker

	// OnPeer represents the Peers state upon establishing the connection to the server.
	// Allows for pre-flight validation before the connection is established.
	OnPeer    func(TCPNode) error
	TCPNodeCh chan TCPNode
}

func NOPOnPeer(v TCPNode) error { return nil }

func NewTCPTransport(Addr string, onPeer func(TCPNode) error) Transport {
	return &TCPTransport{
		ListenAddr:  Addr,
		Decoder:     &NOPDecoder{},
		Logger:      *slog.Default(),
		Handshakefn: &TCPHandshake{},
		OnPeer:      onPeer,
		TCPNodeCh:   make(chan TCPNode),
	}
}

// Consume() implements Transport interface.
func (n *TCPTransport) Consume() <-chan TCPNode {
	return n.TCPNodeCh
}

func (n *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", n.ListenAddr)
	if err != nil {
		return err
	}

	n.listener = ln
	go n.acceptLoop()
	return nil
}

func (n *TCPTransport) acceptLoop() {
	for {
		var conn net.Conn
		var err error
		NodeDetails := TCPNode{}

		defer func() {
			fmt.Println("dropping a connection:", err)
			conn.Close()
		}() // runs when the connection is closed.

		conn, err = n.Handshakefn.HandshakeFn(n.listener)
		if err != nil {
			if errors.Is(err, net.ErrClosed) { // if the network connection had been closed, return.
				conn.Close()
				return
			}

			n.Logger.Error("ERR", "TCP_HANDSHAKE_ERR", err)
			continue
		}

		if n.OnPeer != nil {
			if err = n.OnPeer(NodeDetails); err != nil {
				return
			}
		}

		go n.readLoop(conn, NodeDetails)
	}
}

func (n *TCPTransport) readLoop(conn net.Conn, peer TCPNode) error {
	// readloop
	for {
		n.Logger.Info("INCOMING_CONNECTION", "addr", conn.RemoteAddr().String())

		if err := n.Decoder.Decode(conn, &peer); err != nil {
			conn.Close()
			n.Logger.Error("ERR", "TCP_DECODING_ERR", err)
			continue
		}
		peer.Addr = conn.RemoteAddr()
		peer.Conn = conn

		n.TCPNodeCh <- peer
	}
}

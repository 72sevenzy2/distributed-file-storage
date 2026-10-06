package n2n

import (
	"net"
)

// Node represents a remote connection.
type Node interface {
	Send([]byte) error
	Close() error
	RemoteAddr() net.Addr
}

// TCPNode represents the remote nodes connection methodology.
type Peer struct {
	Conn     net.Conn
	outbound bool
}

func NewPeer(conn net.Conn, outbound bool) *Peer {
	return &Peer{
		Conn:     conn,
		outbound: outbound,
	}
}

// Send() implements the Node interface.
func (n *Peer) Send(b []byte) error {
	_, err := n.Conn.Write(b)
	return err
}

// RemoteAddr implements the Node interface, and will
// return the remote address of the underlying connection.
func (n *Peer) RemoteAddr() net.Addr {
	return n.Conn.RemoteAddr()
}

// Close() implements the Node interface.
func (n *Peer) Close() error {
	if n.Conn == nil { // safeguards against nil connections if it were closed elsewhere.
		return nil
	}

	return n.Conn.Close()
}

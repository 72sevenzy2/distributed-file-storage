package n2n

import (
	"net"
)

// Node represents a remote connection.
type Node interface {
	Close() error
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

// Close() implements the Node interface.
func (n *Peer) Close() error {
	if n.Conn == nil { // safeguards against nil connections if it were closed elsewhere.
		return nil
	}

	return n.Conn.Close()
}

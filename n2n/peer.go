package n2n

import (
	"net"
)

// Node represents a remote connection.
type Node interface {
	Close() error
}

// TCPNode represents the remote nodes connection methodology.
type TCPNode struct {
	Payload []byte
	Addr    net.Addr
	Conn    net.Conn
}

// Close() implements the Node interface.
func (n *TCPNode) Close() error {
	if n.Conn == nil { // safeguards against nil connections if it were closed elsewhere.
		return nil
	}

	return n.Conn.Close()
}

package n2n

import (
	"errors"
	"net"
)

type Handshaker interface {
	HandshakeFn() (net.Conn, error)
}

var TCPHandshakeErr = errors.New("TCP handshake error")

func (n *TCPTransport) HandshakeFn() (net.Conn, error) {
	conn, err := n.listener.Accept()
	if err != nil {
		return nil, TCPHandshakeErr
	}
	return conn, nil
}

package n2n

import (
	"net"
)

type Handshaker interface {
	HandshakeFn(net.Listener) (net.Conn, error)
}

type TCPHandshake struct{}

func (n *TCPHandshake) HandshakeFn(ln net.Listener) (net.Conn, error) {
	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

package n2n

import (
	"errors"
	"fmt"
	"net"
)

type Handshaker interface {
	HandshakeFn(net.Listener) (net.Conn, error)
}

var TCPHandshakeErr = errors.New("TCP handshake error")

type TCPHandshake struct{}

func (n *TCPHandshake) HandshakeFn(ln net.Listener) (net.Conn, error) {
	conn, err := ln.Accept()
	if errors.Is(err, net.ErrClosed) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %v", TCPHandshakeErr, err)
	}
	return conn, nil
}

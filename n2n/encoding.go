package n2n

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(r io.Reader, v *TCPNode) error
}

type GOBDecoder struct{}

func (g *GOBDecoder) Decode(r io.Reader, v *TCPNode) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(&v.Payload)
}

type NOPDecoder struct{}

func (n *NOPDecoder) Decode(r io.Reader, v *TCPNode) error {
	buf := make([]byte, 1024)

	n1, err := r.Read(buf)
	if err != nil {
		return err
	}

	v.Payload = buf[:n1]
	return nil
}

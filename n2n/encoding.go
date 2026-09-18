package n2n

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(r io.Reader, v any) error
}

type GOBDecoder struct{}

func (g *GOBDecoder) Decode(r io.Reader, v any) error {
	dec := gob.NewDecoder(r)
	return dec.Decode(v)
}

type NOPDecoder struct{}

func (n *NOPDecoder) Decode(r io.Reader, v any) error {
	_, err := r.Read(v.([]byte))
	return err
}

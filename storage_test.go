package main

import (
	"bytes"
	"io"
	"testing"
)

func TestStorage(t *testing.T) {
	storeOps := StorageOpts{
		PathTransformFunc: TransformPathFunc,
	}
	store := NewStorage(storeOps)
	key := "somekey"
	bytesData := []byte("some jpeg")

	if err := store.writeToStream(key, bytes.NewReader(bytesData)); err != nil {
		t.Error(err)
	}

	r, err := store.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)
	if string(b) != string(bytesData) {
		t.Errorf("invalid bytes read, have %s, want %s", string(b), string(bytesData))
	}
}

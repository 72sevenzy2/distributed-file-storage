package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	store := newStore()
	key := "somekey"
	bytesData := []byte("some jpeg")

	if err := store.writeToStream(key, bytes.NewReader(bytesData)); err != nil {
		t.Error(err)
	}

	if err := store.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStorage(t *testing.T) {
	store := newStore()
	defer clearAll(t, store)

	key := "somekey"
	bytesData := []byte("some jpeg")

	if err := store.writeToStream(key, bytes.NewReader(bytesData)); err != nil {
		t.Error(err)
	}

	if ok := store.Exists(key); !ok {
		t.Error("unknown key:", key)
	}

	r, err := store.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(bytesData) {
		t.Errorf("invalid bytes read, have %s, want %s", string(b), string(bytesData))
	}
	fmt.Println(string(b))
}

// helper utils
func newStore() *Storage {
	opts := StorageOpts{
		PathTransformFunc: TransformPathFunc,
	}
	return NewStorage(opts)
}

// for clearing directories after each test run.
func clearAll(t *testing.T, s *Storage) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

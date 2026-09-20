package main

import (
	"bytes"
	"testing"
)

func TestStorage(t *testing.T) {
	storeOps := StorageOpts{
		PathTransformFunc: TransformPathFunc,
	}
	store := NewStorage(storeOps)

	if err := store.writeToStream("somekey", bytes.NewReader([]byte("some jpeg"))); err != nil {
		t.Error(err)
	}
}

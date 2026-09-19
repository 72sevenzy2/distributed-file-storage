package main

import "io"

type PathTransformFunc func(string) string

type StorageOpts struct {
	PathTransformFunc PathTransformFunc
}

type Storage struct {
	StorageOpts
}

func NewStorage(s StorageOpts) *Storage {
	return &Storage{
		StorageOpts: s,
	}
}

func (s *Storage) writeToStream(key string, r io.Reader) error {
	return nil
}

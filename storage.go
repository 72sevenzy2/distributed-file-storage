package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

func TransformPathFunc(key string) string {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		src, dist := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[src:dist]
	}

	return strings.Join(paths, "/")
}

type PathTransformFunc func(string) string

type StorageOpts struct {
	PathTransformFunc PathTransformFunc
}

func DefaultPathTransformFunc(v string) string {
	return v
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
	path := s.PathTransformFunc(key)
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	filename := "some file name"
	f, err := os.Create(path + "/" + filename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Println("number of bytes written to disk:", n)

	return nil
}

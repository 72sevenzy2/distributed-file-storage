package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

func TransformPathFunc(key string) pathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		src, dist := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[src:dist]
	}

	return pathKey{
		FileName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

type PathTransformFunc func(string) pathKey

type StorageOpts struct {
	PathTransformFunc PathTransformFunc
}

func DefaultPathTransformFunc(v string) pathKey {
	return pathKey{}
}

type pathKey struct {
	FileName string
	Original string
}

func (p pathKey) Filename() string {
	return fmt.Sprintf("%s/%s", p.FileName, p.Original)
}

type Storage struct {
	StorageOpts
}

func NewStorage(s StorageOpts) *Storage {
	return &Storage{
		StorageOpts: s,
	}
}

func (s *Storage) Exists(key string) bool {
	path := s.PathTransformFunc(key)

	_, err := os.Stat(path.FileName)
	return errors.Is(err, fs.ErrNotExist)
}

func (s *Storage) Delete(key string) error {
	path := s.PathTransformFunc(key)
	return os.RemoveAll(path.FileName)
}

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Storage) readStream(key string) (io.ReadCloser, error) {
	path := s.PathTransformFunc(key)
	return os.Open(path.Filename())
}

func (s *Storage) writeToStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	if err := os.MkdirAll(path.FileName, os.ModePerm); err != nil {
		return err
	}
	filename := path.Filename()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Println("number of bytes written to disk:", n)

	return nil
}

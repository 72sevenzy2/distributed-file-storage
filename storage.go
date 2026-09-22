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

type PathTransformFunc func(string) pathKey

type pathKey struct {
	FileName string
	Original string
}

// DefaultRootFolder is the default parent directory of a file.
const DefaultRootFolder = "defaultRoot"

func DefaultPathTransformFunc(v string) pathKey {
	return pathKey{
		FileName: v,
		Original: v,
	}
}

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

type StorageOpts struct {
	// Root defines the root folder of the nested folders or files.
	Root              string
	PathTransformFunc PathTransformFunc
}

type Storage struct {
	StorageOpts
}

func NewStorage(s StorageOpts) *Storage {
	if s.PathTransformFunc == nil {
		s.PathTransformFunc = DefaultPathTransformFunc
	}
	if s.Root == "" {
		s.Root = DefaultRootFolder
	}

	return &Storage{
		StorageOpts: s,
	}
}

func (s *Storage) Exists(key string) bool {
	path := s.PathTransformFunc(key)
	pathWithRoot := fmt.Sprintf("%s/%s", s.Root, path.Filename())
	fmt.Println(pathWithRoot)

	_, err := os.Stat(pathWithRoot)
	ok := errors.Is(err, fs.ErrNotExist)
	if ok {
		return false
	}
	return true
}

func (s *Storage) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Storage) Delete(key string) error {
	path := s.PathTransformFunc(key)
	defer func() { fmt.Println("deleted path from disk:", path) }()

	pathWithRoot := fmt.Sprintf("%s/%s", s.Root, path.FirstFilepath())
	return os.RemoveAll(pathWithRoot)
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
	pathWithRoot := fmt.Sprintf("%s/%s", s.Root, path.Filename())
	return os.Open(pathWithRoot)
}

func (s *Storage) writeToStream(key string, r io.Reader) error {
	path := s.PathTransformFunc(key)
	pathWithRoot := fmt.Sprintf("%s/%s", s.Root, path.FileName)
	if err := os.MkdirAll(pathWithRoot, os.ModePerm); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/%s", s.Root, path.Filename())
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	fmt.Printf("number of bytes written to disk: %d, to path %s\n", n, filename)

	return nil
}

// FirstFilepath returns the first of the nested directory in .FileName
func (p pathKey) FirstFilepath() string {
	path := strings.Split(p.FileName, "/")
	if len(path) == 0 {
		return ""
	}
	return path[0]
}

func (p pathKey) Filename() string {
	return fmt.Sprintf("%s/%s", p.FileName, p.Original)
}

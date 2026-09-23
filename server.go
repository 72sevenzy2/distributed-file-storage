package main

import (
	"github.com/72sevenzy2/file-storage/n2n"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         n2n.Transport
}

type FileServer struct {
	FileServerOpts
	store *Storage
}

func NewFileServer(fs FileServerOpts) *FileServer {
	storeOpts := StorageOpts{
		Root:              fs.StorageRoot,
		PathTransformFunc: fs.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: fs,
		store:          NewStorage(storeOpts),
	}
}

func (fs *FileServer) Run() error {
	if err := fs.Transport.ListenAndAccept(); err != nil {
		return err
	}
	return nil
}

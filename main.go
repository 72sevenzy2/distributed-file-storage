package main

import (
	"fmt"
	"time"

	"github.com/72sevenzy2/file-storage/n2n"
)

func main() {
	tcpOpts := n2n.TCPTransportOpts{
		ListenAddr:  ":9000",
		Decoder:     &n2n.NOPDecoder{},
		Handshakefn: &n2n.TCPHandshake{},
	}

	tcp := n2n.NewTCPTransport(tcpOpts)

	fileStoreOpts := FileServerOpts{
		StorageRoot:       "some_root",
		PathTransformFunc: TransformPathFunc,
		Transport:         tcp,
	}

	server := NewFileServer(fileStoreOpts)

	go func() {
		time.Sleep(time.Second * 10)
		server.Stop()
	}()

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}
}

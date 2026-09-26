package main

import (
	"fmt"

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

	//go func() {
	//for {
	//msg := <-tcp.Consume()
	//		fmt.Println(msg)
	//}
	//}()

	server := NewFileServer(fileStoreOpts)

	if err := server.Run(); err != nil {
		fmt.Println(err)
		return
	}

	select {}
}

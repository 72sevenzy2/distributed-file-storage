package main

import (
	"fmt"

	"github.com/72sevenzy2/file-storage/n2n"
)

func main() {
	c := n2n.NewTCPTransport(":9000", n2n.NOPOnPeer)
	c.ListenAndAccept()

	go func() {
		for {
			msg := <-c.Consume()
			fmt.Println("message:\n", msg.Payload)
		}
	}()

	select {}
}

package main

import (
	"github.com/72sevenzy2/file-storage/n2n"
)

func main() {
	c := n2n.NewTCPTransport(":9000")
	c.ListenAndAccept()

	select {}
}

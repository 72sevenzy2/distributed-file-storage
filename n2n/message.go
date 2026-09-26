package n2n

// RPC holds the arbitrary data that is being sent over the connection.
type RPC struct {
	From    string
	Payload []byte
}

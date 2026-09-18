package n2n

// Node represents a remote connection.
type Node interface {
	Close() error
}

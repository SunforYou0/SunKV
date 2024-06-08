package resp

// redis protocol reply
type Reply interface {
	ToBytes() []byte
}

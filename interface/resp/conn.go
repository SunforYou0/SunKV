package resp

// redis protocol connection
type Connection interface {
	Write([]byte) error
	GetDBIndex() int
	SelectDB(int)
}

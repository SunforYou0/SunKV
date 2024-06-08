package reply

// const replies
type PongReply struct {
}

var pongbytes = []byte("+PONG\r\n")

func (p PongReply) ToBytes() []byte {
	return pongbytes
}

func MakePongReply() *PongReply {
	//
	return &PongReply{}
}

type OkReply struct {
}

var okBytes = []byte("+OK\r\n")

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

var nullBulkByte = []byte("$-1\r\n")

// empty str reply
type NullBulkReply struct {
}

func (n NullBulkReply) ToBytes() []byte {
	return nullBulkByte
}
func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

type EmptyMultiBulkReply struct {
}

var emptyMultiBulkBytes = []byte("*0\r\n")

func (e EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct {
}

var noBytes = []byte("")

func (n NoReply) ToBytes() []byte {
	return noBytes
}

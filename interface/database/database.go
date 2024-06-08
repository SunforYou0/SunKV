package database

import "sunkv/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply //执行指令
	Close()
	AfterClientClose(c resp.Connection)
}
type DataEntity struct {
	// redis 内部数据结构
	Data interface{}
}

package database

import (
	"sunkv/interface/resp"
	"sunkv/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegistCommand("ping", Ping, 1)
}

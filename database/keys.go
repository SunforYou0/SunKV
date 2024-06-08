package database

import (
	"sunkv/interface/resp"
	"sunkv/lib/utils"
	"sunkv/lib/wildcard"
	"sunkv/resp/reply"
)

// DEL
func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	deleted := db.Removes(keys...)
	if deleted > 0 {
		db.addAof(utils.ToCmdLine2("del", args...))
	}
	return reply.MakeIntReply(int64(deleted))
}

// exists
func execExists(db *DB, args [][]byte) resp.Reply {
	res := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exsit := db.GetEntity(key)
		if exsit {
			res++
		}
	}
	return reply.MakeIntReply(res)
}

// flushDB
func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	db.addAof(utils.ToCmdLine2("flushdb", args...))
	return reply.MakeOkReply()
}

// TYPE (TYPE K1)
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		switch entity.Data.(type) {
		case []byte:
			return reply.MakeStatusReply("string")
		default:
			return reply.UnknownErrReply{}
		}
	}
	return reply.MakeStatusReply("none")
}

// RENAME k1 k2
func execRename(db *DB, args [][]byte) resp.Reply {
	k1 := string(args[0])
	k2 := string(args[1])
	val, exist := db.GetEntity(k1)
	if !exist {
		reply.MakeErrReply("no such key")
	}
	db.PutEntity(k2, val)
	db.Remove(k1)
	db.addAof(utils.ToCmdLine2("rename", args...))
	return reply.MakeOkReply()
}

// RENAMENX
func execRenamenx(db *DB, args [][]byte) resp.Reply {
	k1 := string(args[0])
	k2 := string(args[1])
	_, exist2 := db.GetEntity(k2)
	if exist2 {
		return reply.MakeIntReply(0)
	}
	val, exist := db.GetEntity(k1)
	if !exist {
		return reply.MakeErrReply("no such key")
	}
	db.PutEntity(k2, val)
	db.Remove(k1)
	db.addAof(utils.ToCmdLine2("renamenx", args...))
	return reply.MakeIntReply(1)
}

// keys *
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	res := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			res = append(res, []byte(key))
		}
		return true
	})
	return reply.MakeMultiBulkReply(res)
}
func init() {
	RegistCommand("DEL", execDel, -2)
	RegistCommand("EXISTS", execExists, -2)
	RegistCommand("flushdb", execFlushDB, -1) // ignore other args
	RegistCommand("type", execType, 2)
	RegistCommand("rename", execRename, 3)
	RegistCommand("renamenx", execRenamenx, 3)
	RegistCommand("keys", execKeys, 2) // keys *
}

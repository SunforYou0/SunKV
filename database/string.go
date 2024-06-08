package database

import (
	"fmt"
	"sunkv/interface/database"
	"sunkv/interface/resp"
	"sunkv/lib/utils"
	"sunkv/resp/reply"
)

func (db *DB) getAsString(key string) ([]byte, reply.ErrorReply) {
	entity, ok := db.GetEntity(key)
	if !ok {
		return nil, nil
	}
	fmt.Printf("type of entity.Data:%T\n", entity.Data)
	bytes, ok := entity.Data.([]byte)
	if !ok {

		return nil, &reply.WrongTypeErrReply{}
	}
	return bytes, nil
}

// GET
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	bytes, err := db.getAsString(key)
	if err != nil {
		return err
	}
	if bytes == nil {
		return &reply.NullBulkReply{}
	}
	return reply.MakeBulkReply(bytes)
}

// SET k v
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]
	// 这一句感觉应该放到PutEntity里面
	entity := &database.DataEntity{
		Data: val,
	}
	db.PutEntity(key, entity)
	db.addAof(utils.ToCmdLine2("set", args...))
	db.data.ForEach(func(key string, val interface{}) bool {
		fmt.Printf("key:%s,val:%s\n", key, val.(*database.DataEntity).Data.([]byte))
		return true
	})
	return reply.MakeOkReply()
}

// SETNX
func execSetnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := string(args[1])
	// 这一句感觉应该放到PutEntity里面
	entity := &database.DataEntity{
		Data: val,
	}
	res := db.PutIfAbsent(key, entity)
	db.addAof(utils.ToCmdLine2("setnx", args...))
	return reply.MakeIntReply(int64(res))
}

// GETSET
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := string(args[1])
	entity, exist := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{Data: val})
	if !exist {
		return reply.MakeNullBulkReply()
	}
	db.addAof(utils.ToCmdLine2("getset", args...))
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// STRLEN
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	l := len(entity.Data.([]byte))
	return reply.MakeIntReply(int64(l))
}
func init() {
	RegistCommand("Set", execSet, -3)
	RegistCommand("Get", execGet, 2)
	RegistCommand("GetSet", execGetSet, 3)
	RegistCommand("StrLen", execStrLen, 2)
	RegistCommand("Setnx", execSetnx, 3)
}

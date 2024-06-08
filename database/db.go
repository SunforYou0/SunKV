package database

import (
	"strings"
	"sunkv/datastruct/dict"
	"sunkv/interface/database"
	"sunkv/interface/resp"
	"sunkv/resp/reply"
)

type DB struct {
	index int
	// key -> DataEntity
	data   dict.Dict
	addAof func(line CmdLine)
}
type CmdLine = [][]byte

// kv db 的指令的实现
type ExecFunc func(db *DB, args [][]byte) resp.Reply

func makeDB() *DB {
	db := &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {},
	}
	return db
}
func (db *DB) Exec(c resp.Connection, cmds CmdLine) resp.Reply {
	//ping set setnx
	cmdName := strings.ToLower(string(cmds[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command:" + cmdName)
	}
	if !validateArity(cmd.arity, cmds) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	f := cmd.exector
	// 第一个参数已经没用了，它不是具体函数的参数
	return f(db, cmds[1:])
}
func validateArity(arity int, cmdArgs [][]byte) bool {
	return true
}

/* ---- data Access ----- */

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {

	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity a DataEntity into DB
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists edit an existing DataEntity
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExist(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()

}

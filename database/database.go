package database

import (
	"strconv"
	"strings"
	"sunkv/aof"
	"sunkv/config"
	"sunkv/interface/resp"
	"sunkv/lib/logger"
	"sunkv/resp/reply"
)

type Database struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewDatabase() *Database {
	db := &Database{}
	//一个db就是一个sync.map
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	db.dbSet = make([]*DB, config.Properties.Databases)
	for i := range db.dbSet {
		d := makeDB()
		d.index = i
		db.dbSet[i] = d
	}
	if config.Properties.AppendOnly {
		handler, err := aof.NewAofHandler(db)
		if err != nil {
			panic(err)
		}

		db.aofHandler = handler
		for _, curDb := range db.dbSet {
			//for 循环闭包问题
			singleDB := curDb
			singleDB.addAof = func(line CmdLine) {
				db.aofHandler.AddAof(singleDB.index, line)
			}
			/*
				curDB.addAof = func(line CmdLine) {
						db.aofHandler.AddAof(curDB.index, line)
					// 此curDB引用外部变量curDB，curDB逃逸到堆上
					}
			*/
		}
	}

	return db
}

// 只有select在这一层执行,其余在具体的分库执行
func (db *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) == 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, db, args[1:])
	}
	dbIdx := client.GetDBIndex()
	curDB := db.dbSet[dbIdx]
	return curDB.Exec(client, args)
}

// select 2
func execSelect(client resp.Connection, db *Database, args [][]byte) resp.Reply {
	dbIdx, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}
	if dbIdx >= len(db.dbSet) {
		return reply.MakeErrReply("ERR DB index out of range")
	}
	client.SelectDB(dbIdx)
	return reply.MakeOkReply()
}
func (db *Database) Close() {
	//TODO implement me

}

func (db *Database) AfterClientClose(c resp.Connection) {
	//TODO implement me

}

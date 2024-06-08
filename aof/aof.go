package aof

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sunkv/config"
	"sunkv/interface/database"
	"sunkv/lib/logger"
	"sunkv/lib/utils"
	"sunkv/resp/connection"
	"sunkv/resp/parser"
	"sunkv/resp/reply"
)

const aofBufSize = 1 << 16

type CmdLine [][]byte
type payLoad struct {
	cmdLine CmdLine
	dbIdx   int
}
type AofHandler struct {
	database    database.Database
	aofChan     chan *payLoad // write cache
	aofFile     *os.File
	aofFileName string
	currentDB   int
}

// NewAofHandler
func NewAofHandler(db database.Database) (*AofHandler, error) {
	handler := &AofHandler{}
	handler.aofFileName = config.Properties.AppendFilename
	handler.database = db
	// load aof
	handler.LoadAof()
	aoffile, err := os.OpenFile(handler.aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		fmt.Println("aof file not found")
		return nil, err
	}
	handler.aofFile = aoffile
	//
	handler.aofChan = make(chan *payLoad, aofBufSize)
	go func() { handler.handleAof() }()
	return handler, nil
}

// Add payload->aofChan
func (handler *AofHandler) AddAof(dbIdx int, cmd CmdLine) {
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payLoad{
			cmdLine: cmd,
			dbIdx:   dbIdx,
		}
	}
}

// handleAof payload (set k v) <-aofChan
func (handler *AofHandler) handleAof() {
	//aofChan := handler.aofChan
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIdx != handler.currentDB {
			byteCmd := utils.ToCmdLine("select", strconv.Itoa(p.dbIdx))
			data := reply.MakeMultiBulkReply(byteCmd).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIdx
		}
		d := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(d)
		if err != nil {
			logger.Error(err)
			continue
		}
	}
}

// LoadAof`
func (handler *AofHandler) LoadAof() {
	curFile, err := os.Open(handler.aofFileName)
	if err != nil {
		logger.Error(err)
		return
	}
	defer curFile.Close()
	ch := parser.ParseStream(curFile)
	emptyConn := &connection.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				break
			} else {
				logger.Error(p.Err)
				continue
			}
		}
		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}
		mbReply, ok := p.Data.(reply.MultiBulkReply)
		//类型断言可以返回错误，强转不会
		if !ok {
			logger.Error("need multibulk reply")
			continue
		}
		rep := handler.database.Exec(emptyConn, mbReply.Args)
		if reply.IsErrReply(rep) {
			logger.Error(rep)
		}
	}
}

package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sunkv/lib/logger"
	"sunkv/lib/sync/atomic"
	"sunkv/lib/sync/wait"
	"sync"
	"time"
)

type EchoClient struct {
	// user client
	Conn    net.Conn
	Waiting wait.Wait
}

func (eC *EchoClient) Close() error {
	// wait 10 second if server still do sth
	eC.Waiting.WaitWithTimeout(10 * time.Second)
	_ = eC.Conn.Close()
	return nil
}

type EchoHandler struct {
	// TCP Echo
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeHandler() *EchoHandler {
	return &EchoHandler{}
}
func (eH *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	// reject any client while database closing
	if eH.closing.Get() {
		_ = conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	// record all active client
	eH.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Client close connection")
				eH.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (eH *EchoHandler) Close() error {
	logger.Info("handler shut down")
	eH.closing.Set(true)
	// close all client

	eH.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		// return true to get next key
		return true
	})
	return nil
}

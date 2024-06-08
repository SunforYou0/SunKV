package tcp

import (
	"context"
	"net"
	"os"
	"os/signal"
	"sunkv/interface/tcp"
	"sunkv/lib/logger"
	"sync"
	"syscall"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal)
	// get signal from os
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-sigChan
		switch signal {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	// closeChan will pass the progress stopped signal
	ListenAndServe(listener, handler, closeChan)
	return nil
}
func ListenAndServe(listener net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {
	go func() {
		// if progress killed, close listener and handler
		<-closeChan
		_ = listener.Close()
		_ = handler.Close()
	}()
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var waitDone sync.WaitGroup
	for true {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()

	}
	waitDone.Wait()
}

package main

import (
	"fmt"
	"os"
	"sunkv/config"
	"sunkv/lib/logger"
	"sunkv/resp/handler"
	"sunkv/tcp"
)

const configFile string = "sunkv.conf"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "sunkv",
		Ext:        "log",
		TimeFormat: "2012-01-02",
	})

	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}
	err := tcp.ListenAndServeWithSignal(
		&tcp.Config{Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port)},
		handler.MakeHandler())
	if err != nil {
		logger.Error(err)
	}
}

//*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
//*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n

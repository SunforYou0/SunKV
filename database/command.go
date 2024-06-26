package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	exector ExecFunc
	arity   int //arg num
}

func RegistCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: executor,
		arity:   arity,
	}
}

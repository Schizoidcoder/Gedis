package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc // 执行方法
	arity    int      // 参数数量
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
		arity:    arity,
	}
}

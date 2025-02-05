package database

import (
	"Gedis/datastruct/dict"
	"Gedis/interface/resp"
	"Gedis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict // 不同SyncDict是因为未来可能要换别的实现
}

type CmdLine = [][]byte

type ExecFunc func(db *DB, args [][]byte) resp.Reply

func makeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	//PING SET SETNX
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command " + cmdName)
	}
	//参数个数校验
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	return true

}

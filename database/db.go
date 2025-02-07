package database

import (
	"Gedis/datastruct/dict"
	"Gedis/interface/database"
	"Gedis/interface/resp"
	"Gedis/resp/reply"
	"strings"
)

type DB struct {
	index  int
	data   dict.Dict // 不同SyncDict是因为未来可能要换别的实现
	addAof func(CmdLine)
}

type CmdLine = [][]byte

type ExecFunc func(db *DB, args [][]byte) resp.Reply

func makeDB() *DB {
	db := &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {}, //先给DB一个空实现，因为如果在数据恢复的时候，你调用了set操作，此时还没有初始化好，如果你set那个key在AppendAof里面，那么后面初始化的时候就可能覆盖
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

// SET K V-> arity =3
// EXISTS k1 k2 arity =-2
func validateArity(arity int, cmdArgs [][]byte) bool {
	argNUm := len(cmdArgs)
	if arity >= 0 {
		return argNUm == arity
	}
	return argNUm >= -arity
}

// GetEntity Get k
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			deleted++
			db.Remove(key)
		}
	}
	return deleted
}

func (db *DB) Flush() {
	db.data.Clear()
}

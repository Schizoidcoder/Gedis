package database

import (
	"Gedis/datastruct/dict"
	"Gedis/interface/resp"
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

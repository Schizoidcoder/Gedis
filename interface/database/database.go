package database

import "Gedis/interface/resp"

type CmdLine = [][]byte

type Database interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct { //redis数据结构 string set list
	Data interface{}
}

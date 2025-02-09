package database

import (
	"Gedis/aof"
	"Gedis/config"
	"Gedis/interface/resp"
	"Gedis/lib/logger"
	"Gedis/resp/reply"
	"strconv"
	"strings"
)

type StandaloneDatabase struct {
	dbSet      []*DB
	aofHandler *aof.AofHandler
}

func NewStandaloneDatabase() *StandaloneDatabase {
	database := &StandaloneDatabase{}
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		db := makeDB()
		db.index = i
		database.dbSet[i] = db
	}
	if config.Properties.AppendOnly {
		aofHandler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = aofHandler
		for _, db := range database.dbSet {
			sdb := db
			sdb.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(sdb.index, line)
			}
			/*
				db.addAof = func(line CmdLine) {
					database.aofHandler.AddAof(db.index, line)//注意这里db会逃逸到堆上 而db一直会变
				}
			*/
		}
	}
	return database
}

// set k v
// get k
// select 2

func (database *StandaloneDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply("select")
		}
		return execSelect(client, database, args[1:])
	}
	dbIndex := client.GetDBIndex()
	db := database.dbSet[dbIndex]
	return db.Exec(client, args)
}

func (database *StandaloneDatabase) Close() {

}

func (database *StandaloneDatabase) AfterClientClose(c resp.Connection) {

}

//select 2
func execSelect(c resp.Connection, database *StandaloneDatabase, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid db index")
	}
	if dbIndex < 0 || dbIndex >= len(database.dbSet) {
		return reply.MakeErrReply("ERR db index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.MakeOkReply()
}

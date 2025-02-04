package handler

import (
	databaseface "Gedis/interface/database"
	"Gedis/lib/logger"
	"Gedis/lib/sync/atomic"
	"Gedis/resp/connection"
	"Gedis/resp/parser"
	"Gedis/resp/reply"
	"context"
	"io"
	"net"
	"strings"
	"sync"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

// closeClient 关闭其中一个客户端连接
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		//error
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed" + conn.RemoteAddr().String())
				return
			}
			//protocol error
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection error" + conn.RemoteAddr().String())
				return
			}
			continue
		}
		//exec

	}
}

// Close 关闭所有客户端连接
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key, value interface{}) bool {
			client := key.(*connection.Connection)
			_ = client.Close()
			return true
		},
	)
	r.db.Close()
	return nil
}

package cluster

import (
	"Gedis/config"
	database2 "Gedis/database"
	"Gedis/interface/database"
	"Gedis/interface/resp"
	"Gedis/lib/consistenthash"
	"Gedis/lib/logger"
	"Gedis/resp/reply"
	"context"
	"strings"

	pool "github.com/jolestar/go-commons-pool/v2"
)

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenthash.NodeMap     //集群存放id接口
	peerConnection map[string]*pool.ObjectPool //连接池
	db             database.Database
}

func MakeClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}
	nodes = append(nodes, cluster.self)
	cluster.nodes = nodes
	cluster.peerPicker.AddNode(cluster.nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	return cluster
}

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) (result resp.Reply) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		return reply.MakeErrReply("not supported command")
	}
	result = cmdFunc(cluster, client, args)
	return
}

func (c *ClusterDatabase) Close() {
	c.db.Close()
}

func (c *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	c.db.AfterClientClose(conn)
}

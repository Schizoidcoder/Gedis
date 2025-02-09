package cluster

import (
	"Gedis/config"
	database2 "Gedis/database"
	"Gedis/interface/database"
	"Gedis/interface/resp"
	"Gedis/lib/consistenthash"
	"context"

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
		pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}
	return cluster
}

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	//TODO implement me
	panic("implement me")
}

func (cluster *ClusterDatabase) Close() {
	//TODO implement me
	panic("implement me")
}

func (cluster *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	//TODO implement me
	panic("implement me")
}

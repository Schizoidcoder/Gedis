package cluster

import "Gedis/interface/resp"

type ClusterDatabase struct {
}

func MakeClusterDatabase() *ClusterDatabase {
	return &ClusterDatabase{}
}

func (c *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterDatabase) Close() {
	//TODO implement me
	panic("implement me")
}

func (c *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	//TODO implement me
	panic("implement me")
}

package cluster

import (
	"Gedis/resp/client"
	"context"
	"errors"
)

func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := cluster.peerConnection[peer] //拿到连接池
	if !ok {
		return nil, errors.New("connection not found")
	}
	object, err := pool.BorrowObject(context.Background()) //从池中拿到一个连接
	if err != nil {
		return nil, err
	}
	c, ok := object.(*client.Client) //类型断言
	if !ok {
		return nil, errors.New("connection not found")
	}
	return c, nil
}

func (cluster *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection not found")
	}
	return pool.ReturnObject(context.Background(), peerClient)

}

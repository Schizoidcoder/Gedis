package cluster

import (
	"Gedis/interface/resp"
	"Gedis/resp/reply"
)

//rename k1 k2
func Rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])
	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		//Todo:rename
		return reply.MakeErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, c, cmdArgs)
}

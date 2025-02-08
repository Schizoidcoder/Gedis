package consistenthash

import "hash/crc32"

type HashFunc func(data []byte) uint32
type NodeMap struct {
	hashFunc    HashFunc
	nodeHashs   []int          //这里是为了排序方便要用库的排序函数才用int 因为32位可以uint32直接转int 123434,1245
	nodeHashMap map[int]string // nodeHashMap[21412]= 'A'
}

func NewNodeMap(fn HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    fn,
		nodeHashMap: make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

func (m *NodeMap) IsEmpty() bool {
	return len(m.nodeHashMap) == 0
}

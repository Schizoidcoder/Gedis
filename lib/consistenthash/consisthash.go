package consistenthash

import (
	"hash/crc32"
	"sort"
)

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

func (m *NodeMap) AddNode(keys ...string) { //可能是节点名称或ip
	for _, key := range keys {
		if key == "" {
			continue
		}
		hash := int(m.hashFunc([]byte(key)))
		m.nodeHashs = append(m.nodeHashs, hash)
		m.nodeHashMap[hash] = key
	}
	sort.Ints(m.nodeHashs)
}

func (m *NodeMap) PickNode(keys string) string {
	if m.IsEmpty() {
		return ""
	}
	hash := int(m.hashFunc([]byte(keys)))
	idx := sort.Search(len(m.nodeHashs), func(i int) bool {
		return m.nodeHashs[i] >= hash
	})
	if idx == len(m.nodeHashs) {
		idx = 0
	}
	return m.nodeHashMap[m.nodeHashs[idx]]
}

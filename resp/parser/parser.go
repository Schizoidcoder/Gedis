package parser

import (
	"Gedis/interface/resp"
	"io"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool //解析器解析的是单行还是多行
	expectedArgsCount int  //期望长度
	msgType           byte
	args              [][]byte
	bulkLen           int64 //字节长度
}

func (s *readState) finished() bool {
	return s.expectedArgsCount > 0 && len(s.args) == s.expectedArgsCount
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}

func parse0(reader io.Reader, ch chan<- *Payload) {

}

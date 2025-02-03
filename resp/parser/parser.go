package parser

import (
	"Gedis/interface/resp"
	"bufio"
	"errors"
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
	bulkLen           int64 //bulkLen 存储的是批量字符串的字节长度（$len\r\n）len=bulkLen
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

func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) { //第一个是读数据，第二个是有无io错误，true就是有
	var msg []byte
	var err error
	if state.bulkLen == 0 { //1.\r\n切分
		msg, err = bufReader.ReadBytes('\n')
		if err != nil { // io错误
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
	} else { //2.之前读到了$数字，严格读取字符个数
		msg = make([]byte, state.bulkLen+2) //还要把\r\n读进来
		_, err = io.ReadFull(bufReader, msg)
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error:" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

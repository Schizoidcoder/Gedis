package reply

// UnknownErrReply 未知错误
type UnknownErrReply struct {
}

var unknownErrBytes = []byte("-Err unknown\r\n")

func (u UnknownErrReply) Error() string {
	return "Err unknown"
}

func (u UnknownErrReply) ToBytes() []byte {
	return unknownErrBytes
}

// ArgNumErrReply 参数数量错误
type ArgNumErrReply struct {
	Cmd string
}

func (r *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + r.Cmd + "' command"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{Cmd: cmd}
}

// SyntaxErrReply 语法错误
type SyntaxErrReply struct {
}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

var theSyntaxErrReply = &SyntaxErrReply{}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}

func (u SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (u SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

// WrongTypeErrReply represents operation against a key holding the wrong kind of value
type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")

// ToBytes marshals redis.Reply
func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key holding the wrong kind of value"
}

// ProtocolErr

// ProtocolErrReply represents meeting unexpected byte during parse requests
type ProtocolErrReply struct {
	Msg string
}

// ToBytes marshals redis.Reply
func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "'\r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error '" + r.Msg
}

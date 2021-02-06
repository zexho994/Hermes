package server

import "github.com/panjf2000/gnet"

type EchoServer struct {
	*gnet.EventServer
}

func (es *EchoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

package server

import (
	"flag"
	"github.com/panjf2000/gnet"
	"log"
	"strconv"
)

type tcpServer struct {
	*gnet.EventServer
}

func (es *tcpServer) onInitComplete(src gnet.Server) (action gnet.Action) {
	log.Printf("tcp server is listening on %s (multi-core: %t, loops:%d)\n", src.Addr.String(), src.Multicore, src.NumEventLoop)
	return
}

func (es *tcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	out = frame
	return
}

func InitTcpServer() {
	var port = flag.Int("port", 9002, "--port <pid>")
	// Example command: go run echo.go -port 9000 -multicore=true
	var isMulticore = flag.Bool("multicore", false, "--multicore true")
	flag.Parse()
	echo := new(tcpServer)
	_ = gnet.Serve(echo, "tcp://:"+strconv.Itoa(*port), gnet.WithMulticore(*isMulticore))
}

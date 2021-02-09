package server

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet"
	"log"
)

type tcpServer struct {
	*gnet.EventServer
}

func (es *tcpServer) onInitComplete(src gnet.Server) (action gnet.Action) {
	log.Printf("tcp server is listening on %s (multi-core: %t, loops:%d)\n", src.Addr.String(), src.Multicore, src.NumEventLoop)
	return
}

// 业务代码写在React里
func (es *tcpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	log.Println("new tcp msg")
	out = frame
	return
}


func InitTCPServer() {
	// Example command: go run echo.go -port 9000 -multicore=true
	var port int
	var multicore bool
	flag.IntVar(&port, "port", 8080, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.Parse()
	echo := new(tcpServer)
	log.Fatal(gnet.Serve(echo, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore)))
}

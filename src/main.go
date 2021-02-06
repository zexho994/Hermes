package main

import (
	"github.com/panjf2000/gnet"
	"log"
)
import "github.com/zouzhihao-994/Hermes/src/server"

func main() {
	echo := new(server.EchoServer)
	log.Fatal(gnet.Serve(echo, "tcp://:9000", gnet.WithMulticore(true)))
}

package main

import (
	"encoding/json"
	"github.com/zouzhihao-994/Hermes/codec"
	"io"
	"log"
	"net"
	"sync"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

// Accept 让服务端接收请求，后续处理任务交给协程处理
func (server *Server) Accept(lis net.Listener) {
	log.Println("-> accept")
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error", err)
			return
		}
		go server.ServerConn(conn)
	}
}

// ServerConn 处理连接的事情
func (server *Server) ServerConn(conn io.ReadWriteCloser) {
	log.Println("-> server conn")
	defer func() { _ = conn.Close() }()
	var opt codec.Option
	// 使用json格式解析option协议部分
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}
	// check magicnumber
	if opt.MagicNumber != codec.MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	// 确认编解码方式
	f := codec.NewCodecFuncMap[opt.CodeType]
	if f == nil {
		log.Printf("rpc server: invalic codec type %s", opt.CodeType)
	}
	// 使用指定的编解码方式进行后续解码
	server.serverCodec(f(conn))
}

// invalidRequest is a placeholder for response argv when error occurs
var invalidRequest = struct{}{}

func (server *Server) serverCodec(cc codec.Codec) {
	log.Println("server codec")
	sending := new(sync.Mutex) // make sure to send a complete response
	wg := new(sync.WaitGroup)  // wait until all request are handled
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // it's not possible to recover, so close the connection
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		// 处理请求
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

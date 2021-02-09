package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const(
	// Time allowed to write a message to the peer
	write_wait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pong_wait = 60 * time.Second
	

	ping_period = (pong_wait * 9) / 10 

	max_message_size = 512
)

var (
	newLine = []byte{'\n'}
	space = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type User struct{
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *User) readPump(){
	defer func(){
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(max_message_size)
	c.conn.SetReadDeadline(time.Now().Add(pong_wait))
	c.conn.SetPongHandler(func(string)error{c.conn.SetReadDeadline(time.Now().Add(pong_wait));return nil})
	for {
		_,message,err := c.conn.ReadMessage()
		if err != nil{
			if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure){
				log.Print("error: %v",err)
			}
		}
		message = bytes.TrimSpace(bytes.Replace(message,newLine,space,-1))
		c.hub.broadcast <- message
	}
}

func (c *User) writePump(){
	ticker := time.NewTicker(ping_period)
	defer func(){
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message,ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(write_wait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage,[]byte{})
				return
			}
			w,err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil{
				return
			}
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++{
				w.Write(newLine)
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil{
				return
			}
		case <- ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(write_wait))
			if err := c.conn.WriteMessage(websocket.PingMessage,nil); err != nil{
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter,r *http.Request){
	conn,err := upgrader.Upgrade(w,r,nil)
	if err != nil {
		log.Println(err)
		return
	}
	user := &User{hub: hub,conn: conn,send: make(chan []byte,256)}
	user.hub.register <- user
	go user.writePump()
	go user.readPump()
}
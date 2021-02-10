package server

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartServer() {
	var port int
	flag.IntVar(&port, "port", 8080, "-port <p>")
	flag.Parse()
	log.Printf("start http server, port = %d", port)

	setupRoutes()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// determine whether or not an incoming request from a different domain is allowed to connect
	upgrade.CheckOrigin = func(r *http.Request) bool { return true }

	wsc, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("client connected~")
	reader(wsc)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("client msg: %s", p)
		callMsg := "server call back msg"
		if err := conn.WriteMessage(messageType, []byte(callMsg)); err != nil {
			log.Fatal(err)
		}
	}
}

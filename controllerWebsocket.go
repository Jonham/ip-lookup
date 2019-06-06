package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var websocketConn []*websocket.Conn

func webscoketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		return
	}
	websocketConn = append(websocketConn, conn)

	for {
		t, msg, err := conn.ReadMessage()
		fmt.Println("t", t)
		fmt.Println("msg", msg)

		if err != nil {
			break
		}

		conn.WriteMessage(t, msg)
	}
}

func websocketServerSendHandler(c *gin.Context) {
	proxyEvent, _ := c.GetQuery("event")
	proxyMessage, _ := c.GetQuery("message")

	if len(websocketConn) > 0 {
		for i, conn := range websocketConn {
			conn.WriteMessage(1, []byte(proxyEvent+":"+proxyMessage+" clientID:"+strconv.Itoa(i)))
		}
	}
}

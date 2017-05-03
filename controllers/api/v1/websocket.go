package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WsHandler func
func WsHandler(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}

	for {
		json := struct {
			ID      int    `json:"id"`
			Message string `json:"message"`
			To      int    `json:"to"`
		}{}
		err := conn.ReadJSON(&json)
		if err != nil {
			fmt.Printf("ERR READ: %s", err.Error())
			break
		}

		conn.WriteJSON(json)
	}
}

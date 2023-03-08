package consumer

import (
	"github.com/gohade/hade/framework/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var Web *WebSocket

func Register(r *gin.Engine) error {
	Web = &WebSocket{Upgrader: &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}, Users: &sync.Map{}, Container: r.GetContainer()}

	r.GET("/websocket/:token", Web.CreateConn)

	return nil
}

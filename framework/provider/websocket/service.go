package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

type WebSocketServer struct {
	websocket websocket.Upgrader
	users sync.Map
}
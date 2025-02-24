package websocket

import (
    "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

func RegisterClient(conn *websocket.Conn) {
    clients[conn] = true
}

func BroadcastMessage(msg []byte) {
    for conn := range clients {
        err := conn.WriteMessage(websocket.TextMessage, msg)
        if err != nil {
            conn.Close()
            delete(clients, conn)
        }
    }
}

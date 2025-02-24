package websocket

import (
    "sync"
)

var (
    mutex sync.Mutex
)

func NotifyUsers(message string) {
    mutex.Lock()
    defer mutex.Unlock()

    BroadcastMessage([]byte(message))
}

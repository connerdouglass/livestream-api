package sockets

type Broadcaster interface {
	BroadcastToRoom(namespace, room, event string, args ...interface{}) bool
}

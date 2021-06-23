package services

type Broadcaster interface {
	BroadcastToRoom(namespace, room, event string, args ...interface{}) bool
}

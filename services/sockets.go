package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/godocompany/livestream-api/models"
	socketio "github.com/googollee/go-socket.io"
)

type SocketsService struct {
	Server          *socketio.Server
	StreamsService  *StreamsService
	TelegramService *TelegramService
}

func (s *SocketsService) Setup() {

	// Add handlers to the socket server
	s.Server.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		return nil
	})

	// When a socket disconnects
	s.Server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		fmt.Println("client disconnected: ", conn.RemoteAddr().String())
	})

	// Register all of the event handlers
	s.Server.OnEvent("/", "stream.join", s.OnStreamJoin)
	s.Server.OnEvent("/", "stream.leave", s.OnStreamLeave)

}

// Broadcast broadcasts a message to every member of a room
func (s *SocketsService) Broadcast(room, event string, args ...interface{}) bool {
	return s.Server.BroadcastToRoom("/", room, event, args...)
}

// StreamEnded broadcasts to every viewer of a stream that it has ended
func (s *SocketsService) StreamEnded(stream *models.Stream) {
	s.Broadcast(
		fmt.Sprintf("stream_%s", stream.Identifier),
		"stream.ended",
	)
}

// StreamStarted broadcasts to every viewer of a stream that it has started
func (s *SocketsService) StreamStarted(stream *models.Stream) {
	s.Broadcast(
		fmt.Sprintf("stream_%s", stream.Identifier),
		"stream.started",
	)
}

//====================================================================================================
// stream.join event handler
// Called when a viewer joins a stream
//====================================================================================================

type StreamJoinMsg struct {
	StreamID string `json:"stream_id"`
}

func (s *SocketsService) OnStreamJoin(conn socketio.Conn, data StreamJoinMsg) error {

	// Get the stream with the identifier
	stream, err := s.StreamsService.GetStreamByIdentifier(data.StreamID)
	if err != nil {
		return err
	}
	if stream == nil {
		return errors.New("stream not found")
	}

	// Join the room for the event
	conn.Join(
		fmt.Sprintf("stream_%s", stream.Identifier),
	)

	fmt.Println("joined stream: ", stream.Identifier, conn.RemoteAddr().String())

	go func() {
		for i := 0; i < 10; i++ {
			<-time.After(2 * time.Second)
			s.Broadcast(
				fmt.Sprintf("stream_%s", stream.Identifier),
				"chat.message",
				map[string]interface{}{
					"username": "guest_auto",
					"message":  "This is a great live stream!",
				},
			)
		}
	}()

	return nil

}

func (s *SocketsService) OnStreamLeave(conn socketio.Conn, data StreamJoinMsg) error {

	// Get the stream with the identifier
	stream, err := s.StreamsService.GetStreamByIdentifier(data.StreamID)
	if err != nil {
		return err
	}
	if stream == nil {
		return errors.New("stream not found")
	}

	// Leave the room for the event
	conn.Leave(
		fmt.Sprintf("stream_%s", stream.Identifier),
	)

	fmt.Println("left stream: ", stream.Identifier, conn.RemoteAddr().String())

	return nil

}

//====================================================================================================
// chat.message event handler
// Called when a viewer sends a message in the chat
//====================================================================================================

type ChatMsg struct {
	StreamID string `json:"stream_id"`
	Message  string `json:"message"`
}

func (s *SocketsService) OnChat(conn socketio.Conn, data ChatMsg) error {

	// Get the stream with the identifier
	stream, err := s.StreamsService.GetStreamByIdentifier(data.StreamID)
	if err != nil {
		return err
	}
	if stream == nil {
		return errors.New("stream not found")
	}

	// Get the name of the visitor
	// ...

	// Broadcast the message to the room
	s.Broadcast(
		fmt.Sprintf("stream_%s", stream.Identifier),
		"chat.message",
		map[string]interface{}{
			"username": "guest",
			"message":  data.Message,
		},
	)

	return nil

}

package services

import (
	"errors"
	"fmt"
	"sync"

	"github.com/godocompany/livestream-api/models"
	socketio "github.com/googollee/go-socket.io"
)

type SocketContext struct {
	User *TelegramUser
}

type SocketsService struct {
	Server               *socketio.Server
	StreamsService       *StreamsService
	TelegramService      *TelegramService
	ChatService          *ChatService
	streamChatBuffers    map[uint64]*LiveChatMessageBuffer
	streamChatBuffersMut sync.Mutex
}

func (s *SocketsService) Setup() {

	// Create the buffer
	s.streamChatBuffers = map[uint64]*LiveChatMessageBuffer{}

	// Add handlers to the socket server
	s.Server.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		conn.SetContext(SocketContext{})
		return nil
	})

	// When a socket disconnects
	s.Server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		fmt.Println("client disconnected: ", conn.RemoteAddr().String())
		conn.LeaveAll()
	})

	// Register all of the event handlers
	s.Server.OnEvent("/", "stream.join", s.OnStreamJoin)
	s.Server.OnEvent("/", "stream.leave", s.OnStreamLeave)
	s.Server.OnEvent("/", "telegram.auth", s.OnTelegramAuth)
	s.Server.OnEvent("/", "chat.message", s.OnChat)

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
// telegram.auth event handler
// Called when a user signs in with Telegram auth
//====================================================================================================

type TelegramAuthMsg struct {
	User TelegramUser `json:"user"`
}

func (s *SocketsService) OnTelegramAuth(conn socketio.Conn, data TelegramAuthMsg) error {

	// Get the context for the socket
	ctx, ok := conn.Context().(SocketContext)
	if !ok {
		return errors.New("invalid context for socket connection")
	}

	// Validate the user
	if !s.TelegramService.Verify(&data.User) {
		fmt.Println("Telegram user hash doesn't match!")
		return errors.New("invalid telegram user hash")
	}

	// Set the context
	ctx.User = &data.User

	// Return without error
	return nil

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

	// Emit all the buffered messages to the new viewer, so they don't open the page to
	// a completely empty live chat screen
	bufMsgs := s.copyChatMsgBuffer(stream.ID)
	for _, msg := range bufMsgs {
		conn.Emit(
			"chat.message",
			map[string]interface{}{
				"username":  msg.User.Username,
				"photo_url": msg.User.PhotoUrl,
				"message":   msg.Message,
			},
		)
	}

	fmt.Println("joined stream: ", stream.Identifier, conn.RemoteAddr().String())

	// Update the viewer count
	go s.StreamsService.UpdateViewerCount(
		stream,
		s.Server.RoomLen(
			"/",
			fmt.Sprintf("stream_%s", stream.Identifier),
		),
	)

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

	// Update the viewer count
	go s.StreamsService.UpdateViewerCount(
		stream,
		s.Server.RoomLen(
			"/",
			fmt.Sprintf("stream_%s", stream.Identifier),
		),
	)

	fmt.Println("left stream: ", stream.Identifier, conn.RemoteAddr().String())

	return nil

}

//====================================================================================================
// chat.message event handler
// Called when a viewer sends a message in the chat
//====================================================================================================

type ChatMsg struct {
	StreamID string       `json:"stream_id"`
	Message  string       `json:"message"`
	User     TelegramUser `json:"user"`
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

	// // Get the telegram user from the context
	// ctx, ok := conn.Context().(SocketContext)
	// if !ok {
	// 	return errors.New("invalid context for socket connection")
	// }

	// Validate the telegram user
	if !s.TelegramService.Verify(&data.User) {
		fmt.Println("Verification failed!")
		// return errors.New("invalid Telegram user hash")
	}

	// Check if the user is muted in chat
	muted, err := s.ChatService.IsUserMuted(data.User.Username)
	if err != nil {
		return err
	}
	if muted {
		return errors.New("user is muted in chat")
	}

	// Broadcast the message to the room
	s.Broadcast(
		fmt.Sprintf("stream_%s", stream.Identifier),
		"chat.message",
		map[string]interface{}{
			"username":  data.User.Username,
			"photo_url": data.User.PhotoUrl,
			"message":   data.Message,
		},
	)

	// Push the chat message to the buffer
	// Do it in a goroutine because we don't care about the result and we don't want to block
	// the socket handler just to do this task
	go s.pushChatMsgToBuffer(stream.ID, &data)

	return nil

}

func (s *SocketsService) pushChatMsgToBuffer(streamID uint64, msg *ChatMsg) {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	if !ok {
		buf = &LiveChatMessageBuffer{
			MaxLength: 10,
		}
		s.streamChatBuffers[streamID] = buf
	}

	// Unlock the buffer mutex since we have a pointer to what we need now
	s.streamChatBuffersMut.Unlock()

	// Push the message
	buf.Push(msg)

}

func (s *SocketsService) copyChatMsgBuffer(streamID uint64) []*ChatMsg {

	// Lock on the buffers
	s.streamChatBuffersMut.Lock()

	// Get the buffer for this stream identifier
	buf, ok := s.streamChatBuffers[streamID]
	s.streamChatBuffersMut.Unlock()
	if !ok {
		return nil
	}

	// Copy the values from the buffer
	return buf.GetCopy()

}

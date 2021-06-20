package sockets

import (
	"fmt"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type Server struct {
	SocketSrv *socketio.Server
}

func NewServer() *Server {
	return &Server{
		SocketSrv: socketio.NewServer(nil),
	}
}

// Setup prepares all of the handlers for the socket server and mounts it to the Gin router
func (s *Server) Setup(g *gin.RouterGroup) {

	// Mount the socket server to the Gin router
	g.Use(func(c *gin.Context) {
		s.SocketSrv.ServeHTTP(c.Writer, c.Request)
	})

	// Add handlers to the socket server
	s.SocketSrv.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		return nil
	})

}

func (s *Server) Run() {
	if err := s.SocketSrv.Serve(); err != nil {
		fmt.Println("Socket server error: ", err.Error())
	}
}

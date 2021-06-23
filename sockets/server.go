package sockets

import (
	"fmt"
	"net/http"
	"strings"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func simplifyOrigin(origin string) string {
	origin = strings.TrimRight(origin, "/")
	origin = strings.ToLower(origin)
	return origin
}

func checkOrigin(allowedOrigins []string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		origin = simplifyOrigin(origin)
		for _, allowed := range allowedOrigins {
			if simplifyOrigin(allowed) == origin {
				return true
			}
		}
		return false
	}
}

type Server struct {
	SocketSrv *socketio.Server
}

func NewServer(
	allowedOrigins []string,
) *Server {

	// The function for checking origins
	checkOriginFunc := checkOrigin(allowedOrigins)

	// Create the server
	server := &Server{
		SocketSrv: socketio.NewServer(&engineio.Options{
			Transports: []transport.Transport{
				&polling.Transport{
					CheckOrigin: checkOriginFunc,
				},
				&websocket.Transport{
					CheckOrigin: checkOriginFunc,
				},
			},
		}),
	}

	// Run the setup method
	server.setup()

	// Return the server
	return server

}

// Setup prepares all of the handlers for the socket server and mounts it to the Gin router
func (s *Server) setup() {

	// Add handlers to the socket server
	s.SocketSrv.OnConnect("/", func(conn socketio.Conn) error {
		fmt.Println("client connected: ", conn.RemoteAddr().String())
		return nil
	})

}

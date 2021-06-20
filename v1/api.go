package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/services"
	"github.com/godocompany/livestream-api/v1/hooks"
	"github.com/godocompany/livestream-api/v1/middleware"
)

// Server is the API server instance
type Server struct {
	AccountsService   *services.AccountsService
	AuthTokensService *services.AuthTokensService
	CreatorsService   *services.CreatorsService
	RtmpAuthService   *services.RtmpAuthService
	StreamsService    *services.StreamsService
}

// Setup mounts the API server to the given group
func (s *Server) Setup(g *gin.RouterGroup) {

	// Register middleware for all routes
	g.Use(middleware.CheckAuth(s.AuthTokensService))

	// Register all of the public hooks that require no authentication
	s.setupPublicHooks(g)

	// Register RTMP hooks (called by the privileged RTMP server with a passcode)
	s.setupRtmpHooks(g.Group("rtmp"))

	// Register authenticated hooks
	s.setupAuthenticatedHooks(g)

}

// setupPublicHooks mounts API hooks that are publicly accessible
func (s *Server) setupPublicHooks(g *gin.RouterGroup) {

	// Register public API routes
	g.POST("/auth/login", hooks.AuthLogin(
		s.AccountsService,
		s.AuthTokensService,
		s.CreatorsService,
	))
	g.POST("/creator/get-meta", hooks.GetCreatorMeta(
		s.CreatorsService,
		s.StreamsService,
	))
	g.POST("/stream/get-meta", hooks.GetStreamMeta(
		s.StreamsService,
	))

}

// setupRtmpHooks mounts API hooks used by the RTMP server during streaming
func (s *Server) setupRtmpHooks(g *gin.RouterGroup) {

	// Require the RTMP passcode for these hooks
	g.Use(middleware.RequireRtmpAuth(s.RtmpAuthService))

	// Register RTMP-only hooks here
	g.POST("/stream/get-config", hooks.RtmpGetStreamConfig(
		s.StreamsService,
	))
	g.POST("/stream/set-streaming", hooks.RtmpSetStreaming(
		s.StreamsService,
	))

}

// setupAuthenticatedHooks mounts API hooks that require account authentication
func (s *Server) setupAuthenticatedHooks(g *gin.RouterGroup) {

	// Require login for everything after this
	g.Use(middleware.RequireLogin())

	// Register authenticated API routes
	g.POST("/auth/whoami", hooks.AuthWhoAmI(
		s.AuthTokensService,
		s.CreatorsService,
	))
	g.POST("/studio/stream/set-status", hooks.StudioSetStreamStatus(
		s.AccountsService,
		s.CreatorsService,
		s.StreamsService,
	))
	g.POST("/studio/stream/get", hooks.StudioGetStream(
		s.AccountsService,
		s.CreatorsService,
		s.StreamsService,
	))
	g.POST("/studio/streams/list", hooks.StudioListStreams(
		s.CreatorsService,
		s.StreamsService,
	))

}

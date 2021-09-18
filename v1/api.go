package v1

import (
	"github.com/connerdouglass/livestream-api/services"
	"github.com/connerdouglass/livestream-api/v1/hooks"
	"github.com/connerdouglass/livestream-api/v1/middleware"
	"github.com/gin-gonic/gin"
)

// Server is the API server instance
type Server struct {
	PlatformTitle       string
	MainCreatorUsername string
	SiteConfigService   *services.SiteConfigService
	AccountsService     *services.AccountsService
	AuthTokensService   *services.AuthTokensService
	MembershipService   *services.MembershipService
	CreatorsService     *services.CreatorsService
	RtmpAuthService     *services.RtmpAuthService
	StreamsService      *services.StreamsService
	TelegramService     *services.TelegramService
	BrowserNotifier     *services.BrowserNotifier
	TelegramNotifier    *services.TelegramNotifier
	Notifier            services.Notifier
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
	g.POST("/app/get-state", hooks.AppState(
		s.PlatformTitle,
		s.MainCreatorUsername,
		s.TelegramService,
		s.BrowserNotifier,
	))
	g.POST("/auth/login", hooks.AuthLogin(
		s.AccountsService,
		s.AuthTokensService,
		s.MembershipService,
	))
	g.POST("/creator/get-meta", hooks.GetCreatorMeta(
		s.CreatorsService,
		s.StreamsService,
	))
	g.POST("/stream/get-meta", hooks.GetStreamMeta(
		s.StreamsService,
	))
	g.POST("/notifications/browser/state", hooks.BrowserNotificationsState(
		s.BrowserNotifier,
	))
	g.POST("/notifications/browser/register", hooks.BrowserNotificationsRegister(
		s.BrowserNotifier,
	))
	g.POST("/notifications/browser/update-sub", hooks.BrowserNotificationsUpdateSub(
		s.BrowserNotifier,
	))
	g.POST("/notifications/telegram/state", hooks.TelegramNotificationsState(
		s.TelegramNotifier,
	))
	g.POST("/notifications/telegram/update-sub", hooks.TelegramNotificationsUpdateSub(
		s.TelegramNotifier,
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
		s.MembershipService,
	))
	g.POST("/studio/members/add", hooks.StudioAddMember(
		s.AccountsService,
		s.CreatorsService,
		s.MembershipService,
	))
	g.POST("/studio/members/list", hooks.StudioListMembers(
		s.CreatorsService,
		s.MembershipService,
	))
	g.POST("/studio/stream/set-status", hooks.StudioSetStreamStatus(
		s.StreamsService,
		s.MembershipService,
		s.Notifier,
	))
	g.POST("/studio/stream/get", hooks.StudioGetStream(
		s.CreatorsService,
		s.StreamsService,
		s.MembershipService,
	))
	g.POST("/studio/stream/create", hooks.StudioCreateStream(
		s.CreatorsService,
		s.StreamsService,
		s.MembershipService,
	))
	g.POST("/studio/stream/update", hooks.StudioUpdateStream(
		s.CreatorsService,
		s.StreamsService,
		s.MembershipService,
	))
	g.POST("/studio/streams/list", hooks.StudioListStreams(
		s.CreatorsService,
		s.StreamsService,
		s.MembershipService,
	))

}

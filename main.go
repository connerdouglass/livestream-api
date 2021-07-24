package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/godocompany/livestream-api/models"
	"github.com/godocompany/livestream-api/services"
	v1 "github.com/godocompany/livestream-api/v1"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	//================================================================================
	// Create the database connection
	//================================================================================

	// Get the datbase driver for the database string
	dbDriver := ParseDatabaseDriver(os.Getenv("DB_URL"))
	if dbDriver == nil {
		log.Fatalln("Failed to create database driver. Check DB_URL environment variable")
	}

	// Create the database connection
	db, err := gorm.Open(dbDriver, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.Account{},
		&models.CreatorProfileMember{},
		&models.CreatorProfile{},
		&models.NotificationSubscriber{},
		&models.SiteConfig{},
		&models.Stream{},
	)

	//================================================================================
	// Create all the service instances
	//================================================================================

	// Create the rest of the services
	siteConfigService := &services.SiteConfigService{DB: db}
	telegramService := &services.TelegramService{
		BotAPIKey:   os.Getenv("TELEGRAM_BOT_API_KEY"),
		BotUsername: os.Getenv("TELEGRAM_BOT_USERNAME"),
	}
	accountsService := &services.AccountsService{DB: db}
	authTokensService := &services.AuthTokensService{
		DB:            db,
		SigningPepper: os.Getenv("AUTH_TOKEN_SIGNING_PEPPER"),
	}
	creatorsService := &services.CreatorsService{DB: db}
	rtmpAuthService := &services.RtmpAuthService{
		RtmpServerPasscode: os.Getenv("RTMP_SERVER_PASSCODE"),
	}
	streamsService := &services.StreamsService{DB: db}
	membershipService := &services.MembershipService{DB: db}
	notificationsService := &services.NotificationsService{
		DB:                db,
		SiteConfigService: siteConfigService,
		TelegramService:   telegramService,
	}

	//================================================================================
	// Listen on the Telegram bot channel
	//================================================================================

	go func() {
		if err := telegramService.Listen(); err != nil {
			fmt.Println("Telegram bot error: ", err.Error())
		}
	}()

	//================================================================================
	// Setup the Gin HTTP router
	//================================================================================

	// Create the Gin router
	r := gin.Default()

	// Configure CORS for the API
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = GetAllowedOrigins()
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Accept", "User-Agent", "Authorization")
	r.Use(cors.New(corsCfg))

	// Create the API instance
	api := &v1.Server{
		MainCreatorUsername:  os.Getenv("MAIN_CREATOR_USERNAME"),
		SiteConfigService:    siteConfigService,
		AccountsService:      accountsService,
		AuthTokensService:    authTokensService,
		CreatorsService:      creatorsService,
		MembershipService:    membershipService,
		RtmpAuthService:      rtmpAuthService,
		StreamsService:       streamsService,
		TelegramService:      telegramService,
		NotificationsService: notificationsService,
	}

	// Mount the API routes
	api.Setup(r.Group("v1"))

	// Create a mux to serve both the HTTP and Socket.IO servers
	mux := http.NewServeMux()
	mux.Handle("/", r)

	// Run the server
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Panicln(err)
	}

}

// GetAllowedOrigins gets the slice of allowed CORS origins
func GetAllowedOrigins() []string {

	// Get the list of origins allowed
	env, ok := os.LookupEnv("CORS_ALLOW_ORIGINS")
	if !ok {
		return []string{}
	}

	// Create the slice for it
	origins := []string{}

	// Split up the env value
	originsRaw := strings.Split(env, ",")
	for _, originRaw := range originsRaw {
		origin := strings.TrimSpace(originRaw)
		origins = append(origins, origin)
	}

	// Return the origins slice
	return origins

}

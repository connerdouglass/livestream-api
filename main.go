package main

import (
	"fmt"
	"log"
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
		&models.Stream{},
	)

	//================================================================================
	// Create all the service instances
	//================================================================================

	authTokensService := &services.AuthTokensService{
		DB:            db,
		SigningPepper: os.Getenv("AUTH_TOKEN_SIGNING_PEPPER"),
	}
	rtmpAuthService := &services.RtmpAuthService{
		RtmpServerPasscode: os.Getenv("RTMP_SERVER_PASSCODE"),
	}
	streamsService := &services.StreamsService{DB: db}

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
		AuthTokensService: authTokensService,
		RtmpAuthService:   rtmpAuthService,
		StreamsService:    streamsService,
	}

	// Mount the API routes
	api.Setup(r.Group("v1"))

	// Run the server
	if err := r.Run(); err != nil {
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

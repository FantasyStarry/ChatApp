package main

import (
	"chatapp/config"
	"chatapp/controllers"
	"chatapp/handlers"
	"chatapp/middleware"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRoutes() *gin.Engine {
	// Set Gin mode based on config
	if config.GlobalConfig.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware with config
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", strings.Join(config.GlobalConfig.CORS.AllowedOrigins, ","))
		c.Header("Access-Control-Allow-Methods", strings.Join(config.GlobalConfig.CORS.AllowedMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(config.GlobalConfig.CORS.AllowedHeaders, ","))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/profile", controllers.GetProfile)

		// Chat room routes
		protected.GET("/chatrooms", controllers.GetChatRooms)
		protected.POST("/chatrooms", controllers.CreateChatRoom)
		protected.GET("/chatrooms/:id", controllers.GetChatRoom)
		protected.GET("/chatrooms/:id/messages", controllers.GetChatRoomMessages)

		// WebSocket route
		protected.GET("/ws/:chatroom_id", handlers.HandleWebSocket)
	}

	return r
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	log.Printf("Starting %s v%s in %s mode", cfg.App.Name, cfg.App.Version, cfg.App.Environment)

	// Connect to database
	config.ConnectDatabase()

	// Run migrations
	config.MigrateDatabase()

	// Start WebSocket hub
	handlers.InitWebSocketUpgrader()
	go handlers.GlobalHub.Run()

	// Setup routes
	r := setupRoutes()

	// Start server
	log.Printf("Server starting on %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

package main

import (
	"chatapp/config"
	"chatapp/controllers"
	"chatapp/handlers"
	"chatapp/middleware"
	"chatapp/repository"
	"chatapp/service"
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

	// Initialize repositories
	userRepo := repository.NewUserRepository(config.DB)
	chatRoomRepo := repository.NewChatRoomRepository(config.DB)
	messageRepo := repository.NewMessageRepository(config.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	chatRoomService := service.NewChatRoomService(chatRoomRepo, userRepo)
	messageService := service.NewMessageService(messageRepo, userRepo, chatRoomRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	chatRoomController := controllers.NewChatRoomController(chatRoomService, messageService)

	// Initialize WebSocket hub with message service
	handlers.InitializeHub(messageService)

	// Public routes
	api := r.Group("/api")
	{
		api.POST("/login", authController.Login)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/profile", authController.GetProfile)

		// Chat room routes
		protected.GET("/chatrooms", chatRoomController.GetChatRooms)
		protected.POST("/chatrooms", chatRoomController.CreateChatRoom)
		protected.GET("/chatrooms/:id", chatRoomController.GetChatRoom)
		protected.GET("/chatrooms/:id/messages", chatRoomController.GetChatRoomMessages)

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

	// Setup routes (this initializes the hub)
	r := setupRoutes()

	// Start WebSocket hub after initialization
	handlers.InitWebSocketUpgrader()
	go handlers.GlobalHub.Run()

	// Start server
	log.Printf("Server starting on %s", cfg.Server.Port)
	if err := r.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

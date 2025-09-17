package main

import (
	"chatapp/config"
	"chatapp/models"
	"chatapp/utils"
	"log"
)

func seedData() {
	log.Println("Seeding database with test data...")

	// Create test users
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
		},
		{
			Username: "user1",
			Email:    "user1@example.com",
		},
		{
			Username: "user2",
			Email:    "user2@example.com",
		},
		{
			Username: "user3",
			Email:    "user3@example.com",
		},
	}

	// Hash passwords and create users
	for i := range users {
		hashedPassword, err := utils.HashPassword("password123")
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", users[i].Username, err)
			continue
		}
		users[i].Password = hashedPassword

		// Check if user already exists
		var existingUser models.User
		if err := config.DB.Where("username = ?", users[i].Username).First(&existingUser).Error; err != nil {
			// User doesn't exist, create it
			if err := config.DB.Create(&users[i]).Error; err != nil {
				log.Printf("Failed to create user %s: %v", users[i].Username, err)
			} else {
				log.Printf("Created user: %s", users[i].Username)
			}
		} else {
			log.Printf("User %s already exists", users[i].Username)
		}
	}

	// Create test chat rooms
	chatRooms := []models.ChatRoom{
		{
			Name:        "General",
			Description: "General discussion room",
			CreatedBy:   1, // admin user
		},
		{
			Name:        "Tech Talk",
			Description: "Technical discussions and programming",
			CreatedBy:   1, // admin user
		},
		{
			Name:        "Random",
			Description: "Random conversations and fun stuff",
			CreatedBy:   2, // user1
		},
	}

	for i := range chatRooms {
		var existingRoom models.ChatRoom
		if err := config.DB.Where("name = ?", chatRooms[i].Name).First(&existingRoom).Error; err != nil {
			// Room doesn't exist, create it
			if err := config.DB.Create(&chatRooms[i]).Error; err != nil {
				log.Printf("Failed to create chat room %s: %v", chatRooms[i].Name, err)
			} else {
				log.Printf("Created chat room: %s", chatRooms[i].Name)
			}
		} else {
			log.Printf("Chat room %s already exists", chatRooms[i].Name)
		}
	}

	// Create some sample messages
	messages := []models.Message{
		{
			Content:    "Welcome to the General chat room!",
			UserID:     1, // admin
			ChatRoomID: 1, // General
		},
		{
			Content:    "Hello everyone!",
			UserID:     2, // user1
			ChatRoomID: 1, // General
		},
		{
			Content:    "Let's discuss Go programming here",
			UserID:     1, // admin
			ChatRoomID: 2, // Tech Talk
		},
	}

	for i := range messages {
		if err := config.DB.Create(&messages[i]).Error; err != nil {
			log.Printf("Failed to create message: %v", err)
		} else {
			log.Printf("Created sample message in room %d", messages[i].ChatRoomID)
		}
	}

	log.Println("Database seeding completed!")
}

func main() {
	// Load configuration
	_, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	config.ConnectDatabase()

	// Run migrations
	config.MigrateDatabase()

	// Seed test data
	seedData()

	log.Println("Seed script completed successfully!")
}

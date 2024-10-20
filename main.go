package main

import (
	"fmt"
	"log"
	"time"

	"geospatial--backend/database"
	"geospatial--backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Initialize router
	r := gin.Default()

	// Configure CORS
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Vite's default port
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// Setup routes
	routes.SetupRoutes(r, db)

	// Start server
	if err := r.Run(":8087"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	fmt.Println("Server running at : 8087")
}

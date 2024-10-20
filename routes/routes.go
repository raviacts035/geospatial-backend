package routes

import (
	"geospatial--backend/handlers"
	"geospatial--backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {

	// r.Use(middleware.CORSMiddleware())
	// Public routes
	r.POST("/api/register", handlers.Register(db))
	r.POST("/api/login", handlers.Login(db))
	r.POST("/api/logout", handlers.Logout())

	// for testing purpose
	// r.GET("/api/geodata", handlers.ListGeoData(db))
	// r.POST("/api/geodata", handlers.CreateGeoData(db))
	// r.PUT("/api/geodata/:id", handlers.UpdateGeoData(db))
	// r.DELETE("/api/geodata/:id", handlers.DeleteGeoData(db))

	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/geodata", handlers.ListGeoData(db))
		auth.POST("/geodata", handlers.CreateGeoData(db))
		auth.PUT("/geodata/:id", handlers.UpdateGeoData(db))
		auth.DELETE("/geodata/:id", handlers.DeleteGeoData(db))
	}
}

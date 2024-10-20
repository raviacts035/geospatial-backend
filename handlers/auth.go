package handlers

import (
	"fmt"
	"net/http"
	"time"

	"geospatial--backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Request body : ", requestBody)

		var user models.User
		if username, ok := requestBody["username"].(string); ok {
			user.Email = username
		}

		if username, ok := requestBody["password"].(string); ok {
			user.Password = username
		}
		var existing_user models.User

		// Checking if user exists
		if err := db.Where("email = ?", user.Email).First(&existing_user).Error; err == nil {
			// If user exists, return an error
			c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
			return
		} else if err != gorm.ErrRecordNotFound {
			// Handle any other potential errors from the database query
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": gin.H{"email": user.Email, "id": user.ID}})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUser models.User
		var requestBody map[string]interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println("Request body : ", requestBody)

		if username, ok := requestBody["username"].(string); ok {
			loginUser.Email = username
		}

		if username, ok := requestBody["password"].(string); ok {
			loginUser.Password = username
		}

		var user models.User
		if err := db.Where("email = ?", loginUser.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := user.ComparePassword(loginUser.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		var auth_token string = user.GenerateToken()

		// Set cookie
		c.SetCookie(
			"auth_token",
			auth_token,
			int(time.Hour*24/time.Second), // 24 hours
			"/",
			"localhost",
			true, // Secure
			true, // HTTP only
		)

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": gin.H{"email": user.Email}, "access": auth_token})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			"auth_token",
			"",
			-1,
			"/",
			"localhost",
			true,
			true,
		)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

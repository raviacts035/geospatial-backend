package handlers

import (
	"net/http"
	"strconv"

	"geospatial--backend/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)

// CreateGeoDataInput represents the expected input format
type CreateGeoDataInput struct {
	Data datatypes.JSON
}

func CreateGeoData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateGeoDataInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid input format",
				"details": err.Error(),
			})
			return
		}

		userID, _ := c.Get("user_id")

		geoData := models.GeoData{
			UserID: userID.(uint),
			Data:   input.Data,
		}

		if err := db.Create(&geoData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create GeoData"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":  "GeoData created successfully",
			"geo_data": geoData,
		})
	}
}

func ListGeoData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var geoDataList []models.GeoData
		if err := db.Where("user_id = ?", userID).Find(&geoDataList).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch GeoData"})
			return
		}
		c.JSON(http.StatusOK, geoDataList)
	}
}

func UpdateGeoData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, _ := c.Get("user_id")

		var geoData models.GeoData
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&geoData).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "GeoData not found or unauthorized"})
			return
		}

		var input CreateGeoDataInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		geoData.Data = input.Data
		if err := db.Save(&geoData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update GeoData"})
			return
		}

		c.JSON(http.StatusOK, geoData)
	}
}

func DeleteGeoData(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		db = db.Debug()

		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid ID format",
			})
			return
		}

		userID, _ := c.Get("user_id")

		var geoData models.GeoData
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&geoData).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "GeoData not found or unauthorized",
					"cause": err,
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
				"cause": err,
			})
			return
		}

		if err := db.Unscoped().Delete(&geoData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to delete GeoData",
				"cause": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "GeoData deleted successfully",
			"id":      id,
		})
	}
}

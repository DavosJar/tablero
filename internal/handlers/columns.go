package handlers

import (
	"tablero/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateColumnRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateColumnRequest struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}

func GetColumns(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var columns []models.Column
		if err := db.Order("\"order\"").Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).Find(&columns).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch columns"})
			return
		}
		c.JSON(200, columns)
	}
}

func CreateColumn(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateColumnRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		column := models.Column{
			Name: req.Name,
		}

		if err := db.Create(&column).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create column"})
			return
		}

		c.JSON(201, column)
	}
}

func UpdateColumn(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req UpdateColumnRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		var column models.Column
		if err := db.First(&column, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "column not found"})
			return
		}

		if req.Name != "" {
			column.Name = req.Name
		}
		if req.Order != 0 {
			column.Order = req.Order
		}

		if err := db.Save(&column).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to update column"})
			return
		}

		c.JSON(200, column)
	}
}

func DeleteColumn(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var column models.Column
		if err := db.First(&column, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "column not found"})
			return
		}

		// Check if column has tasks
		var taskCount int64
		if err := db.Model(&models.Task{}).Where("column_id = ?", id).Count(&taskCount).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to check tasks"})
			return
		}

		if taskCount > 0 {
			c.JSON(400, gin.H{"error": "cannot delete column with tasks"})
			return
		}

		if err := db.Delete(&column).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to delete column"})
			return
		}

		c.JSON(200, gin.H{"message": "column deleted"})
	}
}

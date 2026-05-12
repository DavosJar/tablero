package handlers

import (
	"tablero/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     *string `json:"due_date"`
	ColumnID    uint   `json:"column_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	DueDate     *string `json:"due_date"`
	ColumnID    uint   `json:"column_id"`
	Order       int    `json:"order"`
}

func GetTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		if err := db.Order("column_id, \"order\"").Find(&tasks).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to fetch tasks"})
			return
		}
		c.JSON(200, tasks)
	}
}

func CreateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		// Validate column exists
		var column models.Column
		if err := db.First(&column, req.ColumnID).Error; err != nil {
			c.JSON(400, gin.H{"error": "column not found"})
			return
		}

		task := models.Task{
			Title:       req.Title,
			Description: req.Description,
			Priority:    req.Priority,
			ColumnID:    req.ColumnID,
		}

		if err := db.Create(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to create task"})
			return
		}

		c.JSON(201, task)
	}
}

func UpdateTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req UpdateTaskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "task not found"})
			return
		}

		// Update fields
		if req.Title != "" {
			task.Title = req.Title
		}
		if req.Description != "" {
			task.Description = req.Description
		}
		if req.Priority != "" {
			task.Priority = req.Priority
		}
		if req.ColumnID != 0 {
			// Validate column exists
			var column models.Column
			if err := db.First(&column, req.ColumnID).Error; err != nil {
				c.JSON(400, gin.H{"error": "column not found"})
				return
			}
			task.ColumnID = req.ColumnID
		}
		if req.Order != 0 {
			task.Order = req.Order
		}

		if err := db.Save(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to update task"})
			return
		}

		c.JSON(200, task)
	}
}

func DeleteTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var task models.Task
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "task not found"})
			return
		}

		if err := db.Delete(&task).Error; err != nil {
			c.JSON(500, gin.H{"error": "failed to delete task"})
			return
		}

		c.JSON(200, gin.H{"message": "task deleted"})
	}
}

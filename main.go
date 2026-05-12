package main

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"tablero/internal/database"
	"tablero/internal/handlers"
	"tablero/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin
	router := gin.Default()

	// Apply middleware
	router.Use(cors.Default())
	router.Use(middleware.SecurityHeaders())

	// API routes
	apiGroup := router.Group("/api")

	// Public routes
	apiGroup.POST("/login", handlers.Login)
	apiGroup.POST("/logout", handlers.Logout)

	// Protected routes
	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(middleware.AuthMiddleware())

	// Columns routes
	protectedGroup.GET("/columns", handlers.GetColumns(db))
	protectedGroup.POST("/columns", handlers.CreateColumn(db))
	protectedGroup.PUT("/columns/:id", handlers.UpdateColumn(db))
	protectedGroup.DELETE("/columns/:id", handlers.DeleteColumn(db))

	// Tasks routes
	protectedGroup.GET("/tasks", handlers.GetTasks(db))
	protectedGroup.POST("/tasks", handlers.CreateTask(db))
	protectedGroup.PUT("/tasks/:id", handlers.UpdateTask(db))
	protectedGroup.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// Serve embedded frontend
	fsys, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		// Frontend not embedded (development mode)
		fsys = os.DirFS("frontend/dist")
	}

	// Serve static files and SPA with fallback to index.html
	fileServer := http.FileServer(http.FS(fsys))
	router.NoRoute(func(c *gin.Context) {
		// Try to serve the file
		if f, err := fsys.Open(c.Request.URL.Path[1:]); err == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
		} else {
			// Fallback to index.html for SPA routing
			c.Request.URL.Path = "/index.html"
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
	})

	// Start server
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

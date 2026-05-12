package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"tablero/internal/database"
	"tablero/internal/handlers"
	"tablero/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist
var frontendFS embed.FS

func getContentType(path string) string {
	if strings.HasSuffix(path, ".css") {
		return "text/css; charset=utf-8"
	}
	if strings.HasSuffix(path, ".js") {
		return "application/javascript; charset=utf-8"
	}
	if strings.HasSuffix(path, ".html") {
		return "text/html; charset=utf-8"
	}
	if strings.HasSuffix(path, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(path, ".svg") {
		return "image/svg+xml"
	}
	if strings.HasSuffix(path, ".ico") {
		return "image/x-icon"
	}
	return "text/plain; charset=utf-8"
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(middleware.SecurityHeaders())
	router.RedirectTrailingSlash = false

	// API routes
	apiGroup := router.Group("/api")
	apiGroup.GET("/health", handlers.Health)
	apiGroup.POST("/login", handlers.Login)
	apiGroup.POST("/logout", handlers.Logout)

	protectedGroup := apiGroup.Group("")
	protectedGroup.Use(middleware.AuthMiddleware())
	protectedGroup.GET("/columns", handlers.GetColumns(db))
	protectedGroup.POST("/columns", handlers.CreateColumn(db))
	protectedGroup.PUT("/columns/:id", handlers.UpdateColumn(db))
	protectedGroup.DELETE("/columns/:id", handlers.DeleteColumn(db))
	protectedGroup.GET("/tasks", handlers.GetTasks(db))
	protectedGroup.POST("/tasks", handlers.CreateTask(db))
	protectedGroup.PUT("/tasks/:id", handlers.UpdateTask(db))
	protectedGroup.DELETE("/tasks/:id", handlers.DeleteTask(db))

	// Use DirFS in development (disk), embedded FS in production (Docker)
	var distFS fs.FS
	if _, err := os.Stat("frontend/dist"); err == nil {
		// Development: read from disk
		distFS = os.DirFS("frontend/dist")
	} else {
		// Production: use embedded FS
		distFS, _ = fs.Sub(frontendFS, "frontend/dist")
	}

	// Root
	router.GET("/", func(c *gin.Context) {
		http.FileServer(http.FS(distFS)).ServeHTTP(c.Writer, c.Request)
	})

	// NoRoute: static files + SPA fallback
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api") {
			c.Status(http.StatusNotFound)
			return
		}

		// Try static file
		f, err := distFS.Open(strings.TrimPrefix(path, "/"))
		if err == nil {
			defer f.Close()
			stat, _ := f.Stat()
			if !stat.IsDir() {
				http.ServeContent(c.Writer, c.Request, path, stat.ModTime(), f.(io.ReadSeeker))
				return
			}
		}

		// SPA fallback
		f, err = distFS.Open("index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		defer f.Close()
		stat, _ := f.Stat()
		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), f.(io.ReadSeeker))
	})

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

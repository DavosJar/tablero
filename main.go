package main

import (
	"embed"
	"log"
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

// readDistFile reads a file from frontend/dist/. In development reads from disk,
// in production reads from the embedded filesystem.
func readDistFile(path string) ([]byte, error) {
	if _, err := os.Stat("frontend/dist"); err == nil {
		log.Printf("readDistFile: reading from disk: frontend/dist/%s", path)
		return os.ReadFile("frontend/dist/" + path)
	}
	log.Printf("readDistFile: reading from embed: %s", path)
	data, err := frontendFS.ReadFile(path)
	if err != nil {
		log.Printf("readDistFile: ERROR reading from embed: %v", err)
	}
	return data, err
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

	// Debug: list embedded FS contents
	entries, err := frontendFS.ReadDir(".")
	if err != nil {
		log.Printf("DEBUG embed: error reading root: %v", err)
	} else {
		log.Printf("DEBUG embed: root entries:")
		for _, e := range entries {
			log.Printf("  - %s (dir=%v)", e.Name(), e.IsDir())
			if e.IsDir() {
				subs, _ := frontendFS.ReadDir(e.Name())
				for _, s := range subs {
					log.Printf("    - %s/%s", e.Name(), s.Name())
				}
			}
		}
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

	// Root - serve index.html
	router.GET("/", func(c *gin.Context) {
		data, err := readDistFile("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// NoRoute: static files + SPA fallback
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api") {
			c.Status(http.StatusNotFound)
			return
		}

		// Try static file
		filePath := strings.TrimPrefix(path, "/")
		if filePath != "" {
			data, err := readDistFile(filePath)
			if err == nil {
				c.Data(http.StatusOK, getContentType(path), data)
				return
			}
			log.Printf("NoRoute: readDistFile(%q) failed: %v", filePath, err)
		}

		// SPA fallback
		data, err := readDistFile("index.html")
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

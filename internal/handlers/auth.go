package handlers

import (
	"os"
	"sync"
	"time"

	"tablero/internal/utils"

	"github.com/gin-gonic/gin"
)

var (
	loginAttempts = make(map[string][]time.Time)
	attemptsMutex sync.Mutex
)

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	clientIP := c.ClientIP()

	// Rate limiting: máx 5 intentos por minuto
	attemptsMutex.Lock()
	now := time.Now()
	oneMinuteAgo := now.Add(-time.Minute)

	// Limpiar intentos antiguos
	if attempts, exists := loginAttempts[clientIP]; exists {
		validAttempts := []time.Time{}
		for _, attempt := range attempts {
			if attempt.After(oneMinuteAgo) {
				validAttempts = append(validAttempts, attempt)
			}
		}
		loginAttempts[clientIP] = validAttempts

		if len(validAttempts) >= 5 {
			attemptsMutex.Unlock()
			c.JSON(429, gin.H{"error": "too many login attempts"})
			return
		}
	}

	loginAttempts[clientIP] = append(loginAttempts[clientIP], now)
	attemptsMutex.Unlock()

	appPassword := os.Getenv("APP_PASSWORD")
	if appPassword == "" {
		appPassword = "default-password"
	}

	if req.Password != appPassword {
		c.JSON(401, gin.H{"error": "invalid password"})
		return
	}

	token, err := utils.GenerateJWT()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(200, LoginResponse{
		Token:     token,
		ExpiresIn: 86400, // 24 horas
	})
}

func Logout(c *gin.Context) {
	// JWT is stateless, so logout just signals the client to clear the token
	c.JSON(200, gin.H{"message": "logged out"})
}

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fine-track/auth-app/src/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func corsMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"userAgent": c.Request.UserAgent(),
		})
	})

	r.POST("/authorize", verifyAccessTokenMiddleware, handleAuthorize)

	r.GET("/profile", verifyAccessTokenMiddleware, handleGetProfile)

	r.POST("/get-access-token", handleGetAccessToken)

	r.POST("/login", handleLogin)

	r.POST("/logout", verifyAccessTokenMiddleware, handleLogout)

	r.POST("/register", handleRegister)

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	client := db.ConnectToDb()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	r := setupRouter()
	err = r.Run(":8081")
	if err != nil {
		panic(err)
	}
}

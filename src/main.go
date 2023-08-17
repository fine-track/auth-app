package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fine-track/auth-app/src/db"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(corsMiddleware)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":   "pong",
			"userAgent": c.Request.UserAgent(),
		})
	})

	authGrp := r.Group("/auth")
	authGrp.GET("/profile", verifyAccessTokenMiddleware, handleGetProfile)
	authGrp.POST("/authorize", verifyAccessTokenMiddleware, handleAuthorize)
	authGrp.POST("/logout", verifyAccessTokenMiddleware, handleLogout)
	authGrp.POST("/get-access-token", handleGetAccessToken)
	authGrp.POST("/login", handleLogin)
	authGrp.POST("/register", handleRegister)

	r.Use(handleNotFound)

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

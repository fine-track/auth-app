package main

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sifatulrabbi/fine-track/auth-api/src/utils"
)

func verifyAccessTokenMiddleware(c *gin.Context) {
	res := utils.HTTPResponse{}
	authToken := c.Request.Header.Get("authorization")
	if authToken == "" {
		res.Message = "No authorization token found!"
		res.Unauthorized(c)
		c.Abort()
		return
	}

	splitted := strings.Split(authToken, " ")
	if strings.ToLower(splitted[0]) != "bearer" || splitted[1] == "" {
		res.Message = "Invalid authorization token format"
		res.Unauthorized(c)
		c.Abort()
		return
	}

	accessToken := splitted[1]
	claims := AccessTokenClaims{}
	if err := ValidateAccessToken(accessToken, &claims); err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
		c.Abort()
	} else {
		c.Set("sessionData", claims)
		c.Next()
	}
}

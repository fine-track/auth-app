package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fine-track/auth-app/src/db"
	"github.com/fine-track/auth-app/src/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserPayload struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutPayload struct {
	RefreshToken string `json:"refresh_token"`
}

func handleLogin(c *gin.Context) {
	// get the user email and password or the OTP
	res := utils.HTTPResponse{}
	mode := c.DefaultQuery("mode", "password")
	var body LoginPayload
	if err := c.BindJSON(&body); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}
	if body.Email == "" || body.Password == "" {
		res.Message = "Invalid request payload. `Email` and `Password` is required."
		res.Unauthorized(c)
		return
	}

	user := db.User{}
	if err := user.GetUserByEmail(body.Email); err != nil {
		log.Println(err.Error())
		res.Message = "No user account found with the email"
		res.Unauthorized(c)
		return
	}

	if mode != "password" && mode != "otp" {
		res.Message = "Invalid login mode"
		res.Unauthorized(c)
		return
	}
	// if the mode is password verify the user email and the password
	if mode == "password" {
		if valid := utils.IsValidPass(body.Password, user.Password); !valid {
			res.Message = "Invalid password"
			res.Unauthorized(c)
			return
		}
	}
	// if the mode is otp then verify the user email and the OTP
	if mode == "otp" {
		otp := db.OTPSession{}
		if err := otp.GetByEmail(user.Email); err != nil {
			res.Message = "Unable to verify the OTP please try again."
			res.Unauthorized(c)
			return
		}
		if body.Password != otp.Code {
			res.Message = "Invalid OTP please enter the right code."
			res.Unauthorized(c)
			return
		}
	}

	// if verified create a new session
	session := db.Session{}
	session.Email = user.Email
	session.UserId = user.ID
	session.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	session.UserAgent = c.Request.UserAgent()
	if err := session.CreateNew(); err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
		return
	}

	// after creating a new session for the user generate a access token for them
	if accessToken, err := GetAccessTokenForSession(session.Email, session.UserId.Hex()); err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
	} else {
		res.Message = "Login successful"
		res.Data = map[string]string{
			"refresh_token": session.ID.Hex(),
			"access_token":  accessToken,
		}
		res.Ok(c)
	}
}

func handleRegister(c *gin.Context) {
	res := utils.HTTPResponse{}
	var body RegisterUserPayload
	if err := c.BindJSON(&body); err != nil {
		fmt.Println("while getting body -> ", err.Error())
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}
	// validate the request payload
	if body.Fullname == "" {
		res.Message = "`Fullname` is required"
		res.BadRequest(c)
		return
	}
	if body.Email == "" {
		res.Message = "`Email` is required"
		res.BadRequest(c)
		return
	}
	if body.ConfirmPassword == "" {
		res.Message = "Please enter confirm password"
		res.BadRequest(c)
		return
	}
	if body.Password == "" {
		res.Message = "`Password` is required"
		res.BadRequest(c)
		return
	}
	if body.Password != body.ConfirmPassword {
		res.Message = "Passwords don't match"
		res.BadRequest(c)
		return
	}

	// before saving the user info make sure to hash the user password with bcrypt
	if hashedPass, err := utils.HashPass(body.Password); err != nil {
		res.Message = err.Error()
		res.InternalServerError(c)
	} else {
		body.Password = hashedPass
	}

	// create a new user in the db.
	user := db.User{
		Email:    body.Email,
		Fullname: body.Fullname,
		Password: body.Password,
	}
	if err := user.CreateNew(); err != nil {
		fmt.Println("while creating user -> ", err.Error())
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}

	res.Message = "New user registered"
	res.Data = map[string]any{
		"_id":      user.ID,
		"email":    user.Email,
		"fullname": user.Fullname,
	}
	res.Ok(c)
}

func handleLogout(c *gin.Context) {
	res := utils.HTTPResponse{}
	if temp, exists := c.Get("sessionData"); !exists {
		res.Message = "Unable to identify the user"
		res.Unauthorized(c)
		return
	} else {
		if sessionData, ok := temp.(AccessTokenClaims); !ok {
			res.Message = "Invalid session"
			res.Unauthorized(c)
			return
		} else {
			if err := db.RemoveUserSessions(sessionData.Email); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	res.Message = "Successfully logged out"
	res.Ok(c)
}

// verifies the user credentials and returns a refresh token
func handleAuthorize(c *gin.Context) {
	res := utils.HTTPResponse{}
	if temp, exists := c.Get("sessionData"); !exists {
		res.Message = "No session found"
		res.Data = map[string]string{"code": "SESSION_NOT_FOUND"}
		res.Unauthorized(c)
	} else {
		if sessionData, ok := temp.(AccessTokenClaims); !ok {
			res.Message = "Invalid session info"
			res.Data = map[string]string{"code": "SESSION_NOT_FOUND"}
			res.Unauthorized(c)
		} else {
			res.Message = "User authorized"
			res.Data = sessionData
			res.Ok(c)
		}
	}
}

// verifies a refresh token or user session and sends back an access token which is valid for 2 hours
func handleGetAccessToken(c *gin.Context) {
	res := utils.HTTPResponse{}
	authToken := c.Request.Header.Get("authorization")
	if authToken == "" {
		res.Message = "No authorization token found"
		res.Unauthorized(c)
		return
	}

	splitted := strings.Split(authToken, " ")
	if strings.ToLower(splitted[0]) != "bearer" || splitted[1] == "" {
		res.Message = "Invalid authorization token"
		res.Unauthorized(c)
		return
	}

	session := db.Session{}
	if err := session.GetById(splitted[1]); err != nil {
		fmt.Println(splitted[1])
		res.Message = "Unable to verify the authorization token"
		res.Unauthorized(c)
		return
	}
	if accessToken, err := GetAccessTokenForSession(session.Email, session.UserId.Hex()); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
	} else {
		res.Data = map[string]string{"access_token": accessToken}
		res.Message = "New access token which will be valid for the next 1 hour"
		res.Ok(c)
	}
}

func handleGetProfile(c *gin.Context) {
	res := utils.HTTPResponse{}
	temp, exists := c.Get("sessionData")
	if !exists {
		res.Message = "Invalid auth token"
		res.Forbidden(c)
		return
	}
	sessionData, ok := temp.(AccessTokenClaims)
	if !ok {
		res.Message = "Invalid request"
		res.Forbidden(c)
		return
	}
	user := db.User{}
	if err := user.GetById(sessionData.UserId); err != nil {
		res.Message = err.Error()
		res.BadRequest(c)
		return
	}
	// only sending relevant data and omitting sensitive data such as passwords.
	res.Data = gin.H{
		"_id":        user.ID,
		"fullname":   user.Fullname,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	res.Ok(c)
}

func handleNotFound(c *gin.Context) {
	res := utils.HTTPResponse{}
	res.Message = fmt.Sprintf("%s: %s is not found!", c.Request.Method, c.Request.URL.Path)
	res.NotFound(c)
}

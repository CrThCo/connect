package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/MartinResearchSociety/connect-api/controllers"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var hmackey []byte

func init() {
	hmackey = []byte(os.Getenv("H_MAC_KEY"))
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token := parts[1]
		parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return hmackey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			c.Set("UID", claims["sub"])
			c.Next()
		}
	}
}

func routes(r *gin.Engine) {
	post := new(controllers.PostController)
	user := new(controllers.UserController)
	auth := new(controllers.AuthController)

	v1 := r.Group("/v1")
	{
		p := v1.Group("/post")
		p.Use(authMiddleware())
		{
			p.POST("/save", post.New)
			p.POST("/list", post.List)
			p.POST("/vote/:post_id", post.Vote)
		}

		vote := v1.Group("/vote")
		vote.Use(authMiddleware())
		{
			vote.GET("/count/user/:user_id", post.CountVoteByUser)
			vote.GET("/list/user/:user_id", post.GetVoteByUser)
			vote.GET("/count/post/:post_id", post.CountVoteByPost)
			vote.GET("/list/post/:post_id", post.GetVoteByPost)
		}

		v1.POST("/user/signup", user.Signup)
		v1.POST("/user/signin", auth.Signin)

		u := v1.Group("/user")
		u.Use(authMiddleware())
		{
			u.GET("/list", user.List)
			u.GET("/i/:uid", user.Info)
			u.POST("/update", user.Update)
			u.POST("/upload/image", user.Upload)
		}
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Static("/static", os.Getenv("FILE_STORAGE_PATH"))
	r.StaticFile("/favicon.ico", "favicon.ico")
}

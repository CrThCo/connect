package controllers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MartinResearchSociety/connect-api/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// AuthController type
type AuthController struct {
	*Controller
}

var (
	issuer  string
	hmackey []byte
)

func init() {
	issuer = os.Getenv("JWT_TOKEN_ISSUER")
	hmackey = []byte(os.Getenv("H_MAC_KEY"))
}

func (a *AuthController) getJWT(uid string) (string, error) {
	h := md5.New()
	ct := time.Now()
	expireAt := ct.Add(24 * time.Hour)
	sc := jwt.StandardClaims{
		Subject:   uid,
		IssuedAt:  ct.Unix(),
		NotBefore: ct.Unix(),
		ExpiresAt: expireAt.Unix(),
		Issuer:    issuer,
		Id:        fmt.Sprintf("%x", h.Sum(nil)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, sc)
	return token.SignedString(hmackey)
}

// Signin method
func (a *AuthController) Signin(c *gin.Context) {
	user := db.User{}
	if err := c.Bind(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	uid, err := db.GetUserByCredentials(user.Email, user.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	jwtToken, err := a.getJWT(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}

// Refresh method
func (a *AuthController) Refresh(c *gin.Context) {
	uid := c.GetString("UID")
	user := db.User{}
	err := db.GetUser(uid, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	jwtToken, err := a.getJWT(user.ID.Hex())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSONP(http.StatusOK, gin.H{
		"token": jwtToken,
	})
}

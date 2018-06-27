package controllers

import (
	"crypto/md5"
	"time"
	"encoding/json"

	"github.com/MartinResearchSociety/connect/models"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

type LoginController struct {
	beego.Controller
}

// @Title Post 
// @Description create tokens
// @Param	body		body 	models.User	true "body for user content"
// @Success 200 {int} models.JWT
// @Failure 403 body is empty
// @router / [post]
func (l *LoginController) Post() {
	var u models.User
	var err error
	json.Unmarshal(l.Ctx.Input.RequestBody, &u)
	uid, err := models.GetUserByCredentials(u.Username, u.Password)
	currentTimestamp := time.Now().UTC().Unix()
	// md5 of sub & iat
	h := md5.New()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject: string(uid),
		IssuedAt: currentTimestamp,
		NotBefore: currentTimestamp,
		// TODO: make it configurable for production
		ExpiresAt: currentTimestamp + 3600,
		//TODO: change in production - should be configurable
		Issuer: l.Ctx.Input.Domain(),
		Id: string(h.Sum(nil)),
	})

	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("HMACKEY")))

	if err != nil {
    	l.Data["json"] = err.Error()
	}

	l.Data["json"] = map[string]string{"token": tokenString}
	l.ServeJSON()
}

package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MartinResearchSociety/connect/models"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

type LoginController struct {
	beego.Controller
}

func (l *LoginController) getJWTToken(uid string) (string, error) {
	h := md5.New()
	currentTimestamp := time.Now().UTC().Unix()
	standardClaims := jwt.StandardClaims{
		Subject:   uid,
		IssuedAt:  currentTimestamp,
		NotBefore: currentTimestamp,
		// TODO: make it configurable for production
		ExpiresAt: currentTimestamp + 13600,
		//TODO: change in production - should be configurable
		Issuer: l.Ctx.Input.Domain(),
		Id:     fmt.Sprintf("%x", h.Sum(nil)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, standardClaims)
	return token.SignedString([]byte(beego.AppConfig.String("HMACKEY")))
}

//@Title Refresh
//@Description refresh authentication token
//@Success 200 {string} models.JWT
//@Failure 403 body is empty
//@router /token/refresh
func (l *LoginController) Refresh() {
	userID := l.Ctx.Input.Param("userID")
	u, err := models.GetUser(userID)
	if err != nil {
		l.Ctx.Output.SetStatus(401)
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}
	tokenString, err := l.getJWTToken(u.ID.Hex())
	if err != nil {
		l.Ctx.Output.SetStatus(401)
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}
	l.Ctx.Output.SetStatus(200)
	l.Data["json"] = map[string]string{"token": tokenString}
	l.ServeJSON()
}

// @Title Auth
// @Description create tokens
// @Param	body		body 	models.Auth	true "body for user content"
// @Success 200 {int} models.JWT
// @Failure 403 body is empty
// @router /login [post]
func (l *LoginController) Auth() {
	var u models.Auth
	if err := json.Unmarshal(l.Ctx.Input.RequestBody, &u); err != nil {
		l.Ctx.Output.SetStatus(400)
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}

	uid, err := models.GetUserByCredentials(u.Email, u.Password)

	if err != nil {
		l.Ctx.Output.SetStatus(401)
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}

	if uid == "" {
		l.Ctx.Output.SetStatus(401)
		l.Data["json"] = ""
		l.ServeJSON()
		return
	}

	tokenString, err := l.getJWTToken(uid)
	if err != nil {
		l.Ctx.Output.SetStatus(401)
		l.Data["json"] = err.Error()
		l.ServeJSON()
		return
	}

	l.Ctx.Output.SetStatus(200)
	l.Data["json"] = map[string]string{"token": tokenString}
	l.ServeJSON()
}

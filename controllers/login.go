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
	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("HMACKEY")))
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

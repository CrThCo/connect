package main

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/MartinResearchSociety/connect/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	var AuthFilter = func(ctx *context.Context) {

		ctx.Output.Header("Content-Type", "application/json")
		header := strings.Split(ctx.Input.Header("Authorization"), " ")
		if len(header) != 2 || header[0] != "Bearer" {
			ctx.Abort(401, "Not Authorized")
		}

		var tokenString string = ctx.Input.Header("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// TODO: change for production
			return []byte(beego.AppConfig.String("HMACKEY")), nil
		})

		if err != nil {
			ctx.Output.SetStatus(403)
			resBytes, err := json.Marshal(err.Error())
			ctx.Output.Body(resBytes)
			if err != nil {
				panic(err)
			}
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims != nil {
			return
		}
		ctx.Output.SetStatus(403)
		ctx.Output.Body([]byte("Invalid token!"))
		if err != nil {
			panic(err)
		}

	}

	//TODO: everything is filtered?!
	beego.InsertFilter("/asdasd*", beego.BeforeRouter, AuthFilter)

	beego.Run()
}

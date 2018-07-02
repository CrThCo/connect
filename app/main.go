package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	_ "github.com/MartinResearchSociety/connect/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"

	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/dghubble/sessions"
)

const (
	sessionName    = "example-twtter-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "twitterID"
)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	// 1. Register Twitter login and callback handlers
	oauth1Config := &oauth1.Config{
		ConsumerKey:    beego.AppConfig.String("TwitterKey"),
		ConsumerSecret: beego.AppConfig.String("TwitterSecret"),
		CallbackURL:    beego.AppConfig.String("TwitterCallback"),
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
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
		ctx.Output.SetStatus(401)
		ctx.Output.Body([]byte("Invalid token!"))
		if err != nil {
			panic(err)
		}
	}

	//TODO: make it so that all filtered routes lie under this
	beego.InsertFilter("/v1/user/*", beego.BeforeRouter, AuthFilter)

	beego.Handler("/twitter/login", twitter.LoginHandler(oauth1Config, nil))
	beego.Handler("/twitter/callback", twitter.CallbackHandler(oauth1Config, issueSession(), nil))

	beego.Run()
}

// issueSession issues a cookie session after successful Twitter login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Print(twitterUser)
		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = twitterUser.ID
		session.Save(w)
		//TODO: redirect when it works
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

package main

import (
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/MartinResearchSociety/connect/routers"
	"github.com/MartinResearchSociety/connect/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/dgrijalva/jwt-go"

	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"

	"github.com/astaxie/beego/session"
)

// var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

var globalSessions *session.Manager

func init() {
    globalSessions, _ = session.NewManager("memory", &session.ManagerConfig{CookieName: "userId", EnableSetCookie: true, Gclifetime:3600, Maxlifetime: 3600, Secure: false, CookieLifeTime: 3600})
    go globalSessions.GC()
}

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

		// Excluding Loging and Register route
		switch ctx.Request.RequestURI {
		case "/v1/user/signup", "/v1/user/login":
			return
		}

		ctx.Output.Header("Content-Type", "application/json")
		ctx.Output.Header("Access-Control-Allow-Origin", "*")
		header := strings.Split(ctx.Input.Header("Authorization"), " ")
		if len(header) != 2 || header[0] != "Bearer" {
			ctx.Abort(401, "Not Authorized")
		}

		var tokenString string = header[1]

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

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sub, _ := claims["sub"].(string)
			ctx.Input.SetParam("userID", sub)
			return
		}

		ctx.Output.SetStatus(401)
		ctx.Output.Body([]byte("Invalid token!"))
		if err != nil {
			panic(err)
		}
	}

	beego.SetStaticPath("/file", beego.AppConfig.String("FileStoragePath"))
	beego.Handler("/twitter/login", twitter.LoginHandler(oauth1Config, nil))
	beego.Handler("/twitter/callback", twitter.CallbackHandler(oauth1Config, issueSession(), nil))
	//TODO: everything is filtered?!
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"content-type", "authorization"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//TODO: make it so that all filtered routes lie under this
	beego.InsertFilter("/v1/*", beego.BeforeRouter, AuthFilter)

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
		
		// session := sessionStore.New(sessionName)
		// session.Values[sessionUserKey] = twitterUser.ID
		// // session.Values[sessionUserName] = twitterUser.ScreenName
		// session.Save(w)
		sess, _ := globalSessions.SessionStart(w, req)
		defer sess.SessionRelease(w)

		if u, err := models.GetUserByEmail(twitterUser.Email); err != nil {
			fmt.Printf("Error trying to get user by email: %v", err.Error())
		} else if u == nil {
			fmt.Printf("User not added yet: email=%v", twitterUser.Email)
		} else if u != nil {
			// session.Values[sessionUserID] = u.ID.String()
			sess.Set("userId", u.ID.String())
		} else {
			// add user, if not already present
		u := &models.User{
			ID: bson.NewObjectId(), 
			Username: twitterUser.Name, 
			Email: twitterUser.Email,
		}

		if err := u.Insert(); err != nil {
			fmt.Printf("Couldn't add user: email=%v, err=%v", twitterUser.Email, err)
		}

		// session.Values[sessionUserID] = u.ID.String()
		sess.Set("userId", u.ID.String())
		}

		//TODO: redirect when it works
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

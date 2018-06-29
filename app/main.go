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
	. "github.com/kkdai/twitter"
)

var ConsumerKey string
var ConsumerSecret string
var twitterClient *ServerClient

func init() {
	ConsumerKey=    "9b6zZShiwX6VKKHOTZqNq5Phz"
	ConsumerSecret= "5PJUoBMA1D3AhIXQW1KF8VRMH2EDaD2iS2TAaPLpkOC6bmFWHD"
}

const (
	sessionName    = "example-twtter-app"
	sessionSecret  = "example cookie signing secret"
	sessionUserKey = "twitterID"
		//This URL need note as follow:
		// 1. Could not be localhost, change your hosts to a specific domain name
		// 2. This setting must be identical with your app setting on twitter Dev
		CallbackURL string = "http://myserver.local:8030/maketoken"
	)

var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}	
	
	twitterClient = NewServerClient(ConsumerKey, ConsumerSecret)


	// 1. Register Twitter login and callback handlers
	oauth1Config := &oauth1.Config{
		ConsumerKey:    "9b6zZShiwX6VKKHOTZqNq5Phz",
		ConsumerSecret: "5PJUoBMA1D3AhIXQW1KF8VRMH2EDaD2iS2TAaPLpkOC6bmFWHD",
		CallbackURL:    "http://myserver.local:8050/twitter/callback",
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

	beego.Handler("/profile", profileHandler())
	beego.Handler("/twitter/login", blah())
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

func blah() http.Handler {
	return  http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Enter redirect to twitter")
			fmt.Println("Token URL=", CallbackURL)
			requestUrl := twitterClient.GetAuthURL(CallbackURL)
			http.Redirect(w, r, requestUrl, http.StatusTemporaryRedirect)
			fmt.Println("Leave redirect") })
}
func profileHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
})
}

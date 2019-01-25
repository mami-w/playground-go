package main

import (
	_ "crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/session"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"golang.org/x/oauth2"
)

type authConfig struct{
	ClientID string
	ClientSecret string
	Domain string
	CallbackUrl string
}

func newServer() (router *mux.Router) {

	// this shoudl come from config or env variables


	cf := &session.ManagerConfig{CookieName:"gosessionid", Gclifetime:3600}
	sessionManager, _ := session.NewManager("memory", cf)
	go sessionManager.GC()

	config := &authConfig{
		ClientID:"kWwxrezJsEmJn3Iysy5qHFMvM6rLBDJ2",
		ClientSecret:"PkLVlOh4hGse7gQ7TfvI4KydP1LU0EB6j9jLBcY-gC8q67n2mtN4UyOanYRemGc_",
		Domain:"mami-w.auth0.com",
		CallbackUrl:"http://localhost:8002/callback",
	}
	router = mux.NewRouter()

	initRoutes(router, sessionManager, config)

	return router
}

func initRoutes(router *mux.Router, sessionManager *session.Manager, config *authConfig) {

	router.HandleFunc("/", homeHandler(config))
	router.HandleFunc("/callback", callbackHandler(sessionManager, config))
	router.HandleFunc("/user", userHandler(sessionManager))
	router.Use(newAuthMiddleware(sessionManager))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
}

// technically middleware.go
func newAuthMiddleware (sessionManager *session.Manager) mux.MiddlewareFunc {
	return mux.MiddlewareFunc(func(next http.Handler) (http.Handler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// if we are not prefixed by '/user', call next.ServeHTTP(w, r) immediately
			requestUri := r.RequestURI
			if strings.HasPrefix(strings.ToLower(requestUri), "/user") {
				session, _ := sessionManager.SessionStart(w, r)
				defer session.SessionRelease(w)

				if session.Get("profile") == nil {
					http.Redirect(w, r, "/", http.StatusMovedPermanently)
				} else {
					next.ServeHTTP(w, r)
				}
			} else {
				next.ServeHTTP(w,r)
			}
		})
	})
}

// technically home_handler.go
var bodyTemplate = `
	<script src="https://cdn.auth0.com/js/lock/11.13.1/lock.min.js"></script>
	<script type="text/javascript">
		var lock = new Auth0Lock('{{.ClientID}}', '{{.Domain}}');
		function signin() {
			lock.show({
				callbackURL: '{{.CallbackUrl}}'
				, responseType: 'code'
				, authParams: {
					scope: 'openid profile'
				}
			})
		}
	</script>
	<button onclick="window.signin();">Login</button>
`

func homeHandler(config *authConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		t := template.Must(template.New("htmlz").Parse(bodyTemplate))
		t.Execute(w, config)
	}
}

// technically callback_handler.go
func callbackHandler(sessionManager *session.Manager, config *authConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Instantiating the OAuth2 package to exchange the Code for a Token
		conf := &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.CallbackUrl,
			Scopes:       []string{"openid", "name", "email", "picture"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + config.Domain + "/authorize",
				TokenURL: "https://" + config.Domain + "/oauth/token",
			},
		}

		// Getting the Code that we got from Auth0
		code := r.URL.Query().Get("code")

		// Exchanging the code for a token
		token, err := conf.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Getting now the User information
		client := conf.Client(oauth2.NoContext, token)
		resp, err := client.Get("https://" + config.Domain + "/userinfo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Reading the body
		raw, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Unmarshalling the JSON of the Profile
		var profile map[string]interface{}
		if err := json.Unmarshal(raw, &profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, _ := sessionManager.SessionStart(w, r)
		defer session.SessionRelease(w)

		session.Set("id_token", token.Extra("id_token"))
		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)

		// Redirect to logged in page
		http.Redirect(w, r, "/user", http.StatusMovedPermanently)

	}
}
// technically userHandler.go

func userHandler(sessionManager *session.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := sessionManager.SessionStart(w, req)
		defer session.SessionRelease(w)

		profile := session.Get("profile")
		fmt.Fprintf(w, "USER DATA: %+v", profile)
	}
}
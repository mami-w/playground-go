package server

import (
	_ "crypto/sha512"
	"errors"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

const (
	//defaultAddress = "http://ec2-18-217-168-85.us-east-2.compute.amazonaws.com"
	defaultAddress = "http://localhost:8000"
)

func NewServer() (router *mux.Router){
	//Init()

	router = mux.NewRouter()
	initRoutes(router)
	return router
}

func initRoutes(router *mux.Router) {

	trackerEndpoint, apiSecret, err := getCmdlineFlags()
	secret := []byte(apiSecret)
	// todo: better error handling
	if err != nil { panic(err) }

	router.Handle("/api/v1.0/tracker/user", wrap(getAllUsersHandler(trackerEndpoint), authMiddleware(secret))).Methods("GET")

	router.Handle("/api/v1.0/tracker/user/{userid}", wrap(getAllEntriesHandler(trackerEndpoint), authMiddleware(secret))).Methods("GET")
	router.Handle("/api/v1.0/tracker/user/{userid}", wrap(createUserHandler(trackerEndpoint), authMiddleware(secret))).Methods("POST")
	router.Handle("/api/v1.0/tracker/user/{userid}", wrap(updateUserHandler(trackerEndpoint), authMiddleware(secret))).Methods("PUT")
	router.Handle("/api/v1.0/tracker/user/{userid}", wrap(deleteUserHandler(trackerEndpoint), authMiddleware(secret))).Methods("DELETE")

	router.Handle("/api/v1.0/tracker/user/{userid}/entry/{entryid}", wrap(getEntryHandler(trackerEndpoint), authMiddleware(secret))).Methods("GET")
	router.Handle("/api/v1.0/tracker/user/{userid}/entry/{entryid}", wrap(createEntryHandler(trackerEndpoint), authMiddleware(secret))).Methods("POST")
	router.Handle("/api/v1.0/tracker/user/{userid}/entry/{entryid}", wrap(updateEntryHandler(trackerEndpoint), authMiddleware(secret))).Methods("PUT")
	router.Handle("/api/v1.0/tracker/user/{userid}/entry/{entryid}", wrap(deleteEntryHandler(trackerEndpoint), authMiddleware(secret))).Methods("DELETE")

	router.PathPrefix("/").Handler(	http.FileServer(http.Dir("./assets/")))
}

func getCmdlineFlags() (endpoint string, apiSecret string, err error) {

	flag.StringVar(&endpoint, "endpoint", "", "endpoint for tracker data")
	flag.StringVar(&apiSecret, "apisecret", "", "secret for api access token")

	flag.Parse()

	if flag.NFlag() == 0 {
		log.Println("using default entpoint, no secret")
		return defaultAddress, "", nil
	}
	if flag.NFlag() < 2 {
		return "", "", errors.New("too few arguments")
	}

	log.Printf("endpoint - %v", endpoint)
	return endpoint, apiSecret, nil
}

// todo: route bearer header to backend for auth
func authMiddleware (secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) (http.Handler) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var err error
			var tokenString string

			defer func() {
				if err != nil {
					fmt.Println(err)
					fmt.Printf("Token not valid: %v", tokenString)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
				}
			}()

			tokenString, err = getTokenFromRequest(r)
			if err != nil {
				return
			}

			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
			// head of the token to identify which key to use, but the parsed token (head and claims) is provided
			// to the callback, providing flexibility.
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			})
			if err != nil {
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				fmt.Printf("+%v", claims)
				next.ServeHTTP(w, r)
			} else {
				err = errors.New(fmt.Sprintf("Claims are not valid: %v",claims))
			}
		})
	}
}

func wrap(next http.Handler, middleware func(http.Handler) http.Handler) http.Handler {
	return  middleware(next)
}

func getTokenFromRequest(req *http.Request) (string, error) {

	bearerSchema := "Bearer "

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header required")
	}

	if !strings.HasPrefix(authHeader, bearerSchema) {
		return "", errors.New("Authorization requires Basic/Bearer scheme")
	}

	return authHeader[len(bearerSchema):], nil
}

// todo :)
func handleError(err error) {
	log.Printf(err.Error())
}

func reportRemoteError(remoteErrorMsg string) string {
	return fmt.Sprintf("remote server returned error: %v", remoteErrorMsg)
}

/*
func loginHandler(conf *oauth2.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		aud := "https://" + domain + "/userinfo"

		// Generate random state
		b := make([]byte, 32)
		rand.Read(b)
		state := base64.StdEncoding.EncodeToString(b)

		session, err := Store.Get(r, "state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Values["state"] = state
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		audience := oauth2.SetAuthURLParam("audience", aud)
		url := conf.AuthCodeURL(state, audience)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func callbackHandler(conf *oauth2.Config) func (w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {

		state := r.URL.Query().Get("state")
		session, err := Store.Get(r, "state")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if state != session.Values["state"] {
			http.Error(w, "Invalid state parameter", http.StatusInternalServerError)
			return
		}

		code := r.URL.Query().Get("code")

		token, err := conf.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Getting now the userInfo
		client := conf.Client(context.TODO(), token)
		resp, err := client.Get("https://" + domain + "/userinfo")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		var profile map[string]interface{}
		if err = json.NewDecoder(resp.Body).Decode(&profile); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err = Store.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["id_token"] = token.Extra("id_token")
		session.Values["access_token"] = token.AccessToken
		session.Values["profile"] = profile
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to logged in page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
*/

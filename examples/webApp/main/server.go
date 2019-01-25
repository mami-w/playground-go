package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func newServer() (router *mux.Router) {

	router = mux.NewRouter()

	initRoutes(router)
	return router
	}


func initRoutes(router *mux.Router) {

	router.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images/"))))
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./assets/css/"))))
	router.HandleFunc("/", homeHandler())
}

// initial example
/*
func initRoutes(router *mux.Router) {
	webroot := os.Getenv("WEBROOT")
	if len(webroot) == 0 {
		root, err := os.Getwd()
		if err != nil {
			panic ("Cannot find working directory")
		} else {
			webroot = root
		}
	}

	router.HandleFunc("/api/test", testHandler()).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webroot + "/assets/")))
}
*/
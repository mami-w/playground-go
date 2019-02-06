package server

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/mami-w/playground-go/timetracker/logger"
	"log"
	"net/http"
)

const (
	defaultAddress = "http://ec2-18-217-168-85.us-east-2.compute.amazonaws.com"
)

func NewServer() (router *mux.Router){

	router = mux.NewRouter()
	initRoutes(router)
	return router
}

func initRoutes(router *mux.Router) {

	trackerEndpoint, err := getTrackerEndpoint()
	// todo: better error handling
	if err != nil { panic(err) }

	router.HandleFunc("/api/v1.0/tracker/user", getAllUsersHandler(trackerEndpoint)).Methods("GET")

	router.HandleFunc("/api/v1.0/tracker/user/{userid}", getAllEntriesHandler(trackerEndpoint)).Methods("GET")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}", createUserHandler(trackerEndpoint)).Methods("POST")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}", updateUserHandler(trackerEndpoint)).Methods("PUT")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}", deleteUserHandler(trackerEndpoint)).Methods("DELETE")

	router.HandleFunc("/api/v1.0/tracker/user/{userid}/entry/{entryid}", getEntryHandler(trackerEndpoint)).Methods("GET")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}/entry/{entryid}", createEntryHandler(trackerEndpoint)).Methods("POST")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}/entry/{entryid}", updateEntryHandler(trackerEndpoint)).Methods("PUT")
	router.HandleFunc("/api/v1.0/tracker/user/{userid}/entry/{entryid}", deleteEntryHandler(trackerEndpoint)).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
}


func getTrackerEndpoint() (trackerEndpoint string, err error) {

	endpointFlag := flag.String("endpoint", "", "endpoint for ")

	flag.Parse()

	if flag.NFlag() == 0 {
		return defaultAddress, nil
	}

	trackerEndpoint = *endpointFlag

	logger.Get().Printf("endpoint - %v", trackerEndpoint)
	return trackerEndpoint, nil
}

func handleError(err error) {
	log.Printf(err.Error())
}


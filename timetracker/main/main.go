package main

import (
	"errors"
	"flag"
	"github.com/mami-w/playground-go/timetracker/logger"
	"github.com/mami-w/playground-go/timetracker/other"
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"github.com/mami-w/playground-go/timetracker/trackerdata/memoryStorage"
	"github.com/mami-w/playground-go/timetracker/trackerdata/postgresStorage"
	"log"
	"net/http"
	"os"
)

// https
// auth0

func main() {

	var storage trackerdata.Storage
	var err error

	count, endpoint, user, pwd, err := parseCmdlineParameters()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	switch {
	case count == 0 :
		storage, err = memoryStorage.NewStorage()
	default:
		storage, err = postgresStorage.NewPostgresStorage(endpoint, user, pwd)
	}

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1.0/timetracker/user/", other.UserHandler(storage))
	log.Println("starting to listen on port 8000")
	http.ListenAndServe(":8000", nil)
}

func parseCmdlineParameters() (count int, endpoint string, username string, pwd string, err error) {

	endpointFlag := flag.String("endpoint", "", "endpoint for postgres storage")
	usernameFlag := flag.String("username", "", "user of postgres storage")
	pwdFlag := flag.String("pwd", "", "pwd")

	flag.Parse()

	count = flag.NFlag()

    if count > 0 && count < 3 {
	  return count, endpoint, username, pwd, errors.New("Got more than 0 and less than 3 parameters")
    }

	endpoint = *endpointFlag
	username = *usernameFlag
	pwd = *pwdFlag

    logger.Get().Printf("endpoint - %v, username - %v, pwd - ***", endpoint, username)
    return count, endpoint, username, pwd, nil
}

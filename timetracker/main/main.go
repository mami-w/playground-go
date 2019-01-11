package main

import (
	"fmt"
	"github.com/mami-w/playground-go/timetracker/other"
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"log"
	"net/http"
	"os"
)


func main() {
	storage, err := trackerdata.NewStorage()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1.0/timetracker/user/", other.UserHandler(storage))
	fmt.Println("starting to listen on port 8000")
	http.ListenAndServe(":8000", nil)
}


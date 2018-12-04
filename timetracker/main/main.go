package main

import (
	"fmt"
	"github.com/mami-w/timetracker/other"
	"log"
	"net/http"
	"os"
)


func main() {
	users, entries, err := other.InitData()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1.0/timetracker/user/", other.UserHandler(users, entries))
	fmt.Println("starting to listen on port 8000")
	http.ListenAndServe(":8000", nil)
}


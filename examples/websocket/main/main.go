package main

import (
	"log"
	"net/http"
)

func main() {

	svr := newServer()
	err := http.ListenAndServe(":8003", svr)
	log.Print("Listening on port 8003")

	if err != nil {
		log.Fatal(err.Error())
	}
}

// from: https://gist.github.com/tmichel/7390690
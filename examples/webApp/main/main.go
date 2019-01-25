package main

import (
	"log"
	"net/http"
)

func main() {
	svr := newServer()
	err := http.ListenAndServe(":5000", svr)
	if err != nil {
		log.Fatal(err.Error())
	}
}

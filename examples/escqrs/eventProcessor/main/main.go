package main

import (
	"github.com/mami-w/playground-go/examples/escqrs/eventProcessor/server"
	"log"
	"net/http"
)

func main() {

	url := "amqp://guest:guest@0.0.0.0:5672"
	svr := server.NewServer(url)

	err := http.ListenAndServe(":8001", svr)
	if err != nil {
		log.Fatal(err.Error())
	}
}

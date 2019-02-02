package main

import (
	"fmt"
	"github.com/mami-w/playground-go/trackerclient-react/server"
	"net/http"
)

func main() {

	svr := server.NewServer();
	fmt.Println("Listening on port 8006")
	panic(http.ListenAndServe(":8006", svr))
}

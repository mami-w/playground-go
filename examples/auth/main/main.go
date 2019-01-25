package main

import (
	"log"
	"net/http"
)

/*
basic auth (user:pwd)
forms/cookie
windows auth
SAML (works with bearer and saml tokens)
AuthO, OpenID
...

 */

 func main() {

 	svr := newServer()
 	err := http.ListenAndServe(":8002", svr)
 	if err != nil {
 		log.Fatal(err.Error())
	}
 }
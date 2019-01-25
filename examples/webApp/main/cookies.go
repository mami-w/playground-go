package main

import (
	"log"
	"net/http"
	"time"
)

func cookieWriteHandler(w http.ResponseWriter, req *http.Request) {
	expire := time.Now().Add(2 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:"sample", Value:"its a mario", Expires:expire }
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusNoContent)
}

func cookieReadHandler(w http.ResponseWriter, req *http.Request) {
	cookie, _ := req.Cookie("sample")
	log.Printf(cookie.Value)
}

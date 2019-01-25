package main

import (
	"encoding/json"
	"net/http"
	"text/template"
)

type sampleContent struct {
		ID string 	`json:"id"`
	Content string 	`json:"content"`
}

func testHandler() http.HandlerFunc {
	content := sampleContent{ID:"0395834", Content:"This is me :)"}
	return func(w http.ResponseWriter, req* http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body,  _ := json.Marshal(content)
		w.Write(body)
		/*WriteHeader sends an HTTP response header with status code. If WriteHeader is not called explicitly,
		the first call to Write will trigger an implicit WriteHeader(http.StatusOK).
		Thus explicit calls to WriteHeader are mainly used to send error codes.
		 */
	}
}

var t *template.Template

func init() {
	t = template.Must(template.ParseFiles("assets/templates/index.html"))
}

func homeHandler() http.HandlerFunc {
	data := sampleContent{ ID:"034300", Content:"another type of content"}
	return func(w http.ResponseWriter, req *http.Request) {
		t.Execute(w, data)
	}
}
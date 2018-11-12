package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"hello":"world"}`))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	for k, v := range r.URL.Query() {
		fmt.Printf("%s, %s - %v\n", k, v, reflect.TypeOf(v))
	}

	fmt.Println(r.URL)

	reqbody, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s\n", reqbody)

	switch r.Method {
	case "GET":
		fmt.Println("get")
	case "PUT":
		fmt.Println("put")
	default:
		fmt.Println(r.Method)
	}

	// there is no good way to get the argument /api/{id} from the query without extra
}
func main() {
	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/api/", apiHandler)
	fmt.Println("starting to listen on port 8000")
	http.ListenAndServe(":8000", nil)
}

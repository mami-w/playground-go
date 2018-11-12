package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func testGet() {
	response, err := http.Get("https://ifconfig.co")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}

func testPut() {

	postData := strings.NewReader(`{"some":"json"}`)
	response, err := http.Post("https://httpbin.org/post", "application/json", postData)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", body)
}

func testClient() {
	debug := os.Getenv("DEBUG")
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://ifconfig.co", nil)
	if nil != err {
		log.Fatal(err)
	}
	if debug == "1" {
		debugrequest, _ := httputil.DumpRequestOut(request, true)
		fmt.Printf("%s", debugrequest)
	}
	response, err := client.Do(request)

	if nil != err {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if debug == "1" {
		debugresponse, _ := httputil.DumpResponse(response, true)
		fmt.Printf("%s", debugresponse)
	}
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
func main() {

	testGet()
	testPut()
	testClient()
}

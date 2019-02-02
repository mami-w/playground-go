package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func getAllEntriesHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"] // the book title slug
		url := trackerEndpoint + fmt.Sprintf(userFormatString, userID)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// todo: handle error
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w,  string(body), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

func getEntryHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func createEntryHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func updateEntryHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func deleteEntryHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

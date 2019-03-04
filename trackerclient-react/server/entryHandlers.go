package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	entryQueryString = "/api/v1.0/timetracker/user/%s"
	entryFormatString = "/api/v1.0/timetracker/user/%s/entry/%s"
)

func getAllEntriesHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		url := trackerEndpoint + fmt.Sprintf(entryQueryString, userID)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, reportRemoteError(string(body)), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

func getEntryHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		entryID := vars["entryid"]

		url := trackerEndpoint + fmt.Sprintf(entryFormatString, userID, entryID)

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, reportRemoteError(string(body)), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

func createEntryHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		entryID := vars["entryid"]

		url := trackerEndpoint + fmt.Sprintf(entryFormatString, userID, entryID)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusCreated {
			http.Error(w, reportRemoteError(string(body)), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

func updateEntryHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		entryID := vars["entryid"]

		url := trackerEndpoint + fmt.Sprintf(entryFormatString, userID, entryID)

		client := &http.Client{}
		req, _ := http.NewRequest("PUT", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, reportRemoteError(string(body)), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

func deleteEntryHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		entryID := vars["entryid"]

		url := trackerEndpoint + fmt.Sprintf(entryFormatString, userID, entryID)

		client := &http.Client{}
		req, _ := http.NewRequest("DELETE", url, r.Body)

		resp, err := client.Do(req)

		if err != nil {
			handleError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			http.Error(w, reportRemoteError(string(body)), resp.StatusCode)
			return
		}

		w.Write(body)
	}
}

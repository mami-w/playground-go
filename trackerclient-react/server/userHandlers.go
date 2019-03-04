package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	userQueryString = "/api/v1.0/timetracker/user"
	userFormatString = "/api/v1.0/timetracker/user/%s"
)

func getAllUsersHandler(trackerEndpoint string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			url := trackerEndpoint + fmt.Sprintf(userQueryString)

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

func createUserHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		url := trackerEndpoint + fmt.Sprintf(userFormatString, userID)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, r.Body)

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

func updateUserHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		url := trackerEndpoint + fmt.Sprintf(userFormatString, userID)

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

func deleteUserHandler(trackerEndpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["userid"]
		url := trackerEndpoint + fmt.Sprintf(userFormatString, userID)

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



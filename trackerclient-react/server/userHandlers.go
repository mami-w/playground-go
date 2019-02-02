package server

import "net/http"

func getAllUsersHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func createUserHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func updateUserHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func deleteUserHandler(trackerEndpoint string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo:
		w.WriteHeader(http.StatusNotImplemented)
	}
}



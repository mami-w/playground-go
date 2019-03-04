package other

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/mami-w/playground-go/timetracker/logger"
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func UserHandler(storage trackerdata.Storage) func (w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r * http.Request) {

		// queryValues := r.URL.Query() -- todo
		path := r.URL.RequestURI()

		// extract the user id
		pattern := regexp.MustCompile(`.*/api/v1\.0/timetracker/user/([^/]*)/?(entry)?/?([^/]*)/?`)

		matches := pattern.FindStringSubmatch(path)

		var userID, entryID string

		if len(matches) == 0 {
			http.Error(w, "Invalid query format", http.StatusBadRequest)
			return
		}

		userID = matches[1]
 		entry := matches[2]
		entryID = matches[3]

		if entry == "" && entryID != "" {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		switch r.Method {

		case "GET":
			handleGet(storage, w, r, userID, entryID)
		case "PUT":
			handlePut(storage, w, r, userID, entryID)
		case "POST":
			handlePost(storage, w, r, userID, entryID)
		case "DELETE":
			handleDelete(storage, w, r, userID, entryID)
		}
	}
}

func handleGet(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if userID == "" {
		handleGetAllUsers(storage, w, r)
		return;
	}
	_, found, err := storage.GetUser(userID)
	if !found {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if entryID == "" {
		handleGetAllEntries(storage, w, r, userID)
	} else {
		handleGetEntry(storage, w, r, userID, entryID)
	}
}

func handleGetAllUsers(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request) {

	var users, err = storage.GetAllUsers();
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(users)

	if err != nil {
		logger.Get().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)
}

func handleGetAllEntries(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string) {

	var entryValues, err = storage.GetAllEntries(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(entryValues)

	if err != nil {
		logger.Get().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)
}

func handleGetEntry(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	entry, found, err := storage.GetEntry(userID, entryID)
	if !found {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		logger.Get().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(entry)

	if err != nil {
		logger.Get().Println(err)
		http.Error(w, "<todo>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)
}

func handlePut(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if entryID == "" {
		handleUpdateUser(storage, w, r, userID)
		return
	}

	handleUpdateEntry(storage, w, r, userID, entryID)
}

func handleUpdateUser(storage trackerdata.Storage, w http.ResponseWriter, req *http.Request, userID string) {

	if userID == "" {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}

	if _, found, _ := storage.GetUser(userID); !found {
		http.Error(w, "User with user ID does not exist", http.StatusBadRequest)
		return
	}

	reqbody, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var user trackerdata.User
	err = json.Unmarshal(reqbody, &user)
	if err != nil {
		e := fmt.Sprintf("Could not read json request body: %v, value: %v", err.Error(), string(reqbody))
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if (user.ID != "") && (user.ID != userID) {
		http.Error(w, "user id in the body does not match the request", http.StatusBadRequest)
		return
	}

	user.ID = userID

	storage.SetUser(user)
	w.WriteHeader(http.StatusOK)
}

func handleUpdateEntry(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	_, exists, _ := storage.GetUser(userID)
	if !exists {
		http.NotFound(w, r)
		return
	}

	_, exists, _ = storage.GetEntry(userID, entryID)
	if !exists {
		http.NotFound(w, r)
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	entry := trackerdata.Entry{}
	err = json.Unmarshal(reqbody, &entry)

	if err != nil {
		e := fmt.Sprintf("Could not read json request body: %v, value: %v", err.Error(), string(reqbody))
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if (entry.UserID != "") && (userID != entry.UserID) {
		http.Error(w, "User Ids do not match", http.StatusBadRequest)
		return
	}
	if (entryID != entry.ID) && (entry.ID != "") {
		http.Error(w, "entry id's do not match", http.StatusBadRequest)
		return
	}

	entry.UserID = userID;
	entry.ID = entryID;

	storage.SetEntry(entry)

	locationURL := createLocationURL(r, userID, entry.ID)

	w.Header().Set("Location", locationURL)

	body, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, "could not write json body", http.StatusInternalServerError)
	}
	w.Write(body)
}

func handlePost(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if entryID == "" {
		handleCreateUser(storage, w, r, userID)
		return
	}

	handleCreateEntry(storage, w, r, userID, entryID)
}

func handleCreateUser(storage trackerdata.Storage, w http.ResponseWriter, req *http.Request, userID string) {

	if userID == "" {
		http.Error(w, "User ID not specified", http.StatusBadRequest)
		return
	}

	if _, found, _ := storage.GetUser(userID); found {
		http.Error(w, "User with user ID already exists", http.StatusBadRequest)
		return
	}

	reqbody, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var user trackerdata.User
	err = json.Unmarshal(reqbody, &user)
	if err != nil {
		e := fmt.Sprintf("Could not read json request body: %v, value: %v", err.Error(), string(reqbody))
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if (user.ID != "") && (user.ID != userID) {
		http.Error(w, "user id in the body does not match the request", http.StatusBadRequest)
		return
	}

	user.ID = userID

	storage.SetUser(user)
	w.WriteHeader(http.StatusCreated)
}

func handleCreateEntry(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	_, exists, _ := storage.GetUser(userID)
	if !exists {
		http.NotFound(w, r)
		return
	}

	reqbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	entry := trackerdata.Entry{}

	err = json.Unmarshal(reqbody, &entry)

	if err != nil {
		e := fmt.Sprintf("could not read json body: %v, %v", err.Error(), string(reqbody))
		http.Error(w, e, http.StatusBadRequest)
		return
	}

	if (entry.UserID != "") && (userID != entry.UserID) {
		http.Error(w, "User Ids do not match", http.StatusBadRequest)
		return
	}
	if (entryID != entry.ID) && (entry.ID != "") {
		http.Error(w, "entry id's do not match", http.StatusBadRequest)
		return
	}

	entry.ID = entryID
	entry.UserID = userID

	//entry.ID, _ = newUUID()

	storage.SetEntry(entry)

	locationURL := createLocationURL(r, userID, entry.ID)

	w.Header().Set("Location", locationURL)

	body, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, "could not write json body", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func createLocationURL(r *http.Request, userID string, entryID string) string {
	url := r.URL
	path := url.Scheme + "://" + url.Host + "/api/v1.0/timetracker/user/" + userID + "/entry/" + entryID
	return path
}

func handleDelete(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if entryID == "" {
		handleDeleteUser(storage, w, r, userID)
		return;
	}

	handleDeleteEntry(storage, w, r, userID, entryID);
}

func handleDeleteUser(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string) {

	_, exists, _ := storage.GetUser(userID)
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	_, err := storage.DeleteUser(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return;
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleDeleteEntry(storage trackerdata.Storage, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	_, exists, _ := storage.GetUser(userID)
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	_, exists, _ = storage.GetEntry(userID, entryID)
	if !exists {
		http.Error(w, "entry not found", http.StatusNotFound)
		return
	}

	_, err := storage.DeleteEntry(userID, entryID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
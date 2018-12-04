package other

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/mami-w/timetracker/trackerdata"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func UserHandler(users UserDataMap, entries EntryDataMap) func (w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r * http.Request) {

		// queryValues := r.URL.Query() -- todo
		path := r.URL.RequestURI()

		// extract the user id
		pattern := regexp.MustCompile(`.*/api/v1\.0/timetracker/user/([^/]*)*/?(entry)?/?([^/]*).*`)

		matches := pattern.FindStringSubmatch(path)

		var userID, entryID string

		if len(matches) == 0 {
			http.NotFound(w, r)
			return
		}

		if len(matches) > 1 {
			userID = matches[1]
		}

		if len(matches) > 3 {
			entryID = matches[3]
		}

		switch r.Method {

		case "GET":
			handleGet(users, entries, w, r, userID, entryID)
		case "PUT":
			handlePut(users, entries, w, r, userID, entryID)
		case "DELETE":
			handleDelete(users, entries, w, r, userID, entryID)
		}
	}
}

func handleGet(users UserDataMap, entries EntryDataMap, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	_, found := users[userID]
	if !found {
		http.NotFound(w, r)
		return
	}

	userEntries := entries[userID]
	var body []byte
	var err error

	if entryID == "" {
		var entryValues []trackerdata.Entry
		for _,v := range userEntries {
			entryValues = append(entryValues, v)
		}
		body, err = json.Marshal(entryValues)
	} else {
		entry, found := userEntries[entryID]
		if !found {
			http.NotFound(w, r)
			return
		}

		body, err = json.Marshal(entry)

	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "<todo>", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(body)
}

func handlePut(users UserDataMap, entries EntryDataMap, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if entryID == "" {
		handlePutUser(users, w, r, userID)
		return
	}

	handlePutEntry(users, entries, w, r, userID, entryID)
}

func handlePutUser(users UserDataMap, w http.ResponseWriter, r *http.Request, userID string) {

	// unmarshal the body
	reqbody, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}

	var user trackerdata.User
	err = json.Unmarshal(reqbody, &user)
	if err != nil {
		http.Error(w, "Could not read request body json", http.StatusBadRequest)
		return
	}

	if userID != user.ID {
		http.Error(w, "UserID does not match body", http.StatusBadRequest)
		return
	}

	status := http.StatusOK
	// find out if user does already exist
	if _, found := users[userID]; !found {
		status = http.StatusCreated
	}

	users[userID] = user
	w.WriteHeader(status)
}

func handlePutEntry(users UserDataMap, entries EntryDataMap, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	_, exists := users[userID]
	if !exists {
		http.NotFound(w, r)
		return
	}

	userEntries, exists := entries[userID]
	if !exists {
		userEntries = map[string]trackerdata.Entry{}
		entries[userID] = userEntries
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
		http.Error(w, "could not read json body", http.StatusBadRequest)
		return
	}

	if userID != entry.UserID {
		http.Error(w, fmt.Sprintf("user IDs don't match, user: %v - entry: %v", userID, entry.UserID,), http.StatusBadRequest)
		return
	}
	if entryID != entry.ID {
		http.Error(w, fmt.Sprintf("entry IDs don't match, entry: %v - ID: %v", entryID, entry.ID,), http.StatusBadRequest)
		return
	}

	if entry.ID == "" {
		entry.ID, _ = newUUID()
	}

	_, exists = userEntries[entry.ID]
	if exists == false {
		w.WriteHeader(http.StatusCreated)
	}

	userEntries[entry.ID] = entry

	locationURL := createLocationURL(r, userID, entry.ID)

	w.Header().Set("Location", locationURL)

	body, err := json.Marshal(entry)
	if err != nil {
		http.Error(w, "could not write json body", http.StatusInternalServerError)
	}
	w.Write(body)
}

func createLocationURL(r *http.Request, userID string, entryID string) string {
	url := r.URL
	path := url.Scheme + "://" + url.Host + "/api/v1.0/timetracker/user/" + userID + "/entry/" + entryID
	return path
}

func handleDelete(users UserDataMap, entries EntryDataMap, w http.ResponseWriter, r *http.Request, userID string, entryID string) {

	if userID == "" {
		http.Error(w, "no user specified", http.StatusNotFound)
		return
	}

	_, exists := users[userID]
	if !exists {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	entryMap, exists := entries[userID]
	if !exists {
		http.Error(w, "no entries for user found", http.StatusNotFound)
		return
	}

	if entryID == "" {
		entries[userID] = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, exists = entryMap[entryID]
	if !exists {
		http.Error(w, "entry not found", http.StatusNotFound)
		return
	}

	delete(entryMap, entryID)
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
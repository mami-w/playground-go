package other

import (
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetUser1(t *testing.T) {

	getResponse := getResponseTestData(t)
	rr := getResponse("GET", "/api/v1.0/timetracker/user/1", nil, t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `
		[{
    		"id": "a",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		},
		{
    		"id": "b",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}]
`
	expectedFormatted := strings.Join(strings.Fields(expected), "")
	// todo: array could come in different order
	if rr.Body.String() != expectedFormatted {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedFormatted)
	}
}

func TestGetUser1EntryA(t *testing.T) {

	getResponse := getResponseTestData(t)
	rr := getResponse("GET", "/api/v1.0/timetracker/user/1/entry/a", nil, t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `
		{
    		"id": "a",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
`
	expectedFormatted := strings.Join(strings.Fields(expected), "")
	if rr.Body.String() != expectedFormatted {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedFormatted)
	}
}

func TestGetNonExistingEntry(t *testing.T) {

	getResponse := getResponseTestData(t)
	rr := getResponse("GET", "/api/v1.0/timetracker/user/1/entry/c", nil, t)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestAddUser1(t *testing.T) {

	jsonEntry := `{ "ID":""}`
	body := strings.NewReader(jsonEntry)

	getResponse := getResponseEmpty()
	rr := getResponse("POST", "/api/v1.0/timetracker/user/1", body, t)

	responsebody, _ := ioutil.ReadAll(rr.Body)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusCreated, responsebody)
	}
}

func TestAddUserAndEntry(t *testing.T) {

	jsonEntry := `{ "id":""}`
	body := strings.NewReader(jsonEntry)

	getResponse := getResponseEmpty()
	rr := getResponse("POST", "/api/v1.0/timetracker/user/1", body, t)

	responsebody, _ := ioutil.ReadAll(rr.Body)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusCreated, responsebody)
	}

	jsonEntry = `{
    		"id": "",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}`
	body = strings.NewReader(jsonEntry)

	rr = getResponse("POST", "/api/v1.0/timetracker/user/1/entry/a", body, t)

	responsebody, _ = ioutil.ReadAll(rr.Body)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusCreated, responsebody)
	}
}

func TestAddEntryWrongFormat(t *testing.T) {

	jsonEntry := `{
    		"id": "b",
    		"userID": "1" -- wrong format
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}`
	body := strings.NewReader(jsonEntry)

	getResponse := getResponseTestData(t)
	rr := getResponse("PUT", "/api/v1.0/timetracker/user/1/entry/b", body, t)

	responsebody, _ := ioutil.ReadAll(rr.Body)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusBadRequest, responsebody)
	}
}

func TestUpdateEntry(t *testing.T) {

	jsonEntry := `{
    		"id": "a",
    		"userID": "1",
    		"entryType": "2",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}`
	body := strings.NewReader(jsonEntry)

	getResponse := getResponseTestData(t)
	rr := getResponse("PUT", "/api/v1.0/timetracker/user/1/entry/a", body, t)

	responsebody, _ := ioutil.ReadAll(rr.Body)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v, error: %s",
			status, http.StatusOK, responsebody)
	}
}

func TestDeleteEntry(t *testing.T) {
	getResponse := getResponseTestData(t)

	rr := getResponse("DELETE", "/api/v1.0/timetracker/user/1/entry/a", nil, t)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	rr = getResponse("GET", "/api/v1.0/timetracker/user/1/entry/b", nil, t)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteAllEntries(t *testing.T) {
	getResponse := getResponseTestData(t)

	rr := getResponse("DELETE", "/api/v1.0/timetracker/user/1", nil, t)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}

	rr = getResponse("GET", "/api/v1.0/timetracker/user/1/entry/b", nil, t)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

}

func TestDeleteNonExistingEntry(t *testing.T) {
	getResponse := getResponseTestData(t)

	rr := getResponse("DELETE", "/api/v1.0/timetracker/user/1/entry/c", nil, t)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func getResponseEmpty() func (method string, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {

	storage, _ := trackerdata.NewStorage()

	return func (method string, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {
		return getHttpResponse(method, url, body, storage, t)
	}
}

func getResponseTestData(t *testing.T) func (method string, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {

	storage, err := AddTestData()
	if err != nil {
		t.Error(err)
	}
	return func (method string, url string, body io.Reader, t *testing.T) *httptest.ResponseRecorder {
		return getHttpResponse(method, url, body, storage, t)
	}
}


func getHttpResponse(method string, url string, body io.Reader, storage trackerdata.Storage, t *testing.T) *httptest.ResponseRecorder {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UserHandler(storage))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	return rr
}
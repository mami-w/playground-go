package postgresStorage

import (
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"strconv"
	"testing"
	"time"
)

func NewAWSStorage() (storage *PostgresStorage, err error) {

	user := "sa"
	pwd := "passw0rd!"
	endpoint := "postgres-test.c7dqrtwbjpzv.us-east-2.rds.amazonaws.com"

	return NewPostgresStorage(endpoint, user, pwd)
}

func TestPostgresStorage_GetUser(t *testing.T) {

	storage, _ := NewAWSStorage()
	err := storage.TestSetUser("1")

	if err != nil {
		t.Error(err)
	}

	user, found := storage.TestGetUser("1")

	if !found {
		t.Error("Did not find user '1' ")
	}
	if user.ID != "1" {
		t.Errorf("Wrong id, expected: 1, got %v", user.ID)
	}
}

func TestPostgresStorage_SetUser(t *testing.T) {

	storage, _ := NewAWSStorage()
	err := storage.TestSetUser("1")

	if err != nil {
		t.Error(err)
	}

}

func TestPostgresStorage_DeleteUser(t *testing.T) {

	storage, _ := NewAWSStorage()
	err := storage.TestSetUser("1")

	_, err = storage.DeleteUser("1")

	if err != nil {
		t.Error(err.Error())
	}

	users, err := storage.GetAllUsers();

	for _, user := range users {
		if user.ID == "1" {
			t.Error("user 1 not deleted")
		}
	}

}

func TestPostgresStorage_GetAllUser(t *testing.T) {

	storage,_ := NewAWSStorage()
	err  := storage.TestSetUser("1")
	if err != nil {
		t.Error(err)
	}
	err = storage.TestSetUser("2")
	if err != nil {
		t.Error(err)
	}

	users, err := storage.GetAllUsers()

	if len(users) < 2 {
		t.Error("not enough users")
	}

	for i := 1; i < 3; i++ {
		id := strconv.Itoa(i)
		if !userExists(users, id) {
			t.Errorf("cannot find user %v", i)
		}
 	}
}

func TestPostgresStorage_SetEntry(t *testing.T) {

	storage,_ := NewAWSStorage()

	storage.TestSetUser("1")
	err := storage.TestSetEntry("abc", "1")

	if err != nil {
		t.Error(err)
	}
}

func TestPostgresStorage_GetEntry(t *testing.T) {

	storage,_ := NewAWSStorage()

	storage.TestSetUser("1")
	storage.TestSetEntry("abc", "1")

	entry, found, err := storage.GetEntry("1", "abc")

	if err != nil {
		t.Error(err)
	}
	if !found {
		t.Error("did not find entry 'abc' with user '1'")
	}
	if entry.ID != "abc" || entry.UserID != "1" {
		t.Errorf("expected entry 1-abc, got %v-%v", entry.UserID, entry.ID)
	}
}

func TestPostgresStorage_DeleteEntry(t *testing.T) {

	storage,_ := NewAWSStorage()

	storage.TestSetUser("1")
	storage.TestSetEntry("abc", "1")

	_, err := storage.DeleteEntry("1", "abc")

	if err != nil {
		t.Error(err.Error())
	}
}

func TestPostgresStorage_GetAllEntries(t *testing.T) {

	storage,_ := NewAWSStorage()

	storage.TestSetUser("1")
	storage.TestSetEntry("abc", "1")
	storage.TestSetEntry("def", "1")

	entries, err := storage.GetAllEntries("1")

	if err != nil {
		t.Error(err)
	}

	if len(entries) < 2 {
		t.Error("found less than 2 entries")
	}

	e := [2]string{"abc", "def"}
	for _, v := range e {
		if !entryExists(entries, v) {
			t.Error("could not find entry %v", v)
		}
	}
}

func TestPostgresStorage_DeleteAllEntries(t *testing.T) {

	storage,_ := NewAWSStorage()

	storage.TestSetUser("1")
	storage.TestSetEntry("abc", "1")

	_, err := storage.DeleteUser("1")

	if err != nil {
		t.Error(err.Error())
	}
}

func entryExists(entries []trackerdata.Entry, id string) bool {
	for _, entry := range entries {
		if entry.ID == id {
			return true
		}
	}
	return false
}

func (storage *PostgresStorage) TestSetUser(id string) (err error){

	user := trackerdata.User{ID:id}
	_, err = storage.SetUser(user)
	return err
}

func (storage *PostgresStorage) TestGetUser(id string) (user *trackerdata.User, found bool){

	user, found, _ = storage.GetUser(id)
	return user, found
}

func (storage *PostgresStorage) TestSetEntry(id string, userID string) (err error) {

	duration, _ := time.ParseDuration("1h")
	entry := trackerdata.Entry{ID:id, UserID:userID, EntryType:"1", StartTime:time.Now(), Length: duration}
	_, err = storage.SetEntry(entry)
	return err
}

func userExists(users []trackerdata.User, id string) bool {
	for _, u := range users {
		if u.ID == id {
			return true
		}
	}
	return false
}
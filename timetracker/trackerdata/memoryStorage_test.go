package trackerdata

import (
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {

	s, _:= NewStorage()
	if s == nil {
		t.Error("Could not create MemoryStorage")
	}
}

func TestMemoryStorage_GetUser(t *testing.T) {

	s := newTestMemoryStorage()
	user, found, err := s.GetUser("1")

	if !found {
		t.Error(err)
	}

	if user.ID != "1" {
		t.Errorf("User has id %v instead of '1'", user.ID)
	}

}

func TestMemoryStorage_SetUser(t *testing.T) {

	s := newTestMemoryStorage()
	s.SetUser(User{ID:"1"})

	user, found := s.users["1"]

	if !found {
		t.Error("Did not find user 1")
	}

	if user.ID != "1" {
		t.Errorf("Expected user: 1, got user: %v", user.ID)
	}
}

func TestMemoryStorage_GetEntry(t *testing.T) {

	s := newTestMemoryStorage()

	entry, found, err := s.GetEntry("1", "a")

	if !found || err != nil {
		t.Error(err)
	}

	if entry.UserID != "1" {
		t.Errorf("Expected user 1, got user: %v", entry.UserID)
	}

	if entry.ID != "a" {
		t.Errorf("Expected entry a, got entry: %v", entry.ID)
	}

}

func TestMemoryStorage_SetEntry(t *testing.T) {

	duration, _ := time.ParseDuration("1h")
	entry := Entry{ID:"b", UserID:"2", EntryType:"1", StartTime:time.Now(), Length: duration}
	s, _ := NewStorage()
	_, err := s.SetEntry(entry)

	if err == nil {
		t.Error("Expected failure to add entry for non existing user")
	}

	status, _ := s.SetUser(User{ID:"2"})

	if status != StatusCreated {
		t.Errorf("Expected 'created', got %v", status)
	}

	status, _ = s.SetEntry(entry)

	if status != StatusCreated {
		t.Errorf("Expected 'created, got %v", status)
	}

	entry = s.entries["2"]["b"]

	if entry.ID != "b" {
		t.Errorf("Expected 'b', got %v", entry.ID)
	}

	if entry.UserID != "2" {
		t.Errorf("Expected '2', bog %v", entry.UserID)
	}
}

func TestMemoryStorage_GetAllUsers(t *testing.T) {

	s := newTestMemoryStorage()

	users, _ := s.GetAllUsers()

	if len(users) != 1 {
		t.Error("No users found")
	}

	if users[0].ID != "1" {
		t.Errorf("expected '1', got %v", users[0].ID)
	}
}

func TestMemoryStorage_GetAllEntries(t *testing.T) {

	s := newTestMemoryStorage()

	entries, _ := s.GetAllEntries("1")

	if (len(entries) != 1) {
		t.Error("No entries found")
	}

	entry := entries[0]
	if entry.ID  != "a" || entry.UserID != "1" {
		t.Errorf("expected 1-a, got %v - %v", entry.UserID, entry.ID)
	}
}

func newTestMemoryStorage() (s* MemoryStorage) {

	duration, _ := time.ParseDuration("1h")
	s, _ = NewStorage()

	s.users["1"] = User{ID: "1"}
	s.entries["1"] = make(map[string]Entry)
	s.entries["1"]["a"] = Entry{ID: "a",
		UserID:    "1",
		EntryType: "1",
		StartTime: time.Now(),
		Length: duration,
	}

	return s
}
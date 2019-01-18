package memoryStorage

import (
	"fmt"
	"log"
)
import "github.com/mami-w/playground-go/timetracker/trackerdata"

type UserDataMap map[string]trackerdata.User
type EntryDataMap map[string]map[string]trackerdata.Entry

type MemoryStorage struct {
	users UserDataMap
	entries EntryDataMap
	logger *log.Logger // todo: not used yet
}

type Error string

func NewError(s string) Error {
	return Error(s)
}

func (e Error) Error() string {
	return string(e)
}

func NewStorage() (s *MemoryStorage, err error) {

	err = nil

	s = &MemoryStorage{
		users:   UserDataMap{},
		entries: EntryDataMap{},
	}
	return s, err
}

func (s* MemoryStorage) SetUser(u trackerdata.User) (status trackerdata.Status, err error) {

	status = trackerdata.StatusCreated
	err = nil

	_, found := s.users[u.ID]
	if found {
		status = trackerdata.StatusUpdated
	}

	s.users[u.ID] = u;

	return status, err
}

func (s* MemoryStorage) GetUser(id string) (u *trackerdata.User, found bool, err error) {

	err = nil

	user, found := s.users[id]

	return &user, found, nil
}

func (s* MemoryStorage) GetAllUsers() (users []trackerdata.User, err error) {

	err = nil
	for _, v := range s.users {
		users = append(users, v)
	}

	return users, err
}

func (s* MemoryStorage) SetEntry(e trackerdata.Entry) (status trackerdata.Status, err error) {

	err = nil

	_, found := s.users[e.UserID]
	if !found {
		err = NewError(fmt.Sprintf("User %v does not exist", e.UserID))
		return status, err
	}

	status = trackerdata.StatusUpdated
	entryMap, found := s.entries[e.UserID]
	if !found {
		entryMap = make(map[string]trackerdata.Entry)
		s.entries[e.UserID] = entryMap
		status = trackerdata.StatusCreated
	}

	if _, found = entryMap[e.ID]; found {
		status = trackerdata.StatusUpdated
	}

	entryMap[e.ID] = e

	return status, err
}

func (s* MemoryStorage) GetEntry(userID string, id string) (e *trackerdata.Entry, found bool, err error) {

	err = nil

	if _, found = s.users[userID]; !found {
		err = NewError(fmt.Sprintf("User %v does not exist", userID))
		return e, found, err
	}

	entries, found := s.entries[userID];
	if  !found {
		return e, found, err
	}

	entry, found := entries[id]

	return &entry, found, err
}

func (s* MemoryStorage) GetAllEntries(userID string) (entries []trackerdata.Entry, err error) {

	err = nil

	_, found, err := s.GetUser(userID)
	if !found {
		err = NewError(fmt.Sprintf("User not found %v", userID))
		return entries, err
	}

	entryMap, found := s.entries[userID]
	if !found {
		return entries, err
	}
	for _, v := range entryMap {
		entries = append(entries, v)
	}

	return entries, err
}

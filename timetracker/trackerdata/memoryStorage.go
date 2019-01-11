package trackerdata

import "fmt"

type UserDataMap map[string]User
type EntryDataMap map[string]map[string]Entry

type MemoryStorage struct {
	users UserDataMap
	entries EntryDataMap
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

func (s* MemoryStorage) SetUser(u User) (status Status, err error) {

	status = StatusCreated
	err = nil

	_, found := s.users[u.ID]
	if found {
		status = StatusUpdated
	}

	s.users[u.ID] = u;

	return status, err
}

func (s* MemoryStorage) GetUser(id string) (u User, found bool, err error) {

	err = nil

	u, found = s.users[id]

	return u, found, nil
}

func (s* MemoryStorage) GetAllUsers() (users []User, err error) {

	err = nil
	for _, v := range s.users {
		users = append(users, v)
	}

	return users, err
}

func (s* MemoryStorage) SetEntry(e Entry) (status Status, err error) {

	err = nil

	_, found := s.users[e.UserID]
	if !found {
		err = NewError(fmt.Sprintf("User %v does not exist", e.UserID))
		return status, err
	}

	status = StatusUpdated
	entryMap, found := s.entries[e.UserID]
	if !found {
		entryMap = make(map[string]Entry)
		s.entries[e.UserID] = entryMap
		status = StatusCreated
	}

	if _, found = entryMap[e.ID]; found {
		status = StatusUpdated
	}

	entryMap[e.ID] = e

	return status, err
}

func (s* MemoryStorage) GetEntry(userID string, id string) (e Entry, found bool, err error) {

	err = nil

	if _, found = s.users[userID]; !found {
		err = NewError(fmt.Sprintf("User %v does not exist", userID))
		return e, found, err
	}

	entries, found := s.entries[userID];
	if  !found {
		return e, found, err
	}

	e, found = entries[id]

	return e, found, err
}

func (s* MemoryStorage) GetAllEntries(userID string) (entries []Entry, err error) {

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

package other

import (
	"encoding/json"
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"github.com/mami-w/playground-go/timetracker/trackerdata/memoryStorage"
)

func AddTestData() (storage trackerdata.Storage, err error) {

	storage, _ = memoryStorage.NewStorage()
	user := trackerdata.User{}

	err = json.Unmarshal([]byte(user1), &user)
	if err != nil {
		return nil, err
	}
	storage.SetUser(user)

	err = json.Unmarshal([]byte(user2), &user)
	storage.SetUser(user)

	entry := trackerdata.Entry{}

	err = json.Unmarshal([]byte(entry1a), &entry)
	storage.SetEntry(entry)
	err = json.Unmarshal([]byte(entry1b), &entry)
	storage.SetEntry(entry)
	err = json.Unmarshal([]byte(entry2c), &entry)
	storage.SetEntry(entry)
	err = json.Unmarshal([]byte(entry2d), &entry)
	storage.SetEntry(entry)

	return storage, err
}

const (
	user1 = `{ "ID":"1"}`
	user2 = `{ "ID":"2"}`
	entry1a = `
		{
    		"id": "a",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
`
	entry1b = `
		{
    		"id": "b",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
`
	entry2c = `
	{
    		"id": "c",
    		"userID": "2",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
`
	entry2d = `
		{
    		"id": "d",
    		"userID": "2",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
`
)

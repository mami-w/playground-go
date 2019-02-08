package other

import (
	"encoding/json"
	"fmt"
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

func Test1() {
	entry := trackerdata.Entry{}
	err := json.Unmarshal([]byte(entryGuid), &entry)
	fmt.Println(entry)
	fmt.Println(err.Error())
	fmt.Println("done")
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

	entryGuid = `
{
"entryType": "1",
"id": "7a11dae1-e533-463f-a2f0-e8c80e9eede9",
"length": 3600000000000,
"startTime": "2019-02-06T21:36:44.496Z",
"userid": "4b680358-015d-4d70-bf4c-4f28145e7b5b"
}
`
)

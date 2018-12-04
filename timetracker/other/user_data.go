package other

import (
	"encoding/json"
	"github.com/mami-w/timetracker/trackerdata"
)


type UserDataMap map[string]trackerdata.User
type EntryDataMap map[string]map[string]trackerdata.Entry

func InitData() (users	 UserDataMap, entries EntryDataMap, err error) {

	users = UserDataMap{}
	entries = EntryDataMap{}

	return users, entries, nil
}
func AddTestData() (userCache UserDataMap, entryCache EntryDataMap, err error) {

	err = json.Unmarshal([]byte(userdata), &userCache)
	if err != nil {
		return nil, nil, err
	}

	return userCache, entryCache, json.Unmarshal([]byte(entrydata), &entryCache)
}

const (
	userdata = `
{
	"1": { "ID":"1"},
	"2": { "ID":"2"}
}
`

	entrydata = `
{
	"1":{
		"a":{
    		"id": "a",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		},
		"b":{
    		"id": "b",
    		"userID": "1",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
	},
	"2":{
		"c":{
    		"id": "c",
    		"userID": "2",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		},
		"d":{
    		"id": "d",
    		"userID": "2",
    		"entryType": "1",
    		"startTime": "2018-12-02T19:14:53.785417-08:00",
    		"length": 90000000000
		}
	}
}
`
)

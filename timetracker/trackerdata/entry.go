package trackerdata

import "time"

// Entry is a time entry
type Entry struct {
	ID        string        `json:"id"`
	UserID    string        `json:"userid"`
	EntryType string        `json:"entryType"`
	StartTime time.Time     `json:"startTime"`
	Length    time.Duration `json:"length"`
}

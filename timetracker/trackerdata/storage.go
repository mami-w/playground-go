package trackerdata

type Status int

const (
	StatusCreated   Status = 0
	StatusUpdated   Status = 1
)

type Storage interface {
	SetUser(u User) (Status, error)
	GetUser(id string) (*User, bool, error)
	DeleteUser(id string) (bool, error)
	GetAllUsers() ([]User, error)

	SetEntry(e Entry) (Status, error)
	GetEntry(userID string, id string) (*Entry, bool, error)
	DeleteEntry(userID string, id string) (bool, error)
	GetAllEntries(userID string) ([]Entry, error)
}

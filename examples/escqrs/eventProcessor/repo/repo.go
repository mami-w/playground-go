package repo

import "log"

type Repository struct {
	// todo
}

type RepositoryInterface interface {
	Save(event interface{}) (err error)
}

func InitRepository() (repo *Repository) {
	repo = &Repository{}
	return repo
}

func (repo *Repository) Save(data interface{}) (err error) {
	// nothing...
	log.Printf("%+v", data)
	return nil
}
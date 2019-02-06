package postgresStorage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mami-w/playground-go/timetracker/logger"
	"github.com/mami-w/playground-go/timetracker/trackerdata"
	"log"
)

/*
type Storage interface {
	SetUser(u User) (Status, error)
	GetUser(id string) (User, bool, error)
	GetAllUsers() ([]User, error)

	SetEntry(e Entry) (Status, error)
	GetEntry(userID string, id string) (Entry, bool, error)
	GetAllEntries(userID string) ([]Entry, error)
}
*/

type PostgresStorage struct {
	endpoint string
	user string
	pwd string
	dnsStr string
}

const (
	dbname = "trackerdata"
	driverName = "postgres"
	dnsFormatString = "postgres://%s:%s@%s/%s"

	queryUsersSql = `SELECT id from data.users`

	queryUserSql = `SELECT id from data.users where id = $1`

	queryEntriesSql = `
SELECT id, userid, entrytype, starttime, duration
FROM data.entries
WHERE userid = $1
`

	queryEntrySql = `
SELECT id, userid, entrytype, starttime, duration
FROM data.entries
WHERE id = $1 and userid = $2
`

	upsertUserSql = `
INSERT INTO data.users (id)
VALUES ($1)
ON CONFLICT (id)
DO NOTHING
RETURNING id`

	queryEntryExistSql = `SELECT count(1) FROM data.entries WHERE id = $1;`

	upsertEntrySql = `
INSERT INTO data.entries (id, userid, entrytype, starttime, duration)
VALUES($1, $2, $3, $4, $5)
ON CONFLICT (id)
DO UPDATE SET entrytype = EXCLUDED.entrytype, starttime = EXCLUDED.starttime, duration = EXCLUDED.duration;
`

	deleteUserSql = `
DELETE FROM data.users
WHERE id = $1
`

	deleteAllEntriesSql = `
DELETE FROM data.entries
WHERE userid = $1
`

	deleteEntrySql = `
DELETE FROM data.entries
WHERE id = $1 AND userid = $2
`
	)

func NewPostgresStorage(endpoint string, user string, pwd string) (storage *PostgresStorage, err error) {

	dnsStr := fmt.Sprintf(dnsFormatString, user, pwd, endpoint, dbname)
	storage = &PostgresStorage{endpoint:endpoint, user:user, pwd:pwd, dnsStr:dnsStr}
	// todo additional validation

	return storage, err
}

func (storage *PostgresStorage) SetUser(u trackerdata.User) (status trackerdata.Status, err error) {

	// get a db connection
	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()

	var updated string
	err = db.QueryRow(upsertUserSql, u.ID).Scan(&updated)

	switch {
		  	case err == sql.ErrNoRows:
		  		status = trackerdata.StatusUpdated
		  		err = nil
		  	case err != nil:
		  		log.Fatal(err)
		  	default:
		  		status = trackerdata.StatusCreated
		  	}

	return status, err
}

func (storage *PostgresStorage) GetUser(id string) (user *trackerdata.User, found bool, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			fmt.Println(err.Error())
		}
	}()

	if err != nil {
		return user, false, err
	}

	var userID string
	err = db.QueryRow(queryUserSql, id).Scan(&userID)

	if (err == sql.ErrNoRows) {
		return nil, false, nil
	}

	user = &trackerdata.User{ ID:userID }

	return user, true, nil
}

func (storage *PostgresStorage) DeleteUser(id string) (success bool, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}

	{
		stmt, err := tx.Prepare(deleteAllEntriesSql)
		if err != nil {
			tx.Rollback()
			return false, err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(id); err != nil {
			tx.Rollback() // return an error too, we may want to wrap them
			return false, err
		}
	}

	{
		stmt, err := tx.Prepare(deleteUserSql)
		if err != nil {
			tx.Rollback()
			return false, err
		}
		defer stmt.Close()

		if _, err := stmt.Exec(id); err != nil {
			tx.Rollback() // return an error too, we may want to wrap them
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (storage *PostgresStorage) GetAllUsers() (users []trackerdata.User, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	if err != nil {
		return  nil, err
	}

	rows, err := db.Query(queryUsersSql)
	if err != nil {
		return nil, err
	}

	users = make([]trackerdata.User,0)
	var id string
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		user := trackerdata.User{ID:id}
		users = append(users, user)
	}

	return users, err
}

func (storage *PostgresStorage) SetEntry(e trackerdata.Entry) (status trackerdata.Status, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	var entrycount int
	err = db.QueryRow(queryEntryExistSql, e.ID).Scan(&entrycount)
	if err != nil {
		return status, err
	}

	_, err = db.Exec(upsertEntrySql, e.ID, e.UserID, e.EntryType, e.StartTime, int64(e.Length))

	if err != nil {
		return status, err
	}

	if entrycount >  0 {
		return trackerdata.StatusUpdated, err
	}

	return trackerdata.StatusCreated, err
}

func (storage *PostgresStorage) GetEntry(userID string, id string) (*trackerdata.Entry, bool, error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	entry := trackerdata.Entry{}
	err = db.QueryRow(queryEntrySql, id, userID).Scan(&entry.ID, &entry.UserID, &entry.EntryType, &entry.StartTime, &entry.Length)

	switch {
	case err == sql.ErrNoRows:
		return nil, false, nil
	case err != nil:
		return nil, false, err
	}

	return &entry, true, err

}

func (storage *PostgresStorage) DeleteEntry(userID string, id string) (success bool, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	var entrycount int
	err = db.QueryRow(deleteEntrySql, id, userID).Scan(&entrycount)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (storage *PostgresStorage) GetAllEntries(userID string) (entries []trackerdata.Entry, err error) {

	db, err := sql.Open(driverName, storage.dnsStr)

	defer db.Close()
	defer func() {
		if err != nil {
			logger.Get().Println(err.Error())
		}
	}()

	rows, err := db.Query(queryEntriesSql, userID)

	switch {
	case err == sql.ErrNoRows:
		return entries, nil
	case err != nil:
		return entries, err
	}

	entries = make([]trackerdata.Entry,0)
	var entry trackerdata.Entry
	for rows.Next() {
		err = rows.Scan(&entry.ID, &entry.UserID, &entry.EntryType, &entry.StartTime, &entry.Length)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, err
}

/*
INSERT INTO customers (name, email)
VALUES
 (
 'Microsoft',
 'hotline@microsoft.com'
 )
ON CONFLICT (name)
DO
 UPDATE
   SET email = EXCLUDED.email || ';' || customers.email;
 */


 /*
 sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)
RETURNING id`
  id := 0
  err = db.QueryRow(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun").Scan(&id)
  */


  /*
   tx, err := db.Begin()
    if err != nil {
        return err
    }

    {
        stmt, err := tx.Prepare(`INSERT INTO table_1 (thing_1, whatever)
                     VALUES($1,$2);`)
        if err != nil {
            tx.Rollback()
            return err
        }
        defer stmt.Close()

        if _, err := stmt.Exec(thing_1, whatever); err != nil {
            tx.Rollback() // return an error too, we may want to wrap them
            return err
        }
    }

    {
        stmt, err := tx.Prepare(`INSERT INTO table_2 (thing_2, whatever)
                     VALUES($1, $2);`)
        if err != nil {
            tx.Rollback()
            return err
        }
        defer stmt.Close()

        if _, err := stmt.Exec(thing_2, whatever); err != nil {
            tx.Rollback() // return an error too, we may want to wrap them
            return err
        }
    }

    return tx.Commit()
}
   */
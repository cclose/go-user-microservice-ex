package database

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

type PostGresDB struct {
	PgDbSession      *sql.DB
	connectionString string
}

// ErrDuplicateKey is used to signify a duplicate key error
var ErrDuplicateKey = errors.New("database.postgres.duplicatekey")
var duplicateKeyErrRegex = regexp.MustCompile("pq: duplicate key value violates unique constraint \"([^\"]*)\"")

// DuplicateKeyError tests an error to see if it is a Duplicate Key Error
func DuplicateKeyError(err error) bool {
	return duplicateKeyErrRegex.MatchString(err.Error())
}

func Connect(host, user, pass, db, port string) (*PostGresDB, error) {

	connectString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, db)

	dbh, err := sql.Open("postgres", connectString)
	if err != nil {
		return nil, err
	} else if dbh == nil {
		return nil, errors.New("postgres Open returned nil")
	}

	err = dbh.Ping()
	if err != nil {
		return nil, err
	}

	pgdbh := &PostGresDB{dbh, connectString}

	return pgdbh, nil
}

// Disconnect close down our database connection
func (pgdbh *PostGresDB) Disconnect() {
	_ = pgdbh.PgDbSession.Close()
}

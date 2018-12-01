package db

import (
	"database/sql"
	"fmt"

	// postgres db
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/odacremolbap/rest-demo/pkg/log"
)

const postgresDriver = "postgres"

// PersistenceManager exposes entities persistence methods
// for TODO list
type PersistenceManager struct {
	db *sql.DB
}

// Manager is the global reference for persistence methods
// and should be initialized at app start
var Manager *PersistenceManager

// ConnectPostgressDB returns a connected db object
func ConnectPostgressDB(host string, port int, user, pass, database string, dbSSL bool) (*sql.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%d user=%s dbname=%s",
		host, port, user, database)
	if dbSSL {
		dataSource += " sslmode=require"
	} else {
		dataSource += " sslmode=disable"
	}
	log.V(5).Info("connecting to database",
		"datasource (password not shown)", dataSource)

	dataSource += fmt.Sprintf(" password=%s", pass)

	db, err := sql.Open(postgresDriver, dataSource)
	if err != nil {
		return nil, errors.Wrapf(err, "user %s couldn't open database %s:%d/%s",
			user, host, port, database)
	}

	log.V(5).Info("pinging database", "database", database)
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, errors.Wrapf(err, "error pinging database %s:%d/%s",
			host, port, database)
	}
	return db, nil
}

// NewTODOPersistenceManager returns a persistence manager for TODO
func NewTODOPersistenceManager(db *sql.DB) *PersistenceManager {
	return &PersistenceManager{db: db}
}

package storage

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationsPath = "file://migrations"
const defaultDbPath = "uptime-monitor.db"

type SQLiteStorer struct {
	db *sqlx.DB
	dbPath string
}

func NewSQLiteStorer(dataSourceName string) (*SQLiteStorer, error) {
	if dataSourceName == "" {
		dataSourceName = defaultDbPath
	}

    db, err := sqlx.Connect("sqlite3", dataSourceName)
    if err != nil {
        return nil, err
    }

    storer := &SQLiteStorer{db: db, dbPath: dataSourceName}
    if err := storer.runMigrations(); err != nil {
        return nil, err
    }

    return storer, nil
}

func (s *SQLiteStorer) runMigrations() error {
    m, err := migrate.New(
        "file://"+migrationsPath, // file path to migration files
        fmt.Sprintf("sqlite://%s", s.dbPath), // database URL
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}

func (s *SQLiteStorer) StoreWebsiteStatus(website string, status string, latency time.Duration) error {
	query := `INSERT INTO website_status (website, status, lentency) VALUES (?, ?, ?)`
	
	latencyMs := latency.Milliseconds()

	_, err := s.db.Exec(query, website, status, latencyMs)
	if(err != nil) {
		return err
	}

	return nil
}
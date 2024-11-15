package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Service interface {
	Health() error
	Migrate() error
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	dbUrl      = os.Getenv("BLUEPRINT_DB_URL")
	authToken  = os.Getenv("BLUEPRINT_DB_AUTH_TOKEN")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sql.Open("libsql", dbUrl+"?authToken="+authToken)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Health() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	return nil
}

func (s *service) Close() error {
	log.Printf("Disconnected from database")
	return s.db.Close()
}

func (s *service) Migrate() error {
	// Create a migrations table if it doesn't exist
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return err
	}

	migrations := []struct {
		name string
		stmt string
	}{
		{
			name: "create players",
			stmt: `
        CREATE TABLE IF NOT EXISTS players (
          id INTEGER PRIMARY KEY AUTOINCREMENT,
          nickname TEXT NOT NULL,
          champion TEXT NOT NULL,
          frag TEXT NOT NULL,
          date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )`,
		},
	}

	for _, m := range migrations {
		var count int
		err = s.db.QueryRow("SELECT COUNT (*) FROM migrations WHERE name = ?", m.name).Scan(&count)
		if err != nil {
			return err
		}

		if count > 0 {
			fmt.Printf("Migration '%s' already applied, skipping...\n", m.name)
			continue
		}

		_, err := s.db.Exec(m.stmt)
		if err != nil {
			return err
		}

		_, err = s.db.Exec("INSERT INTO migrations (name) VALUES (?)", m.name)
		if err != nil {
			return err
		}

		fmt.Printf("Applied migration '%s'\n", m.name)
	}

	return nil
}

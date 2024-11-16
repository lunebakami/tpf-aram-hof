package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Service interface {
	Health() error
	Migrate() error
	Close() error

	CreatePlayer(Player) (sql.Result, error)
  GetPlayers() ([]Player, error)
  DeletePlayer(int) (sql.Result, error)
}

type service struct {
	db *sql.DB
}

type Player struct {
	ID          int       `json:"id"`
	Nickname    string    `json:"nickname"`
	Champion    string    `json:"champion"`
	Description string    `json:"description"`
	GameMode    string    `json:"game_mode"`
	Frag        string    `json:"frag"`
	Date        time.Time `json:"date"`
}

var (
	dbUrl      = os.Getenv("DATABASE_URL")
	authToken  = os.Getenv("DATABASE_AUTH_TOKEN")
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

	dbInstance = &service{
		db: db,
	}

  err = dbInstance.Migrate()
  if err != nil {
    log.Fatalf("Error migrating database: %e", err)
  }
	return dbInstance
}

func (s *service) CreatePlayer(p Player) (sql.Result, error) {
	nickname := p.Nickname
	champion := p.Champion
	description := p.Description
	gameMode := p.GameMode
	frag := p.Frag
	date := p.Date

	result, err := s.db.Exec(`
    INSERT INTO players (nickname, champion, description, game_mode, frag, date)
    VALUES (?, ?, ?, ?, ?, ?)
    `,
		nickname, champion, description, gameMode, frag, date)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetPlayers() ([]Player, error) {
  rows, err := s.db.Query(`
    SELECT id, nickname, champion, description, game_mode, frag, date
    FROM players
  `)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  players := []Player{}

  for rows.Next() {
    var p Player
    err := rows.Scan(&p.ID, &p.Nickname, &p.Champion, &p.Description, &p.GameMode, &p.Frag, &p.Date)
    if err != nil {
      log.Printf("Error scanning player: %e\n", err)
      continue
    }

    players = append(players, p)
  }

  return players, nil
}

func (s *service) DeletePlayer(id int) (sql.Result, error) {
  result, err := s.db.Exec(`
    DELETE FROM players
    WHERE id = ?
  `, id)
  if err != nil {
    return nil, err
  }

  return result, nil
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
		{
			name: "add description and game mode to players",
			stmt: `
        ALTER TABLE players
        ADD COLUMN description TEXT;

        ALTER TABLE players
        ADD COLUMN game_mode TEXT;
      `,
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

package main

import (
	"database/sql"
	"f1-results-rod/data"
	"log"
	"time"
)

const maxDBLifetime = 5 * time.Minute

func initPostgresDB(dsn string) (data.Repository, error) {
	log.Println("connecting to DB")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(maxDBLifetime)

	repo := data.NewPostgresDB(db)
	err = repo.CreateTables()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

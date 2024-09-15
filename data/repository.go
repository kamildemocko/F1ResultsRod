package data

import "database/sql"

type Repository interface {
	CreateTables() error
	Close() error
	CountTrack(string, int) (int, error)
	GetTrackID(string, int) (int64, error)
	GetTracksYear(int) ([]Track, error)
	InsertTrack(string, string, int) (int64, error)
	InsertResultPlaces([]ResultPlace, int64) error
}

type postgresRepository struct {
	DB *sql.DB
}

func NewPostgresDB(db *sql.DB) Repository {
	return &postgresRepository{
		DB: db,
	}
}

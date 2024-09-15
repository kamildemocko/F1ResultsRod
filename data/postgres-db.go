package data

import (
	"context"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const connTimeout = 10 * time.Second

// func NewPostgresDB(dsn string) (Repository, error) {
// 	db, err := sql.Open("pgx", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err = db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	db.SetMaxIdleConns(maxIdleDBConn)
// 	db.SetMaxOpenConns(maxOpenDbConn)
// 	db.SetConnMaxLifetime(maxDBLifetime)

// 	log.Println("connected to DB")

// 	pgDb := &PostgresDB{db: db}

// 	err = pgDb.CreateTables()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pgDb, nil
// }

func (p *postgresRepository) Close() error {
	return p.DB.Close()
}

func (p *postgresRepository) CreateTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	queryTracks := `
		CREATE TABLE IF NOT EXISTS f1scrap.tracks (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			link VARCHAR(255),
			year INTEGER,
			CONSTRAINT unique_name_year UNIQUE (name, year)
		);`

	queryResult := `
		CREATE TABLE IF NOT EXISTS f1scrap.results (
			id SERIAL PRIMARY KEY,
			position INTEGER,
			driver_no INTEGER,
			driver VARCHAR(255),
			car VARCHAR(255),
			laps INTEGER,
			time_or_retired VARCHAR(255),
			points INTEGER,
			track_id INTEGER,
			FOREIGN KEY (track_id) REFERENCES f1scrap.tracks(id),
			CONSTRAINT unique_driver_track_id UNIQUE (driver, track_id)
		);`

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryTracks)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, queryResult)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (p *postgresRepository) CountTrack(trackName string, year int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM f1scrap.tracks
	WHERE name=$1 AND year=$2;`

	row := p.DB.QueryRowContext(ctx, query, trackName, year)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (p *postgresRepository) GetTrackID(trackName string, year int) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	query := `SELECT id FROM f1scrap.tracks
	WHERE name=$1 AND year=$2;`

	row := p.DB.QueryRowContext(ctx, query, trackName, year)

	var tid int64
	err := row.Scan(&tid)
	if err != nil {
		return -1, err
	}

	return tid, nil
}

func (p *postgresRepository) GetTracksYear(year int) ([]Track, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	query := `SELECT id, name, link, year FROM f1scrap.tracks WHERE year=$1`

	rows, err := p.DB.QueryContext(ctx, query, year)
	if err != nil {
		return nil, err
	}

	var tracks []Track

	for rows.Next() {
		var track Track

		err = rows.Scan(&track.ID, &track.Name, &track.Link, &track.Year)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (p *postgresRepository) InsertTrack(trackName, link string, year int) (int64, error) {
	log.Printf("..inserting track %s\n", trackName)

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	query := `INSERT INTO f1scrap.tracks (name, link, year) VALUES ($1, $2, $3) 
	ON CONFLICT (name, year) DO NOTHING 
	RETURNING  id;`

	var tid int64
	err := p.DB.QueryRowContext(ctx, query, trackName, link, year).Scan(&tid)
	if err != nil {
		return 0, err
	}

	return tid, nil
}

func (p *postgresRepository) InsertResultPlaces(values []ResultPlace, trackId int64) error {
	log.Printf("inserting result place for track ID: %d\n", trackId)

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout)
	defer cancel()

	query := `INSERT INTO f1scrap.results (position, driver_no, driver, car, laps, time_or_retired, points, track_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
	ON CONFLICT (driver, track_id) DO UPDATE SET 
	position=$1, driver_no=$2, driver=$3, car=$4, laps=$5, time_or_retired=$6, points=$7, track_id=$8;`

	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}

	for _, rp := range values {
		_, err = tx.ExecContext(
			ctx, query,
			rp.Position,
			rp.DriverNo,
			rp.Driver,
			rp.Car,
			rp.Laps,
			rp.TimeOrRetired,
			rp.Points,
			trackId,
		)

		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

package data

import "time"

const (
	maxOpenDbConn = 25
	maxIdleDBConn = 25
	maxDBLifetime = 5 * time.Minute
	connTimeout   = 3 * time.Second
)

type DBIntf interface {
	CreateTables() error
	Close() error
	CountTrack(string, int) (int, error)
	GetTrackID(string, int) (int64, error)
	GetTracksYear(int) ([]Track, error)
	InsertTrack(string, string, int) (int64, error)
	InsertResultPlaces([]ResultPlace, int64) error
}

type Track struct {
	ID   int64
	Name string
	Link string
	Year int
}

type ResultPlace struct {
	Position      int    `json:"position"`
	DriverNo      int    `json:"driverNo"`
	Driver        string `json:"driver"`
	Car           string `json:"car"`
	Laps          int    `json:"laps"`
	TimeOrRetired string `json:"timeRetired"`
	Points        int    `json:"points"`
}

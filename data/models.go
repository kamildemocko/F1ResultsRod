package data

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

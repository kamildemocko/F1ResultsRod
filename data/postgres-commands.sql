CREATE TABLE IF NOT EXISTS f1scrap.tracks (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255),
	link VARCHAR(255),
	year INTEGER,
	CONSTRAINT unique_name_year UNIQUE (name, year)
);

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
);

SELECT * FROM f1scrap.tracks

SELECT id, name, link, year FROM f1scrap.tracks
WHERE year=2024

SELECT COUNT(*) FROM f1scrap.tracks
WHERE name='FORMULA 1 GRAND PRIX DE MONACO 2024' AND year=2024

INSERT INTO f1scrap.tracks (name, year) 
VALUES ('test', 2024)
ON CONFLICT (name, year) DO NOTHING;

SELECT * FROM f1scrap.results

INSERT INTO f1scrap.results (position, driver_no, driver, car, laps, time_or_retired, points, track_id) 
VALUES (1, 20, 'test', 'test', 30, 'test', 25, 1)
ON CONFLICT (driver, track_id) DO NOTHING;

SELECT * FROM f1scrap.results
WHERE track_id = (
	SELECT id FROM f1scrap.tracks
	WHERE name='FORMULA 1 GRAND PRIX DE MONACO 2024' AND year=2024
)
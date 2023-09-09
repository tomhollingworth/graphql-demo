package graph

import (
	"database/sql" // add this

	_ "github.com/lib/pq" // add this
)

var (
	connectionString = "postgresql://postgres:postgres@postgres/equipment?sslmode=disable"
	db               *sql.DB
)

func init() {
	// initialize the database connection
	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	// check that the tables exist, if not create them
	_, err = db.Query("SELECT * FROM equipment LIMIT 1;")
	if err != nil {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS equipment (id text PRIMARY KEY, name TEXT NOT NULL, description TEXT)"); err != nil {
			panic(err)
		}
	}
	_, err = db.Query("SELECT * FROM equipment_property LIMIT 1")
	if err != nil {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS equipment_property (id text PRIMARY KEY, description TEXT NOT NULL, equipment_id text NOT NULL, FOREIGN KEY (equipment_id) REFERENCES equipment (id))"); err != nil {
			panic(err)
		}
	}
}

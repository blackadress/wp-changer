package wals

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func CreateWalTable(db *sql.DB) {
	walls_table := `
		CREATE TABLE IF NOT EXISTS wals (
			id integer NOT NULL PRIMARY KEY autoincrement,
			file text
		);
	`
	query, err := db.Prepare(walls_table)
	if err != nil {
		log.Fatal(err)
	}
	query.Exec()
}

func AreWallsInDB(db *sql.DB) bool {
	query := `SELECT * FROM wals where id = 1;`
	wals, err := db.Query(query)
	if err != nil {
		fmt.Println("Error haciendo select de wals en db: ", err)
		return false
	}

	defer wals.Close()
	if wals.Next() {
		// fmt.Println("Ya hay wallpapers en la BD")
		return true
	}

	return false
}

func GetWalls(db *sql.DB) []string {
	query := `SELECT * FROM wals ORDER BY id;`
	var results []string
	wals, err := db.Query(query)
	if err != nil {
		log.Println("No se pudo obtener registros de wals: ", err)
	}
	defer wals.Close()

	for wals.Next() {
		var id int
		var file string

		wals.Scan(&id, &file)
		results = append(results, file)
		// log.Printf("Id: %d, filepath: %s\n", id, file)
	}

	return results
}

func GetWallById(db *sql.DB, id int) string {
	var file string
	query := `
		SELECT file FROM wals
		WHERE id = ?;
	`
	err := db.QueryRow(query, id).Scan(&file)
	if err != nil {
		log.Printf("GetWallById: error obteniendo wall con id %d de DB: %s\n", id, err)
	}

	return file
}

func InsertWals(db *sql.DB, path string) {
	query := `
		INSERT INTO wals (file)
			VALUES (?);
	`
	records, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparando query de insercion a la tabla wals, error: ", err)
	}
	_, err = records.Exec(path)
	if err != nil {
		log.Println("Error insertando valores a la tabla wals, error: ", err)
	}
}

func DeleteWal(db *sql.DB, id int) error {
	query := `
		DELETE FROM wals WHERE ID = ?;
	`
	_, err := db.Exec(query, id)

	return err
}

func DeleteAllWals(db *sql.DB) error {
	query := `
		DELETE FROM wals;
	`

	_, err := db.Exec(query)

	return err
}

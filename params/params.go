package params

import (
	"database/sql"
	"fmt"
	"log"
)

type Params struct {
	ID       int
	Random   bool
	Curr_wal int
	Last_wal int
	Previous bool
}

func CreateParamTable(db *sql.DB) {
	params_table := `
		CREATE TABLE IF NOT EXISTS params (
			id integer NOT NULL PRIMARY KEY autoincrement,
			random boolean,
			curr_wal integer,
			last_wal integer,
			FOREIGN KEY (curr_wal) REFERENCES wals (id)
		);
	`
	query, err := db.Prepare(params_table)
	if err != nil {
		log.Fatal("Error creando tabla de parametros ", err)
	}

	query.Exec()
}

func InsertParam(db *sql.DB, p Params) Params {
	query := `
		INSERT INTO params (id, random, curr_wal, last_wal)
			VALUES (?, ?, ?, ?)
		RETURNING id;
	`
	db.QueryRow(query, 1, p.Random, p.Curr_wal, p.Last_wal).Scan(&p.ID)

	return p
}

func UpdateParam(db *sql.DB, p Params) error {
	query := `
		UPDATE
			params
		SET
			random = ?,
			curr_wal = ?,
			last_wal = ?
		WHERE
			id = 1;
	`
	_, err := db.Exec(query, p.Random, p.Curr_wal, p.Last_wal)

	return err
}

func AreParamsInDB(db *sql.DB) bool {
	query := `SELECT * FROM params WHERE id = 1;`
	params, err := db.Query(query)
	if err != nil {
		fmt.Println("Error haciendo select de params en db: ", err)
		return false
	}
	defer params.Close()
	if params.Next() {
		// fmt.Println("Ya hay params en la BD")
		return true
	}

	return false
}

func GetParams(db *sql.DB) Params {
	p := Params{}

	query := `SELECT * FROM params WHERE id = 1;`
	db.QueryRow(query).Scan(&p.ID, &p.Random, &p.Curr_wal, &p.Last_wal)

	return p
}

func DeleteParam(db *sql.DB) error {
	query := `DELETE FROM params WHERE id = 1`
	_, err := db.Exec(query)

	return err
}

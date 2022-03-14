package error_files

import (
	"database/sql"
	"log"
)

type ErrorFiles struct {
	ID    int
	File  string
	Fecha string
}

func CreateErrorFilesTable(db *sql.DB) {
	errors_table := `
		CREATE TABLE IF NOT EXISTS error_files (
			id integer NOT NULL PRIMARY KEY autoincrement,
			file text,
			fecha text
		);
	`
	query, err := db.Prepare(errors_table)
	if err != nil {
		log.Fatal("Error creando tabla de errores", err)
	}

	query.Exec()
}

func InsertErrorFile(db *sql.DB, e ErrorFiles) ErrorFiles {
	query := `
		INSERT INTO error_files (file, fecha)
			VALUES (?, ?)
		RETURNING
			id;
	`
	db.QueryRow(query, e.File, e.Fecha).Scan(&e.ID)

	return e
}

func GetErrorFiles(db *sql.DB) []ErrorFiles {
	query := `SELECT * FROM error_files ORDER BY id;`
	var result []ErrorFiles
	efs, err := db.Query(query)
	if err != nil {
		log.Println("No se pudo obtener registros de error_files: ", err)
	}
	defer efs.Close()

	for efs.Next() {
		var ef ErrorFiles

		efs.Scan(&ef.ID, &ef.File, &ef.Fecha)
		result = append(result, ef)
	}

	return result
}

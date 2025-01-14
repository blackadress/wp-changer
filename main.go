package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"wallchanger/pkg/error_files"
	"wallchanger/pkg/params"
	"wallchanger/pkg/wals"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	home := os.Getenv("HOME")
	database_name := home + "/.config/wallpaper-go/database.db"
	create_database(database_name)

	// var random bool
	var prev bool
	// var next bool
	var directory string

	// flag.BoolVar(&random, "r", false, "If set then wallpapers will change randomly. Default is off.")
	flag.BoolVar(&prev, "p", false, "Will change wallpaper to the previous one in queue.")
	// flag.BoolVar(&next, "n", true, "Will change wallpaper to the next one in queue.")
	flag.StringVar(&directory, "d", "./example", "The directory that has the wallpapers")

	flag.Parse()

	// si el path es relativo, hacerlo completo
	var path string = directory
	if "." == directory[0:1] {
		current_directory, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		path = fmt.Sprintf("%s/%s", current_directory, directory[2:])
	}

	files, err := list_files(path)
	if err != nil {
		log.Println(err)
		return
	}

	db, err := sql.Open("sqlite3", database_name)
	if err != nil {
		log.Fatal("database connection error: ", err)
		return
	}

	initWals(db, files)
	initErrorFiles(db)
	p := params.Params{Last_wal: len(files), Previous: prev}
	p = initParams(db, p)

	setWallpaper(db, p)

}

func setWallpaper(db *sql.DB, p params.Params) {
	// if prev is set, then change wallpaper to previous one
	if p.Previous {
		p.Curr_wal -= 1
	} else {
		p.Curr_wal += 1
	}

	// si se llego al final del index, entonces dar la vuelta desde el inicio
	if p.Curr_wal-1 == p.Last_wal {
		p.Curr_wal = 1
	}

	// si es el primer wallpaper y se da la opcion de prev, set el ultimo wallpaper
	if p.Curr_wal+1 == 1 && p.Previous {
		p.Curr_wal = p.Last_wal
	}

	// get current wallpaper by id
	w := wals.GetWallById(db, p.Curr_wal)

	out, err := exec.Command(
		"feh",
		"--bg-fill",
		w,
	).Output()

	if err != nil {
		// si hay error cambiar al siguiente wallpaper
		log.Printf("setWallpaper: error poniendo el wallpaper %d, error: %s out: %s\n", p.Curr_wal, err, out)

		// insertar el archivo que dio error a la tabla de errores
		now := time.Now().Format(time.RFC1123Z)
		file_name := wals.GetWallById(db, p.Curr_wal)
		e := error_files.ErrorFiles{File: file_name, Fecha: now}
		error_files.InsertErrorFile(db, e)

		p.Curr_wal += 1
		setWallpaper(db, p)
	}

	params.UpdateParam(db, p)
}

func initWals(db *sql.DB, files []string) {
	wals.CreateWalTable(db)
	if wals.AreWallsInDB(db) {
		return
	}

	fmt.Println("Insertando los wallpapers a la BD...")

	extensiones_validas := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"webp": true,
	}

	for _, v := range files {
		extension := v[strings.LastIndex(v, ".")+1:]
		if extensiones_validas[extension] {
			wals.InsertWals(db, v)
		}
	}
}

func updateWallpapersFolder(db *sql.DB, files []string) {
	err := wals.DeleteAllWals(db)
	if err != nil {
		fmt.Println("Error eliminando la tabla de wallpapers: ", err)
	}

	for _, v := range files {
		wals.InsertWals(db, v)
		// fmt.Println(v)
	}
}

func initParams(db *sql.DB, p params.Params) params.Params {
	params.CreateParamTable(db)
	if params.AreParamsInDB(db) {
		db_param := params.GetParams(db)
		p.Curr_wal = db_param.Curr_wal
		err := params.UpdateParam(db, p)
		if err != nil {
			log.Fatalf("mergas, algo paso intentando updatear los params %s\n", err)
		}
		return p
	}
	fmt.Println("Insertando los params a la BD")
	p.Curr_wal = 1

	p = params.InsertParam(db, p)
	return p
}

func initErrorFiles(db *sql.DB) {
	error_files.CreateErrorFilesTable(db)
}

func create_database(database_name string) {
	_, err := os.Open(database_name)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("Database file does not exist, creating one...")
		file, err := os.Create(database_name)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	} else if err != nil {
		log.Printf("Something unexpected happened %s/n", err)
	}
}

func list_files(path string) ([]string, error) {
	var answer []string
	f, err := os.Open(path)
	if err != nil {
		return answer, err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return answer, err
	}

	for _, v := range files {
		if !v.IsDir() {
			full_path := filepath.Join(path, v.Name())
			answer = append(answer, full_path)
		} else {
			nested_path := filepath.Join(path, v.Name())
			nested_files, err := list_files(nested_path)
			if err != nil {
				log.Println(err)
			}
			answer = append(answer, nested_files...)
		}
	}
	return answer, nil
}

func scramble(wallpapers []string) {
	// shuffle
	rand.Shuffle(len(wallpapers), func(i, j int) {
		wallpapers[i], wallpapers[j] = wallpapers[j], wallpapers[i]
	})
}

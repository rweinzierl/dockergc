package lib

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const typContainer = "container"
const typImage = "image"

func exitOnDbError(err error) {
	if err != nil {
		log.Printf("Error accessing database %s: %s", dbPath(), err)
		os.Exit(1)
	}
}

func readAll(typ string) *[]string {
	rows, err := getDB().Query("SELECT name from pin where typ = ?", typ)
	exitOnDbError(err)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}
	return &names
}

func add(typ string, name string) {
	stmt, err := getDB().Prepare("insert or ignore into pin(typ, name) values(?, ?)")
	exitOnDbError(err)
	stmt.Exec(typ, name)
	defer stmt.Close()
}

func remove(typ string, name string) {
	stmt, err := getDB().Prepare("delete from pin where typ = ? and name = ?")
	exitOnDbError(err)
	stmt.Exec(typ, name)
	defer stmt.Close()
}

func removeAll() {
	stmt, err := getDB().Prepare("delete from pin")
	exitOnDbError(err)
	stmt.Exec()
	defer stmt.Close()
}

const VarNameDbPath = "DOCKERGC_DB"

func dbPathDefault() string {
	usr, err := user.Current()
	exitOnDbError(err)
	return filepath.Join(usr.HomeDir, ".dockercg", "db.sqlite")
}

func dbPath() string {
	if value, ok := os.LookupEnv(VarNameDbPath); ok && value != "" {
		return value
	} else {
		return dbPathDefault()
	}
}

func connect() *sql.DB {
	path := dbPath()
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	database, err := sql.Open("sqlite3", path)
	exitOnDbError(err)
	statement, err := database.Prepare("create table if not exists pin (typ text, name text, primary key (typ, name))")
	exitOnDbError(err)
	statement.Exec()
	return database
}

var theDb *sql.DB

func getDB() *sql.DB {
	if theDb == nil {
		theDb = connect()
	}
	return theDb
}

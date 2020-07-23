package main

import (
	"database/sql"
	"os"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func ReadAll(typ string) []string {
	rows, _ := Db.Query("SELECT name from pin where typ = ?", typ)
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}
	return names
}

func Add(typ string, name string) {
	stmt, _ := Db.Prepare("insert or ignore into pin(typ, name) values(?, ?)")
	stmt.Exec(typ, name)
	defer stmt.Close()
}

func Remove(typ string, name string) {
	stmt, _ := Db.Prepare("delete from pin where typ = ? and name = ?")
	stmt.Exec(typ, name)
	defer stmt.Close()
}

func RemoveAll() {
	stmt, _ := Db.Prepare("delete from pin")
	stmt.Exec()
	defer stmt.Close()
}

const VarNameDbPath = "DOCKERGC_DB"

func DbPathDefault() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".dockercg", "db.sqlite")
}

func DbPath() string {
	if value, ok := os.LookupEnv(VarNameDbPath); ok {
		return value
	} else {
		return DbPathDefault()
	}
}

func db() *sql.DB {
	path := DbPath()
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	database, _ := sql.Open("sqlite3", path)
	statement, _ := database.Prepare("create table if not exists pin (typ text, name text, primary key (typ, name))")
	statement.Exec()
	return database
}

var Db = db()

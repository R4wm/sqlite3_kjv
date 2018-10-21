package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("Starting sqlite3 db creation. ")
	database, _ := sql.Open("sqlite3", "data/kjv.sqlite3.db")
	defer database.Close()

	statement, _ := database.Prepare("create table if not exists kjv(book string not null, chapter int, verse int, text string)")
	statement.Exec()

	fmt.Println("Database done.")
}

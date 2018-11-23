package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Verse the complete verse context
type Verse struct {
	IsNumberedBook bool
	Book           string
	Chapter        int
	Verse          int
	Text           string
}

//ParseChapterVerse extract chapter and verse from x:x format
func ParseChapterVerse(colonJoined string) (int, int) {
	fmt.Printf("colonJoined: %v\n", colonJoined)

	splitChapterVerse := strings.Split(colonJoined, ":")

	chapter, err := strconv.Atoi(splitChapterVerse[0])
	if err != nil {
		panic(err)
	}

	verseNum, err := strconv.Atoi(splitChapterVerse[1])
	if err != nil {
		panic(err)
	}

	return chapter, verseNum

}

// IsNumberedBook determines if this is numbered book like 1John or 2Timothy.
func IsNumberedBook(firstPart string) bool {
	// firstPart is the very first element in the parsed string.
	if _, err := strconv.Atoi(firstPart); err == nil {
		return true
	}
	return false
}

//PrepareDB Inserts verse context into database. Note old db WILL be deleted.
func PrepareDB(verse <-chan Verse, dbPath string) {

	// Delete old existing database
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Could not remove old database: %s", dbPath)
		}
	}

	//Create new database
	database, _ := sql.Open("sqlite3", dbPath) // "data/kjv.sqlite3.db")
	defer database.Close()

	//Prep new database
	statement, _ := database.Prepare("create table if not exists kjv(book string not null, chapter int, verse int, text string)")
	statement.Exec()

	sqlInsertStr := `INSERT OR REPLACE INTO kjv(book, chapter, verse, text) values(?, ?, ?, ?)`
	stmt, err := database.Prepare(sqlInsertStr)
	if err != nil {
		panic(err)
	}

	//Populate, put into database as they come
	defer stmt.Close()
	for v := range verse {
		stmt.Exec(v.Book, v.Chapter, v.Verse, v.Text)
	}
}

//CreateKJVDB pulls down KJV raw text file, parses and creates database
func CreateKJVDB() {
	fmt.Println("Starting sqlite3 db creation. ")

	url := "https://raw.githubusercontent.com/R4wm/bible/master/data/bible.txt"
	dbInsert := make(chan Verse)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	go PrepareDB(dbInsert, "/tmp/kjv.db")

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		verse := Verse{}
		brokenString := strings.Fields(scanner.Text())
		fmt.Println("broken: ", brokenString)

		if brokenString[0] == "Song" {
			//This is Song of Solomon book, special case where book name has multiple words
			verse.Book = fmt.Sprintf("%s %s %s", brokenString[0], brokenString[1], brokenString[2])
			verse.Chapter, verse.Verse = ParseChapterVerse(brokenString[3])
			verse.Text = strings.Join(brokenString[4:], " ")
		} else if IsNumberedBook(brokenString[0]) {
			verse.Book = brokenString[0] + brokenString[1]
			verse.Chapter, verse.Verse = ParseChapterVerse(brokenString[2])
			verse.Text = strings.Join(brokenString[3:], " ")
		} else {
			verse.Book = brokenString[0]
			verse.Chapter, verse.Verse = ParseChapterVerse(brokenString[1])
			verse.Text = strings.Join(brokenString[2:], " ")
		}

		fmt.Printf("verse: %v\n", verse)

		dbInsert <- verse
	}

	close(dbInsert)
}

func main() {
	CreateKJVDB()
}

package sqlite3_kjv

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	sqlInsertStr = `INSERT OR REPLACE INTO kjv(book, chapter, verse, text, ordinal_verse, ordinal_book, testament) values(?, ?, ?, ?, ?, ?, ?)`
)

var url = "https://raw.githubusercontent.com/R4wm/sqlite3_kjv/KJV_PCE/data/TEXT-KJV-PCE-127-TAB.txt"

// Verse the complete verse context
type Verse struct {
	isNumberedBook bool
	Book           string `json:"book"`
	Chapter        int    `json:"chapter"`
	Verse          int    `json:"verse"`
	Text           string `json:"text"`
	Testament      string `json:"canonical_testament"`
	OrdinalVerse   int    `json:"ordinal_verse"`
	OrdinalBook    int    `json:"ordinal_book"`
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

// isNumberedBook determines if this is numbered book like 1John or 2Timothy.
func isNumberedBook(firstPart string) bool {
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
	statement, _ := database.Prepare("create table if not exists kjv(book string not null, chapter int, verse int, text string, ordinal_verse int, ordinal_book int, testament string)")
	statement.Exec()

	stmt, err := database.Prepare(sqlInsertStr)
	if err != nil {
		panic(err)
	}

	//Populate, put into database as they come
	defer stmt.Close()
	for v := range verse {
		stmt.Exec(v.Book, v.Chapter, v.Verse, v.Text, v.OrdinalVerse, v.OrdinalBook, v.Testament)
	}
}

//CreateKJVDB pulls down KJV raw text file, parses and creates database
func CreateKJVDB(dbpath string) (string, error) {
	fmt.Println("Starting sqlite3 db creation. ")

	dbInsert := make(chan Verse)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s path does not exist.\n", url)
	}
	defer resp.Body.Close()

	go PrepareDB(dbInsert, dbpath)
	defer close(dbInsert)

	scanner := bufio.NewScanner(resp.Body)

	verseCount := 0
	bookTestament := "OLD"

	for scanner.Scan() {
		verseCount++
		verse := Verse{}
		brokenString := strings.Fields(scanner.Text())
		fmt.Println("Broken up string: ", brokenString)

		if brokenString[2] == "Matthew" {
			bookTestament = "NEW"
		}

		// Special use case where title of book is multiple strings
		if brokenString[1] == "Song" {
			//This is Song of Solomon book, special case where book name has multiple words
			verse.Book = fmt.Sprintf("%s %s %s",
				strings.ToUpper(brokenString[2]), // SONG
				strings.ToUpper(brokenString[3]), // OF
				strings.ToUpper(brokenString[4])) // SOLOMON
			verse.Chapter, _ = strconv.Atoi(brokenString[5])
			verse.Verse, _ = strconv.Atoi(brokenString[6])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[7:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		} else if isNumberedBook(brokenString[2]) {
			verse.Book = strings.ToUpper(brokenString[2] + brokenString[3])
			verse.Chapter, _ = strconv.Atoi(brokenString[4])
			verse.Verse, _ = strconv.Atoi(brokenString[5])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[6:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		} else {
			verse.Book = strings.ToUpper(brokenString[2])
			verse.Chapter, _ = strconv.Atoi(brokenString[3])
			verse.Verse, _ = strconv.Atoi(brokenString[4])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[5:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		}

		verse.OrdinalBook, _ = strconv.Atoi(brokenString[0])
		verse.OrdinalVerse = verseCount
		verse.Testament = bookTestament

		fmt.Printf("verse: %v\n", verse)
		dbInsert <- verse

	}

	return dbpath, nil
}

//CreateKJVJson pulls down KJV raw text file, parses and creates database
func CreateKJVJson(jsonPath string) (string, error) {
	fmt.Println("Starting json creation. ")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s path does not exist.\n", url)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	verseCount := 0
	bookTestament := "OLD"

	jsonText := []Verse{}
	for scanner.Scan() {
		verseCount++
		verse := Verse{}
		brokenString := strings.Fields(scanner.Text())
		fmt.Println("Broken up string: ", brokenString)

		if brokenString[2] == "Matthew" {
			bookTestament = "NEW"
		}

		// Special use case where title of book is multiple strings
		if brokenString[1] == "Song" {
			//This is Song of Solomon book, special case where book name has multiple words
			verse.Book = fmt.Sprintf("%s %s %s",
				strings.ToUpper(brokenString[2]), // SONG
				strings.ToUpper(brokenString[3]), // OF
				strings.ToUpper(brokenString[4])) // SOLOMON
			verse.Chapter, _ = strconv.Atoi(brokenString[5])
			verse.Verse, _ = strconv.Atoi(brokenString[6])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[7:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		} else if isNumberedBook(brokenString[2]) {
			verse.Book = strings.ToUpper(brokenString[2] + brokenString[3])
			verse.Chapter, _ = strconv.Atoi(brokenString[4])
			verse.Verse, _ = strconv.Atoi(brokenString[5])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[6:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		} else {
			verse.Book = strings.ToUpper(brokenString[2])
			verse.Chapter, _ = strconv.Atoi(brokenString[3])
			verse.Verse, _ = strconv.Atoi(brokenString[4])
			verse.Text = strings.ReplaceAll(strings.Join(brokenString[5:], " "), "<<", "-|")
			verse.Text = strings.ReplaceAll(verse.Text, ">>", "|-")
		}

		verse.OrdinalBook, _ = strconv.Atoi(brokenString[0])
		verse.OrdinalVerse = verseCount
		verse.Testament = bookTestament
		jsonText = append(jsonText, verse)
		fmt.Printf("verse: %v\n", verse)
	}

	jsonByte, err := json.Marshal(jsonText)
	if err != nil {
		return "", err
	}

	fmt.Println(jsonByte)

	f, err := os.Create(jsonPath)
	defer f.Close()
	n2, err := f.Write(jsonByte)
	f.Sync()
	fmt.Println("jsonPath: ", jsonPath)
	if err != nil {
		fmt.Println("failed to f.Write: ", err)
		return "", err
	}
	fmt.Println("wrote: ", n2)
	return jsonPath, nil
}

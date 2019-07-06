package sqlite3_kjv

import (
	"fmt"
	"testing"
)

func TestGetChapter(t *testing.T) {
	var verses []KJVVerse
	book := "Genesis"
	chapter := 1
	verses = append(verses,
		KJVVerse{Verse: 1,
			Text: "In the beginning God created the heaven and the earth."})

	verses = append(verses,
		KJVVerse{Verse: 2,
			Text: "And the earth was without form, and void; and darkness was upon the face of the deep. And the Spirit of God moved upon the face of the waters."})

	result := GetChapter(book, chapter, verses)
	if fmt.Sprintf("%T", result) != "*sqlite3_kjv.KJVBook" {
		t.Errorf("Expected *sqlite3_kjv.KJVBook, got %T\n", result)
	}

	if result.Book != book {
		t.Errorf("Expected %s, got %v\n", book, result.Book)
	}
}

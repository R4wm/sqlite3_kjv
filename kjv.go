package kjvapi

//KJVVerse simple container for verses
type KJVVerse struct {
	Verse int    `json:"verse"`
	Text  string `json:"text"`
}

//KJVChapter simple container for chapters
type KJVChapter struct {
	Chapter int `json:"chapter"`
	Verses  []KJVVerse
}

//KJVBook simpla container for book from Bible
type KJVBook struct {
	Book     string `json:"book"`
	Chapters []KJVChapter
}

// GetChapter compose a bible chapter with verses
func GetChapter(book string, chapter int, verses []KJVVerse) *KJVBook {
	return &KJVBook{
		Book: book,
		Chapters: []KJVChapter{
			KJVChapter{
				Chapter: chapter,
				Verses:  verses}},
	}
}

// func main() {
// 	var verses []KJVVerse

// 	verse := KJVVerse{1, "In the beginning God created the heaven and the earth."}
// 	verses = append(verses, verse)
// 	verse2 := KJVVerse{2, "And the earth was without form, and void; and darkness was upon the face of the deep. And the Spirit of God moved upon the face of the waters."}
// 	verses = append(verses, verse2)

// 	query := GetChapter("Genesis", 1, verses)
// 	jsonQuery, _ := json.MarshalIndent(query, "", "  ")
// 	fmt.Printf("%s\n", jsonQuery)
// }

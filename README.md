# sqlite3_kjv

# Basic usage
```bash
ᚱ@elastic_kjv $ go get -v -u github.com/r4wm/elastic_kjv
github.com/r4wm/elastic_kjv (download)
github.com/mattn/go-sqlite3 (download)
ᚱ@elastic_kjv $ 
ᚱ@elastic_kjv $ 
ᚱ@elastic_kjv $ go run cmd/main.go -dbPath /tmp/kjv.db 2&> /dev/null 
ᚱ@elastic_kjv $ 
ᚱ@elastic_kjv $ 
ᚱ@elastic_kjv $ sqlite3 /tmp/kjv.db ".schema kjv"
CREATE TABLE kjv(book string not null, chapter int, verse int, text string, ordinal_verse int, ordinal_book int, testament string);
ᚱ@elastic_kjv $ 
ᚱ@elastic_kjv $ sqlite3 /tmp/kjv.db "select text from kjv where book=\"GENESIS\" and chapter=1 and verse=1"
In the beginning God created the heaven and the earth.
ᚱ@elastic_kjv $ 
```

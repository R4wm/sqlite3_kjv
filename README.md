# sqlite3_kjv

# Basic usage
```bash
ᚱ@cmd $ 
ᚱ@cmd $ go get -v -u github.com/r4wm/sqlite3_kjv
github.com/r4wm/sqlite3_kjv (download)
github.com/mattn/go-sqlite3 (download)
ᚱ@cmd $ 
ᚱ@cmd $ 
ᚱ@cmd $ 
ᚱ@cmd $ go run create_sqlite3_kjv.go -dbPath /tmp/kjv.db 2&> /dev/null 
ᚱ@cmd $ 
ᚱ@cmd $ 
ᚱ@cmd $ 
ᚱ@cmd $ sqlite3 /tmp/kjv.db "select text from kjv where book=\"GENESIS\" and chapter=1 and verse=1"
In the beginning God created the heaven and the earth.
ᚱ@cmd $ 
ᚱ@cmd $ 
ᚱ@cmd $ 

```

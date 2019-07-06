package main

import (
	"flag"
	"fmt"

	"github.com/r4wm/sqlite3_kjv"
)

// main: Create the kjv database at desired path for mintz5/deploy.go
func main() {
	var dbPath = flag.String("dbPath", "/tmp/kjv.db", "Path where DB should be created.")
	flag.Parse()

	fmt.Printf("Creating kjv db to %s\n", *dbPath)
	_, err := sqlite3_kjv.CreateKJVDB(*dbPath)

	if err != nil {
		panic(err)
	}

}

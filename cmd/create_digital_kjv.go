package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/r4wm/sqlite3_kjv"
)

// main: Create the kjv database at desired path for mintz5/deploy.go
func main() {
	var dbPath string
	var jsonPath string
	flag.StringVar(&dbPath, "dbPath", "", "DB Path")
	flag.StringVar(&jsonPath, "jsonPath", "", "Write Json path")
	flag.Parse()

	fmt.Println("dbPath: ", dbPath)
	fmt.Println("jsonPath: ", jsonPath)
	if dbPath != "" && jsonPath != "" {
		log.Fatal("args are mutually exclusive")
	}

	if dbPath != "" {
		fmt.Printf("Creating kjv db to %s\n", dbPath)
		_, err := sqlite3_kjv.CreateKJVDB(dbPath)

		if err != nil {
			panic(err)
		}
	} else if jsonPath != "" {
		fmt.Printf("Creating kjvjson to %s\n", jsonPath)
		_, err := sqlite3_kjv.CreateKJVJson(jsonPath)
		if err != nil {
			panic(err)
		}
	} else {
		log.Fatal("Nothing to do")
	}

}

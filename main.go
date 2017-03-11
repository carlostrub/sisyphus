package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	// Processed is a map of e-mail IDs and the value set to true if Junk
	Processed map[string]bool
)

func main() {

	// Get the Maildir to be handled and the DB
	wd, err := os.Getwd()
	if err != nil {
		return log.Fatalf("get working directory: %s", err)
	}
	maildir := flag.String("maildir", wd+"/Maildir", "Path to the Maildir")
	database := flag.String("database", wd+"/sisyphus.db", "Path to the sisyphus database")
	flag.Parse()

	// Load the Maildir
	mails, err := Index(*maildir)
	if err != nil {
		return log.Fatalf("load Maildir content: %s", err)
	}

	fmt.Println(mails)

	// Open the database
	db, err := openDB("sisyphus.db", 0600, nil)
	if err != nil {
		return log.Fatalf("open database: %s", err)
	}
	defer db.Close()

}

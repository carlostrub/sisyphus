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
	// Get the Maildir to be handled
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	maildir := flag.String("d", wd+"/Maildir", "Path of the Maildir to be handled")
	flag.Parse()

	// Load the Maildir content
	mails, err := Index(*maildir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(mails)

	// Create a classifier
	//classifier := bayesian.NewClassifier(Good, Junk)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jbrukh/bayesian"
)

const (
	// good is the class of good mails that are not supposed to be Spam
	good bayesian.Class = "Good"
	// junk is the class of Spam mails
	junk bayesian.Class = "Junk"
)

var (
	// Processed is a map of e-mail IDs and the value set to true if Junk
	Processed map[string]bool
)

// Classifiers contains the classifiers for mail subjects and bodies
type Classifiers struct {
	Subject, Body *bayesian.Classifier
}

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

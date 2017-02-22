package main

import (
	"bufio"
	"log"
	"strings"

	"github.com/jbrukh/bayesian"
	"github.com/luksen/maildir"
)

const (
	// good is the class of good mails that are not supposed to be Spam
	good bayesian.Class = "Good"
	// junk is the class of Spam mails
	junk bayesian.Class = "Junk"
)

var (
	// Maildir holds a set of mail directories to handle.
	Maildir = "/usr/home/cs/Maildir.TEST"

	// processed is a map of e-mail IDs and true if processed already.
	processed map[string]bool
)

// Mails contains the keys of all mails in the Junk.cur and cur directories.
type Mails struct {
	Junk, Good []string
}

// Classifiers contains the classifiers for mail subjects and bodies
type Classifiers struct {
	Subject, Body *bayesian.Classifier
}

// LoadMails loads all mail keys from the Maildir directory for processing.
func LoadMails() (m Mails, err error) {

	m.Junk, err = maildir.Dir(Maildir + "/.Junk").Keys()
	if err != nil {
		return m, err
	}

	m.Good, err = maildir.Dir(Maildir).Keys()
	if err != nil {
		return m, err
	}

	return m, nil
}

// Learn initially classifies all mails and returns the respective classifiers.
func (m Mails) Learn() (c Classifiers, err error) {
	return
}

func cleanText(t string) (c string, err error) {
	return
}

// getContent reads mails' subjects and bodies and returns the respective
// slices of strings
func getContent(keys []string) (s, b []string, err error) {
	for _, k := range keys {

		message, err := maildir.Dir(Maildir).Message(k)
		if err != nil {
			return s, b, err
		}

		// get Subject
		subject := message.Header.Get("Subject")
		s = append(s, strings.Split(subject, " ")...)

		// get Body
		bScanner := bufio.NewScanner(message.Body)
		for bScanner.Scan() {
			b = append(b, strings.Split(bScanner.Text(), " ")...)
		}
	}

	return s, b, nil
}

func main() {

	_, err := LoadMails()
	if err != nil {
		log.Fatal(err)
	}

	// Create a classifier
	//classifier := bayesian.NewClassifier(Good, Junk)
}

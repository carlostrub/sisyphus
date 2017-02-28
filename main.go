package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"mime/quotedprintable"
	"os"
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
	// Processed is a map of e-mail IDs and the value set to true if Junk
	Processed map[string]bool
)

// Mail includes the key of a mail in Maildir
type Mail struct {
	Key           string
	Subject, Body *string
	Junk          bool
}

// Classifiers contains the classifiers for mail subjects and bodies
type Classifiers struct {
	Subject, Body *bayesian.Classifier
}

// Index loads all mail keys from the Maildir directory for processing.
func Index(d string) (m []*Mail, err error) {

	g, err := maildir.Dir(d).Keys()
	if err != nil {
		return m, err
	}
	for _, val := range g {
		var new Mail
		new.Key = val
		m = append(m, &new)
	}

	j, err := maildir.Dir(d + "/.Junk").Keys()
	if err != nil {
		return m, err
	}
	for _, val := range j {
		var new Mail
		new.Key = val
		new.Junk = true
		m = append(m, &new)
	}

	return m, nil
}

// Learn initially classifies all mails and returns the respective classifiers.
func (m *Mail) Learn() (c Classifiers, err error) {
	return
}

func cleanString(i string) (s string, err error) {
	s = strings.ToLower(i)
	s = strings.Replace(s, "!", " ", -1)
	s = strings.Replace(s, "#", " ", -1)
	s = strings.Replace(s, "$", " ", -1)
	s = strings.Replace(s, "%", " ", -1)
	s = strings.Replace(s, "&", " ", -1)
	s = strings.Replace(s, "(", " ", -1)
	s = strings.Replace(s, ")", " ", -1)
	s = strings.Replace(s, "*", " ", -1)
	s = strings.Replace(s, "+", " ", -1)
	s = strings.Replace(s, ",", " ", -1)
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Replace(s, "/", " ", -1)
	s = strings.Replace(s, ":", " ", -1)
	s = strings.Replace(s, ";", " ", -1)
	s = strings.Replace(s, "<", " ", -1)
	s = strings.Replace(s, "=", " ", -1)
	s = strings.Replace(s, ">", " ", -1)
	s = strings.Replace(s, "@", " ", -1)
	s = strings.Replace(s, "[", " ", -1)
	s = strings.Replace(s, "\"", " ", -1)
	s = strings.Replace(s, "\\", " ", -1)
	s = strings.Replace(s, "\n", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	s = strings.Replace(s, "]", " ", -1)
	s = strings.Replace(s, "^", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Replace(s, "{", " ", -1)
	s = strings.Replace(s, "|", " ", -1)
	s = strings.Replace(s, "}", " ", -1)

	return s, nil
}

// Clean prepares the mail's subject and body for training
func (m *Mail) Clean() error {
	if m.Subject != nil {
		s, err := cleanString(*m.Subject)
		if err != nil {
			return err
		}
		m.Subject = &s
	}

	if m.Body != nil {
		b, err := cleanString(*m.Body)
		if err != nil {
			return err
		}
		m.Body = &b
	}
	return nil
}

// Load reads a mail's subject and body
func (m *Mail) Load(d string) error {

	message, err := maildir.Dir(d).Message(m.Key)
	if err != nil {
		return err
	}

	// get Subject
	if m.Subject != nil {
		return errors.New("there is already a subject")
	}
	subject := message.Header.Get("Subject")
	m.Subject = &subject

	// get Body
	bQ := quotedprintable.NewReader(message.Body)
	var b []string
	bScanner := bufio.NewScanner(bQ)
	for bScanner.Scan() {
		raw := bScanner.Text()
		b = append(b, raw)
	}

	body := strings.Join(b, " ")
	if m.Body != nil {
		return errors.New("there is already a body")
	}
	m.Body = &body

	return nil
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

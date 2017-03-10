package main

import (
	"github.com/jbrukh/bayesian"
)

const (
	// good is the class of good mails that are not supposed to be Spam
	good bayesian.Class = "Good"
	// junk is the class of Spam mails
	junk bayesian.Class = "Junk"
)

// Classifiers contains the classifiers for mail subjects and bodies
type Classifiers struct {
	Subject, Body *bayesian.Classifier
}

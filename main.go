package main

import (
	"fmt"
	"log"

	"github.com/luksen/maildir"
)

var (
	// Maildirs holds a set of mail directories to handle.
	Maildirs []string
)

func main() {
	Maildirs = []string{"/usr/home/cs/Maildir.TEST"}

	var err error
	var Bad, Good []string

	for _, dir := range Maildirs {
		var keysBad, keysGood []string
		keysBad, err = maildir.Dir(dir + "/.Junk").Keys()
		if err != nil {
			log.Fatal(err)
		}

		Bad = append(Bad, keysBad...)

		keysGood, err = maildir.Dir(dir).Keys()
		if err != nil {
			log.Fatal(err)
		}

		Good = append(Good, keysGood...)
	}

	fmt.Println("Bad guys:")
	fmt.Println(Bad)
	fmt.Println("Good guys:")
	fmt.Println(Good)
}

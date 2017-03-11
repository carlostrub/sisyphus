package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	// Processed is a map of e-mail IDs and the value set to true if Junk
	Processed map[string]bool
)

func main() {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		return log.Fatalf("get working directory: %s", err)
	}

	var maildir, database *string

	// Define App
	app := cli.NewApp()
	app.Name = "sisyphus"
	app.HelpName = "Intelligent Junk and Spam Mail Handler"
	app.Version = "0.0.0"
	app.Authors = []Author{
		Author{
			Name:  "Carlo Strub",
			Email: "cs@carlostrub.ch",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "maildir",
			Value:       wd + "/Maildir",
			Usage:       "Path to the Maildir directory",
			Destination: &maildir,
		},
		cli.StringFlag{
			Name:        "database",
			Value:       wd + "/sisyphus.db",
			Usage:       "Path to the sisyphus database",
			Destination: &database,
		},
	}

	app.Action = func(c *cli.Context) error {
		if database == nil {
			return errors.New("no database selected")
		}
		if maildir == nil {
			return errors.New("no maildir selected")
		}

		// Load the Maildir
		mails, err := Index(*maildir)
		if err != nil {
			return log.Fatalf("load Maildir content: %s", err)
		}

		fmt.Println(mails)

		// Open the database
		db, err := openDB(*database, 0600, nil)
		if err != nil {
			return log.Fatalf("open database: %s", err)
		}
		defer db.Close()

		return nil
	}

	app.Run(os.Args)
}

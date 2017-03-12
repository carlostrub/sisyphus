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
		log.Fatal(err)
	}

	var maildir, database *string

	// Define App
	app := cli.NewApp()
	app.Name = "Sisyphus"
	app.Usage = "Intelligent Junk and Spam Mail Handler"
	app.UsageText = `This application applies artificial intelligence to
	filter Junk mail in an unobtrusive way. Both, classification and
	learning operate directly on the Maildir of a user in a fully
	transparent mode, without any need for configuration or active
	operation.`
	app.HelpName = "Intelligent Junk and Spam Mail Handler"
	app.Version = "0.0.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Carlo Strub",
			Email: "cs@carlostrub.ch",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "maildir",
			Value:       wd + "/Maildir",
			Usage:       "Path to the Maildir directory",
			Destination: maildir,
		},
		cli.StringFlag{
			Name:        "database",
			Value:       wd + "/sisyphus.db",
			Usage:       "Path to the sisyphus database",
			Destination: database,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start sisyphus daemon",
			Action: func(c *cli.Context) error {
				if database == nil {
					return errors.New("no database selected")
				}
				if maildir == nil {
					return errors.New("no maildir selected")
				}

				// Load the Maildir
				mails, err := Index(*maildir)
				if err != nil {
					log.Fatalf("load Maildir content: %s", err)
				}

				fmt.Println(mails)

				// Open the database
				db, err := openDB(*database)
				if err != nil {
					log.Fatalf("open database: %s", err)
				}
				defer db.Close()

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"e"},
			Usage:   "stop sisyphus daemon",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart sisyphus daemon",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"i"},
			Usage:   "status of sisyphus",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"d"},
			Usage:   "short description of sisyphus",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	app.Run(os.Args)
}

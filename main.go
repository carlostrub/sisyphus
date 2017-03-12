package main

import (
	"errors"
	"fmt"
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
		panic(err)
	}

	var maildir []string

	// Define App
	app := cli.NewApp()
	app.Name = "Sisyphus"
	app.Usage = "Intelligent Junk and Spam Mail Handler"
	app.UsageText = `Sisyphus applies artificial intelligence to filter
	Junk mail in an unobtrusive way. Both, classification and learning
	operate directly on the Maildir of a user in a fully transparent mode,
	without any need for configuration or active operation.`
	app.HelpName = "Intelligent Junk and Spam Mail Handler"
	app.Version = "0.0.0"
	app.Copyright = "(c) 2017, Carlo Strub. All rights reserved. This binary is licensed under a BSD 3-Clause License."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Carlo Strub",
			Email: "cs@carlostrub.ch",
		},
	}

	maildirPaths := cli.StringSlice([]string{
		wd + "/Maildir",
	})
	app.Flags = []cli.Flag{

		cli.StringSliceFlag{
			Name:   "maildirs, d",
			Value:  &maildirPaths,
			EnvVar: "SISYPHUS_DIRS",
			Usage:  "Comma separated list of paths to the Maildir directories",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start sisyphus daemon",
			Action: func(c *cli.Context) error {
				if maildir == nil {
					return errors.New("no maildir selected")
				}

				// Load the Maildir
				mails, err := Index(maildirPaths[0])
				if err != nil {
					return cli.NewExitError(err, 66)
				}

				fmt.Println(mails)

				// Open the database
				db, err := openDB(maildirPaths[0])
				if err != nil {
					return cli.NewExitError(err, 66)
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

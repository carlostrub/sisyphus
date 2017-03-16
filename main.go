package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/urfave/cli"
)

const (
	good = "0"
	junk = "1"
)

func main() {

	// Define App
	app := cli.NewApp()
	app.Name = "Sisyphus"
	app.Usage = "Intelligent Junk and Spam Mail Handler"
	app.UsageText = `sisyphus [global options] command [command options]
	
	Sisyphus applies artificial intelligence to filter
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

	maildirPaths := cli.StringSlice([]string{})

	var pidfile *string
	pidfile = new(string)

	app.Flags = []cli.Flag{

		cli.StringSliceFlag{
			Name:   "maildirs, d",
			Value:  &maildirPaths,
			EnvVar: "SISYPHUS_DIRS",
			Usage:  "Comma separated list of paths to the Maildir directories",
		},
		cli.StringFlag{
			Name:        "pidfile, p",
			Value:       "/tmp/sisyphus.pid",
			EnvVar:      "SISYPHUS_PID",
			Usage:       "Location of PID file",
			Destination: pidfile,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"u"},
			Usage:   "run sisyphus",
			Action: func(c *cli.Context) {

				// check if daemon already running.
				if _, err := os.Stat(*pidfile); err == nil {
					log.Fatal("sisyphus running or " + *pidfile + " file exists.")
				}

				fmt.Print(`


	███████╗██╗███████╗██╗   ██╗██████╗ ██╗  ██╗██╗   ██╗███████╗
	██╔════╝██║██╔════╝╚██╗ ██╔╝██╔══██╗██║  ██║██║   ██║██╔════╝
	███████╗██║███████╗ ╚████╔╝ ██████╔╝███████║██║   ██║███████╗
	╚════██║██║╚════██║  ╚██╔╝  ██╔═══╝ ██╔══██║██║   ██║╚════██║
	███████║██║███████║   ██║   ██║     ██║  ██║╚██████╔╝███████║
	╚══════╝╚═╝╚══════╝   ╚═╝   ╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚══════╝


`)
				// Make arrangement to remove PID file upon receiving the SIGTERM from kill command
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

				go func() {
					signalType := <-ch
					signal.Stop(ch)
					log.Println("Exit command received. Exiting sisyphus...")

					// this is a good place to flush everything to disk
					// before terminating.
					log.Println("Received signal type: ", signalType)

					// remove PID file
					os.Remove(*pidfile)

					os.Exit(0)

				}()

				// Load the Maildir
				if len(maildirPaths) < 1 {
					log.Fatal("No Maildir set.")
				}
				if len(maildirPaths) > 1 {
					log.Fatal("Sorry... only one Maildir supported as of today.")
				}

				log.Println("create empty junk directory if required")
				os.MkdirAll(maildirPaths[0]+"/.Junk/cur", 0700)
				log.Println("loading mails")
				mails, err := Index(maildirPaths[0])
				if err != nil {
					log.Fatal("Wrong path to Maildir")
				}
				log.Println("mails loaded")

				// Open the database
				log.Println("loading database")
				db, err := openDB(maildirPaths[0])
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()
				log.Println("database loaded")

				// Check for unprocessed mail
				var unprocessedGood []string
				for i := range mails {
					db.View(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("Processed"))
						v := b.Get([]byte(mails[i].Key))
						if len(v) == 0 {
							unprocessedGood = append(unprocessedGood, mails[i].Key)
						}
						if string(v) == junk {
							unprocessedGood = append(unprocessedGood, mails[i].Key)
						}
						return nil
					})
				}

				// Classify and learn unprocessed mail

				mux := http.NewServeMux()
				log.Fatalln(http.ListenAndServe(":8080", mux))
			},
		},
		{
			// See
			// https://www.socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example
			// for the process we are using to daemonize
			Name:    "start",
			Aliases: []string{"s"},
			Usage:   "start sisyphus daemon in the background",
			Action: func(c *cli.Context) error {

				err := daemonStart(*pidfile)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"e"},
			Usage:   "stop sisyphus daemon",
			Action: func(c *cli.Context) error {

				err := daemonStop(*pidfile)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart sisyphus daemon",
			Action: func(c *cli.Context) error {

				err := daemonRestart(*pidfile)
				if err != nil {
					log.Fatal(err)
				}

				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"i"},
			Usage:   "status of sisyphus",
			Action: func(c *cli.Context) error {
				log.Println("here, we should get statistics from the db, TBD...")
				return nil
			},
		},
	}

	app.Run(os.Args)
}

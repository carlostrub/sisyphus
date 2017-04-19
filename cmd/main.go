package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/boltdb/bolt"
	"github.com/carlostrub/sisyphus"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
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
		{
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

	by Carlo Strub <cs@carlostrub.ch>


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

				sisyphus.CreateDirs(maildirPaths[0])

				mails, err := sisyphus.Index(maildirPaths[0])
				if err != nil {
					log.Fatal("Wrong path to Maildir")
				}

				// Open the database
				db, err := sisyphus.OpenDB(maildirPaths[0])
				if err != nil {
					log.Fatal(err)
				}
				defer db.Close()

				// Handle all mails after startup
				for i := range mails {
					db.View(func(tx *bolt.Tx) error {
						b := tx.Bucket([]byte("Processed"))
						bMails := b.Bucket([]byte("Mails"))
						v := bMails.Get([]byte(mails[i].Key))
						if len(v) == 0 {
							err = mails[i].Classify(db)
							if err != nil {
								log.Print(err)
							}
							err = mails[i].Learn(db)
							if err != nil {
								log.Print(err)
							}
						}
						if string(v) == sisyphus.Good && mails[i].Junk == true {
							err = mails[i].Learn(db)
							if err != nil {
								log.Print(err)
							}
						}
						if string(v) == sisyphus.Junk && mails[i].Junk == false {
							err = mails[i].Learn(db)
							if err != nil {
								log.Print(err)
							}
						}
						return nil
					})
				}

				// Handle mails as the arrive
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					log.Fatal(err)
				}
				defer watcher.Close()

				done := make(chan bool)
				go func() {
					for {
						select {
						case event := <-watcher.Events:
							if event.Op&fsnotify.Create == fsnotify.Create {
								mailName := strings.Split(event.Name, "/")
								m := sisyphus.Mail{
									Key: mailName[len(mailName)-1],
								}

								if mailName[len(mailName)-2] == "new" {
									err = m.Classify(db)
									if err != nil {
										log.Print(err)
									}
								} else {
									err = m.Learn(db)
									if err != nil {
										log.Print(err)
									}
								}

							}
						case err := <-watcher.Errors:
							log.Println("error:", err)
						}
					}
				}()

				err = watcher.Add(maildirPaths[0] + "/cur")
				if err != nil {
					log.Fatal(err)
				}
				err = watcher.Add(maildirPaths[0] + "/new")
				if err != nil {
					log.Fatal(err)
				}
				err = watcher.Add(maildirPaths[0] + "/.Junk/cur")
				if err != nil {
					log.Fatal(err)
				}
				<-done
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

				err := sisyphus.Pidfile(*pidfile).DaemonStart()
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

				err := sisyphus.Pidfile(*pidfile).DaemonStop()
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

				err := sisyphus.Pidfile(*pidfile).DaemonRestart()
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

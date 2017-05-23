package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/carlostrub/sisyphus"
	"github.com/urfave/cli"
)

var (
	version string
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
	app.Version = version
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
					log.WithFields(log.Fields{
						"pidfile": *pidfile,
					}).Fatal("Already running or pidfile exists")
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
					log.Info("Exit command received. Exiting sisyphus...")

					// this is a good place to flush everything to disk
					// before terminating.
					log.WithFields(log.Fields{
						"signal": signalType,
					}).Info("Received signal")

					// remove PID file
					os.Remove(*pidfile)

					os.Exit(0)

				}()

				if len(maildirPaths) < 1 {
					log.Fatal("No Maildir set. Please check the manual.")
				}

				// Populate maildir with the maildirs given by setting the flag.
				var maildirs []sisyphus.Maildir
				for _, val := range maildirPaths {
					maildirs = append(maildirs, sisyphus.Maildir(val))
				}

				// Load all mails
				mails, err := sisyphus.LoadMails(maildirs)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Fatal("Cannot load mails")
				}

				// Open all databases
				dbs, err := sisyphus.LoadDatabases(maildirs)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Fatal("Cannot load databases")
				}
				defer sisyphus.CloseDatabases(dbs)

				// Learn at startup
				for _, d := range maildirs {
					db := dbs[d]
					m := mails[d]
					for _, val := range m {
						err := val.Learn(db)
						if err != nil {
							log.WithFields(log.Fields{
								"err": err,
							}).Warning("Cannot learn mail")
						}
					}
				}

				//				// Classify on arrival
				//				watcher, err := fsnotify.NewWatcher()
				//				if err != nil {
				//log.WithFields(log.Fields{
				//	"err": err,
				//}).Warning("Cannot create directory watcher")
				//				}
				//				defer watcher.Close()
				//
				//				done := make(chan bool)
				//				go func() {
				//					for {
				//						select {
				//						case event := <-watcher.Events:
				//							if event.Op&fsnotify.Create == fsnotify.Create {
				//								mailName := strings.Split(event.Name, "/")
				//								m := sisyphus.Mail{
				//									Key: mailName[len(mailName)-1],
				//								}
				//
				//								if mailName[len(mailName)-2] == "new" {
				//									err = m.Classify(db)
				//									if err != nil {
				//										log.Print(err)
				//									}
				//								} else {
				//									err = m.Learn(db)
				//									if err != nil {
				//										log.Print(err)
				//									}
				//								}
				//
				//							}
				//						case err := <-watcher.Errors:
				//							log.Println("error:", err)
				//						}
				//					}
				//				}()
				//
				//				err = watcher.Add(maildirPaths[0] + "/cur")
				//				if err != nil {
				//					log.Fatal(err)
				//				}
				//				err = watcher.Add(maildirPaths[0] + "/new")
				//				if err != nil {
				//					log.Fatal(err)
				//				}
				//				err = watcher.Add(maildirPaths[0] + "/.Junk/cur")
				//				if err != nil {
				//					log.Fatal(err)
				//				}
				//				<-done
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

				sisyphus.Pidfile(*pidfile).DaemonStart()

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"e"},
			Usage:   "stop sisyphus daemon",
			Action: func(c *cli.Context) error {

				sisyphus.Pidfile(*pidfile).DaemonStop()

				return nil
			},
		},
		{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart sisyphus daemon",
			Action: func(c *cli.Context) error {

				sisyphus.Pidfile(*pidfile).DaemonRestart()

				return nil
			},
		},
		{
			Name:    "status",
			Aliases: []string{"i"},
			Usage:   "status of sisyphus",
			Action: func(c *cli.Context) error {
				log.Info("here, we should get statistics from the db, TBD...")
				return nil
			},
		},
	}

	app.Run(os.Args)
}

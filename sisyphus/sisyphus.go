package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"

	"github.com/carlostrub/sisyphus"
	"gopkg.in/urfave/cli.v2"
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
			Name:   "maildir, d",
			Value:  &maildirPaths,
			EnvVar: "SISYPHUS_DIRS",
			Usage:  "Call multiple Maildirs by repeating this flag, i.e. --maildir \"./Maildir\" --maildir \"./Maildir2\"",
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
						err := val.Learn(db, d)
						if err != nil {
							log.WithFields(log.Fields{
								"err":  err,
								"mail": val.Key,
							}).Error("Cannot learn mail")
						}
					}
				}

				// Classify whenever a mail arrives in "new"
				watcher, err := fsnotify.NewWatcher()
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Fatal("Cannot setup directory watcher")
				}
				defer watcher.Close()

				done := make(chan bool)
				go func() {
					for {
						select {
						case event := <-watcher.Events:
							if event.Op&fsnotify.Create == fsnotify.Create {
								path := strings.Split(event.Name, "/new/")
								m := sisyphus.Mail{
									Key: path[1],
								}

								err = m.Classify(dbs[sisyphus.Maildir(path[0])], sisyphus.Maildir(path[0]))
								if err != nil {
									log.WithFields(log.Fields{
										"err": err,
									}).Error("Classify mail")
								}

							}
						case err := <-watcher.Errors:
							log.WithFields(log.Fields{
								"err": err,
							}).Error("Problem with directory watcher")
						}
					}
				}()

				for _, val := range maildirPaths {
					err = watcher.Add(val + "/new")
					if err != nil {
						log.WithFields(log.Fields{
							"err": err,
							"dir": val + "/new",
						}).Error("Cannot watch directory")
					}
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

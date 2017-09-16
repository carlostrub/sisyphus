package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/boltdb/bolt"
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
	app.Usage = "Intelligent Junk Mail Handler"
	app.UsageText = `sisyphus [global options] command [command options]
	
	Sisyphus applies artificial intelligence to filter
	Junk mail in an unobtrusive way. Both, classification and learning
	operate directly on the Maildir of a user in a fully transparent mode,
	without any need for configuration or active operation.`
	app.HelpName = "Intelligent Junk Mail Handler"
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

	var learnafter *string
	learnafter = new(string)

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
		cli.StringFlag{
			Name:        "learn",
			Value:       "12h",
			EnvVar:      "SISYPHUS_DURATION",
			Usage:       "Time interval between to learn cycles",
			Destination: learnafter,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"u"},
			Usage:   "run sisyphus",
			Action: func(c *cli.Context) {

				fmt.Print(`


                                                                                           
               #####                                             
              #     # #  ####  #   # #####  #    # #    #  ####  
              #       # #       # #  #    # #    # #    # #      
               #####  #  ####    #   #    # ###### #    #  ####  
                    # #      #   #   #####  #    # #    #      # 
              #     # # #    #   #   #      #    # #    # #    # 
               #####  #  ####    #   #      #    #  ####   ####  

              by Carlo Strub <cs@carlostrub.ch>


`)

				if len(maildirPaths) < 1 {
					log.Fatal("No Maildir set. Please check the manual.")
				}

				// Populate maildir with the maildirs given by setting the flag.
				var maildirs []sisyphus.Maildir
				for _, val := range maildirPaths {
					maildirs = append(maildirs, sisyphus.Maildir(val))
				}

				// Create missing Maildirs
				err := sisyphus.LoadMaildirs(maildirs)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Fatal("Cannot load maildirs")
				}

				// Open all databases
				dbs, err := sisyphus.LoadDatabases(maildirs)
				if err != nil {
					log.WithFields(log.Fields{
						"err": err,
					}).Fatal("Cannot load databases")
				}
				defer sisyphus.CloseDatabases(dbs)

				// Learn at startup and regular intervals
				go func() {
					for {
						duration, err := time.ParseDuration(*learnafter)
						if err != nil {
							log.Fatal("Cannot parse duration for learning intervals.")
						}

						learn(maildirs, dbs)
						time.Sleep(duration)
					}
				}()

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
			Name:    "stats",
			Aliases: []string{"i"},
			Usage:   "Statistics from Sisyphus",
			Action: func(c *cli.Context) error {
				log.Info("here, we should get statistics from the db, TBD...")
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func learn(maildirs []sisyphus.Maildir, dbs map[sisyphus.Maildir]*bolt.DB) {
	mails, err := sisyphus.LoadMails(maildirs)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Cannot load mails")
	}
	for _, d := range maildirs {
		db := dbs[d]
		m := mails[d]
		for _, val := range m {
			err := val.Learn(db, d)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"mail": val.Key,
				}).Warning("Cannot learn mail")
			}
		}
	}
	log.Info("All mails learned")

	return
}

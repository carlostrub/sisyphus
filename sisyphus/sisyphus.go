package main

import (
	"bufio"
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
	app.UsageText = `
	Sisyphus applies artificial intelligence to filter Junk mail in an
	unobtrusive way. Both, classification and learning operate directly on
	the Maildir of a user in a fully transparent mode, without any need for
	configuration or active operation.`
	app.HelpName = "Intelligent Junk Mail Handler"
	app.Version = version
	app.Copyright = "(c) 2017, Carlo Strub. All rights reserved. This binary is licensed under a BSD 3-Clause License."
	app.Authors = []cli.Author{
		{
			Name:  "Carlo Strub",
			Email: "cs@carlostrub.ch",
		},
	}
	app.ExtraInfo = func() map[string]string {
		return map[string]string{
			"ENVIRONMENT VARIABLES": `For configuration, set the following environment
  variables:
  
  SISYPHUS_DIRS:     Comma-separated list of maildirs,
                     e.g. ./Maildir,/home/JohnDoe/Maildir

  SISYPHUS_DURATION: Interval between learning periods, e.g. 12h
			`,
		}
	}
	app.CustomAppHelpTemplate = `NAME:
  {{.Name}} - {{.Usage}}

USAGE:
  sisyphus {{if .VisibleFlags}}[FLAGS] {{end}}COMMAND{{if .VisibleFlags}}{{end}}
  {{.UsageText}}

COMMANDS:
  {{range .VisibleCommands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
{{range $key, $value := ExtraInfo}}
{{$key}}:
  {{$value}}
{{end}}VERSION:
  {{.Version}}

AUTHOR:{{range .Authors}}
  {{.}}{{end}}

COPYRIGHT:
  {{.Copyright}}
`

	dirsRaw, ok := os.LookupEnv("SISYPHUS_DIRS")
	if !ok {
		log.Fatal("Environment variable SISYPHUS_DIRS not set.")
	}
	dirsSplit := strings.Split(dirsRaw, ",")

	var maildirs []sisyphus.Maildir
	for i := 0; i < len(dirsSplit); i++ {
		maildirs = append(maildirs, sisyphus.Maildir(dirsSplit[i]))
	}

	_, ok = os.LookupEnv("SISYPHUS_DURATION")
	if !ok {
		log.Fatal("Environment variable SISYPHUS_DURATION not set.")
	}

	//	app.Flags = []cli.Flag{
	//
	//		&cli.StringSliceFlag{
	//			Name:    "maildir, d",
	//			Value:   &maildirPaths,
	//			EnvVars: []string{"SISYPHUS_DIRS"},
	//			Usage:   "Call multiple Maildirs by repeating this flag, i.e. --maildir \"./Maildir\" --maildir \"./Maildir2\"",
	//		},
	//		&cli.StringFlag{
	//			Name:        "learn",
	//			Value:       "12h",
	//			EnvVars:     []string{"SISYPHUS_DURATION"},
	//			Usage:       "Time interval between to learn cycles",
	//			Destination: learnafter,
	//		},
	//	}

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
						duration, err := time.ParseDuration(os.Getenv("SISYPHUS_DURATION"))
						if err != nil {
							log.Fatal("Cannot parse duration for learning intervals.")
						}

						backup(maildirs, dbs)
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

				for _, val := range maildirs {
					err = watcher.Add(string(val) + "/new")
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
			Usage:   "show statistics",
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

func backup(maildirs []sisyphus.Maildir, dbs map[sisyphus.Maildir]*bolt.DB) {
	for _, d := range maildirs {
		db := dbs[d]

		backup, err := os.Create(string(d) + "/sisyphus.db.backup")
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Backup creation")
		}
		defer backup.Close()

		w := bufio.NewWriter(backup)

		err = db.View(func(tx *bolt.Tx) error {
			_, err := tx.WriteTo(w)
			return err
		})
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("Backup creation")
		}

		w.Flush()
	}

	log.Info("All databases backed up successfully.")

	return
}

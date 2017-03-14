package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/urfave/cli"
)

var (
// Processed is a map of e-mail IDs and the value set to true if Junk
// Processed map[string]bool
)

func main() {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

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

	maildirPaths := cli.StringSlice([]string{
		wd + "/Maildir",
	})

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

				log.Print("App runs..........")
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
				// var maildir []string
				//				if maildir == nil {
				//					return errors.New("no maildir selected")
				//				}
				//
				//				// Load the Maildir
				//				mails, err := Index(maildirPaths[0])
				//				if err != nil {
				//					return cli.NewExitError(err, 66)
				//				}
				//
				//				fmt.Println(mails)
				//
				//				// Open the database
				//				db, err := openDB(maildirPaths[0])
				//				if err != nil {
				//					return cli.NewExitError(err, 66)
				//				}
				//				defer db.Close()

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
				// check if daemon already running.
				if _, err := os.Stat(*pidfile); err == nil {
					return cli.NewExitError("sisyphus running or "+*pidfile+" file exists.", 69)
				}

				cmd := exec.Command(os.Args[0], "run")
				cmd.Start()
				log.Printf("starting sisyphus process ID [%v]\n", cmd.Process.Pid)
				log.Println("sisyphus started")
				err := savePID(*pidfile, cmd.Process.Pid)
				if err != nil {
					return cli.NewExitError(err, 73)
				}

				return nil
			},
		},
		{
			Name:    "stop",
			Aliases: []string{"e"},
			Usage:   "stop sisyphus daemon",
			Action: func(c *cli.Context) error {

				_, err := os.Stat(*pidfile)
				if err != nil {
					return cli.NewExitError("sisyphus is not running", 64)
				}

				processIDRaw, err := ioutil.ReadFile(*pidfile)
				if err != nil {
					return cli.NewExitError("sisyphus is not running", 64)
				}

				processID, err := strconv.Atoi(string(processIDRaw))
				if err != nil {
					return cli.NewExitError("unable to read and parse process id found in "+*pidfile, 74)
				}

				process, err := os.FindProcess(processID)

				if err != nil {
					e := fmt.Sprintf("Unable to find process ID [%v] with error %v \n", processID, err)
					return cli.NewExitError(e, 71)
				}

				// remove PID file
				os.Remove(*pidfile)

				log.Printf("stopping sisyphus process ID [%v]\n", processID)
				// kill process and exit immediately
				err = process.Kill()

				if err != nil {
					e := fmt.Sprintf("Unable to kill process ID [%v] with error %v \n", processID, err)
					return cli.NewExitError(e, 71)
				}

				log.Println("sisyphus stopped")
				os.Exit(0)

				return nil
			},
		},
		{
			Name:    "restart",
			Aliases: []string{"r"},
			Usage:   "restart sisyphus daemon",
			Action: func(c *cli.Context) error {
				_, err := os.Stat(*pidfile)
				if err != nil {
					return cli.NewExitError("sisyphus is not running", 64)
				}

				pid, err := ioutil.ReadFile(*pidfile)
				if err != nil {
					return cli.NewExitError("sisyphus is not running", 64)
				}

				cmd := exec.Command(os.Args[0], "stop")
				err = cmd.Start()
				if err != nil {
					return cli.NewExitError(err, 64)
				}
				log.Printf("stopping sisyphus process ID [%v]\n", string(pid))

				cmd = exec.Command(os.Args[0], "start")
				err = cmd.Start()
				if err != nil {
					return cli.NewExitError(err, 64)
				}

				log.Println("sisyphus restarted")

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

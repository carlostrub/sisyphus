package sisyphus

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// See
// https://www.socketloop.com/tutorials/golang-daemonizing-a-simple-web-server-process-example
// for the process we are using to daemonize

// Pidfile holds the Process ID file of sisyphus
type Pidfile string

// savePID stores a pidfile
func (p Pidfile) savePID(process int) error {
	file, err := os.Create(string(p))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(process))
	if err != nil {
		return err
	}

	err = file.Sync()

	return err
}

// DaemonStart starts sisyphus as a backgound process
func (p Pidfile) DaemonStart() {
	// check if daemon already running.
	if _, err := os.Stat(string(p)); err == nil {

		log.WithFields(log.Fields{
			"pidfile": p,
		}).Fatal("Already running or pidfile exists")

	}

	log.Info("Starting sisyphus daemon")
	cmd := exec.Command(os.Args[0], "run")
	cmd.Start()

	log.WithFields(log.Fields{
		"pid": cmd.Process.Pid,
	}).Info("Sisyphus started")
	err := (p).savePID(cmd.Process.Pid)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Save process ID file")
	}
	log.WithFields(log.Fields{
		"pidfile": p,
	}).Info("Process ID file stored")

	return
}

// DaemonStop stops a running sisyphus background process
func (p Pidfile) DaemonStop() {

	_, err := os.Stat(string(p))
	if err != nil {
		log.Fatal("Sisyphus is not running")
	}

	processIDRaw, err := ioutil.ReadFile(string(p))
	if err != nil {
		log.Fatal("Sisyphus is not running")
	}

	processID, err := strconv.Atoi(string(processIDRaw))
	if err != nil {
		log.WithFields(log.Fields{
			"pid": p,
		}).Fatal("Unable to read process ID")
	}

	process, err := os.FindProcess(processID)
	if err != nil {
		log.WithFields(log.Fields{
			"pid": p,
			"err": err,
		}).Fatal("Unable to find process ID")
	}

	// remove PID file
	err = os.Remove(string(p))
	if err != nil {
		log.Error("Unable to remove process ID file")
	}

	log.WithFields(log.Fields{
		"pid": processID,
	}).Info("Stopping sisyphus process")
	// kill process and exit immediately
	err = process.Kill()
	if err != nil {
		log.WithFields(log.Fields{
			"pid": processID,
			"err": err,
		}).Fatal("Unable to kill sisyphus process")
	}

	log.Info("Sisyphus stopped")
	os.Exit(0)

	return
}

// DaemonRestart restarts a running sisyphus background process
func (p Pidfile) DaemonRestart() {
	_, err := os.Stat(string(p))
	if err != nil {
		log.Fatal("Sisyphus not running")
	}

	pid, err := ioutil.ReadFile(string(p))
	if err != nil {
		log.Fatal("Sisyphus not running")
	}

	log.WithFields(log.Fields{
		"pid": string(pid),
	}).Info("Stopping sisyphus process")
	cmd := exec.Command(os.Args[0], "stop")
	cmd.Start()

	cmd = exec.Command(os.Args[0], "start")
	cmd.Start()

	log.Info("Sisyphus restarted")

	return
}

package sisyphus

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
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
func (p Pidfile) DaemonStart() error {
	// check if daemon already running.
	if _, err := os.Stat(string(p)); err == nil {
		return errors.New("sisyphus running or " + string(p) + " file exists.")
	}

	cmd := exec.Command(os.Args[0], "run")
	err := cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("starting sisyphus process ID [%v]\n", cmd.Process.Pid)
	log.Println("sisyphus started")
	err = (p).savePID(cmd.Process.Pid)

	return err
}

// DaemonStop stops a running sisyphus background process
func (p Pidfile) DaemonStop() error {

	_, err := os.Stat(string(p))
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	processIDRaw, err := ioutil.ReadFile(string(p))
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	processID, err := strconv.Atoi(string(processIDRaw))
	if err != nil {
		return errors.New("unable to read and parse process id found in " + string(p))
	}

	process, err := os.FindProcess(processID)
	if err != nil {
		e := fmt.Sprintf("Unable to find process ID [%v] with error %v \n", processID, err)
		return errors.New(e)
	}

	// remove PID file
	err = os.Remove(string(p))
	if err != nil {
		return err
	}

	log.Printf("stopping sisyphus process ID [%v]\n", processID)
	// kill process and exit immediately
	err = process.Kill()
	if err != nil {
		e := fmt.Sprintf("Unable to kill process ID [%v] with error %v \n", processID, err)
		return errors.New(e)
	}

	log.Println("sisyphus stopped")
	os.Exit(0)

	return nil
}

// DaemonRestart restarts a running sisyphus background process
func (p Pidfile) DaemonRestart() error {
	_, err := os.Stat(string(p))
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	pid, err := ioutil.ReadFile(string(p))
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	cmd := exec.Command(os.Args[0], "stop")
	err = cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("stopping sisyphus process ID [%v]\n", string(pid))

	cmd = exec.Command(os.Args[0], "start")
	err = cmd.Start()
	if err != nil {
		return err
	}

	log.Println("sisyphus restarted")

	return nil
}

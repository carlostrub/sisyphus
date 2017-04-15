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

// savePID stores a pidfile
func savePID(pidfile string, p int) error {
	file, err := os.Create(pidfile)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(p))
	if err != nil {
		return err
	}

	file.Sync()

	return nil
}

func daemonStart(pidfile string) error {
	// check if daemon already running.
	if _, err := os.Stat(pidfile); err == nil {
		return errors.New("sisyphus running or " + pidfile + " file exists.")
	}

	cmd := exec.Command(os.Args[0], "run")
	cmd.Start()
	log.Printf("starting sisyphus process ID [%v]\n", cmd.Process.Pid)
	log.Println("sisyphus started")
	err := savePID(pidfile, cmd.Process.Pid)
	if err != nil {
		return err
	}

	return nil
}

func daemonStop(pidfile string) error {

	_, err := os.Stat(pidfile)
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	processIDRaw, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	processID, err := strconv.Atoi(string(processIDRaw))
	if err != nil {
		return errors.New("unable to read and parse process id found in " + pidfile)
	}

	process, err := os.FindProcess(processID)

	if err != nil {
		e := fmt.Sprintf("Unable to find process ID [%v] with error %v \n", processID, err)
		return errors.New(e)
	}

	// remove PID file
	os.Remove(pidfile)

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

func daemonRestart(pidfile string) error {
	_, err := os.Stat(pidfile)
	if err != nil {
		return errors.New("sisyphus is not running")
	}

	pid, err := ioutil.ReadFile(pidfile)
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

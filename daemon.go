package main

import (
	"os"
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

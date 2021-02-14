package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// catch signals and terminate the app
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// monitor for signals in the background
	go func() {
		s := <-sigc
		fmt.Println("\nreceived signal:", s)
		os.Exit(0)
	}()



}
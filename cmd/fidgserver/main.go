package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log/level"

	"timefidget/pkg/fidgserver"
	"timefidget/pkg/util"
)

func main() {

	// catch signals and terminate the app
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	fs, err := fidgserver.NewFidgserver()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// monitor for signals in the background
	go func() {
		s := <-sigc
		log.Println("msg", "fidgserver stopping", "signal", s)
		fs.Stop()
		log.Println("msg", "fidgserver stopped")
		os.Exit(0)
	}()

	level.Info(util.Logger).Log("msg", "fidgserver running")

	for {
	}

}

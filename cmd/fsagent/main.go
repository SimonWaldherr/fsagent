package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/kardianos/service"
	"simonwaldherr.de/go/fsagent"
)

var logger service.Logger

func main() {
	svcConfig := &service.Config{
		Name:        "fsagent",
		DisplayName: "FileSystem Agent",
		Description: "this service can monitor a folder and do configurable things on filesystem triggers.",
	}

	prg := &fsagent.Program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Signal: %v\n", sig)
			prg.Stop(s)
		}
	}()

	fmt.Println("run ...")
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

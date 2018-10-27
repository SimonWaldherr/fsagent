package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"

	"simonwaldherr.de/go/golibs/file"
)

type program struct {
	stop chan struct{}
}

func (p *program) Start(s service.Service) error {
	// done is a global channel
	done = make(chan bool)
	p.stop = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	close(p.stop)
	<-done
	wg.Wait()
	return nil
}

func (p *program) run() {
	var config []Config
	str, _ := file.Read(os.Args[1])

	err := json.Unmarshal([]byte(str), &config)

	if err != nil {
		log.Println(err)
	}

	watcher := make(map[string]*fsnotify.Watcher)

	for i, conf := range config {
		fmt.Println("load config ...")

		go runConfig(&wg, conf, i, p.stop, watcher)
	}
	done <- true
}

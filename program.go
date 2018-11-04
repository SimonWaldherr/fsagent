package fsagent

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"

	"simonwaldherr.de/go/golibs/file"
)

var Done chan bool
var WG sync.WaitGroup

type Program struct {
	stop chan struct{}
}

func (p *Program) Start(s service.Service) error {
	// done is a global channel
	Done = make(chan bool)
	p.stop = make(chan struct{})
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	close(p.stop)
	<-Done
	WG.Wait()
	return nil
}

func (p *Program) run() {
	var config []Config
	str, _ := file.Read(os.Args[1])

	err := json.Unmarshal([]byte(str), &config)

	if err != nil {
		log.Println(err)
	}

	watcher := make(map[string]*fsnotify.Watcher)

	for i, conf := range config {
		fmt.Println("load config ...")

		go runConfig(&WG, conf, i, p.stop, watcher)
	}
	Done <- true
}

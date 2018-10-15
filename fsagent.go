package main

import (
	"encoding/json"
	"fmt"
	"github.com/SimonWaldherr/gwv"
	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/service"
	"io/ioutil"
	"strings"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"simonwaldherr.de/go/fsagent/modules"
	"simonwaldherr.de/go/golibs/cache"
	"simonwaldherr.de/go/golibs/file"
	"simonwaldherr.de/go/golibs/regex"
	"simonwaldherr.de/go/golibs/xtime"
	"sync"
	"time"
)

// Config represents an element of the application configuration.
type Config struct {
	Verbose  bool   `json:"verbose"`
	Debounce bool   `json:"debounce"`
	Folder   string `json:"folder"`
	Trigger  string `json:"trigger"`
	Ticker   int    `json:"ticker"`
	OnlyNew  bool   `json:"onlynew"`
	Match    string `json:"match"`
	Action   Action `json:"action"`
}

// Action is something that should be performed.
type Action []struct {
	Do        string          `json:"do"`
	Config    json.RawMessage `json:"config"`
	Onsuccess Action          `json:"onSuccess"`
	Onfailure Action          `json:"onFailure"`
}

// Actionable defines a common set of methods each action should have.
type Actionable interface {
	Name() string
	EmptyConfig() interface{}
	Perform(config interface{}, fileName string) error
}

// Actions is a list of Actionable types that can be added to.
var Actions = []Actionable{
	&modules.Mail{},
	&modules.HTTP{},
	&modules.Sleep{},
	&modules.Delete{},
	&modules.Move{},
	&modules.Decompress{},
	&modules.Compress{},
	&modules.IsFile{},
}

func do(act Action, file string) {
	for _, a := range act {
		var err error
		log.Printf("Do \"%v\" on file \"%v\"\n", a.Do, file)

		for _, action := range Actions {
			if a.Do == action.Name() {
				go func() {
					fmt.Printf("run action %v on file %v at %v", a.Do, file, xtime.Fmt("%Y-%m-%d %H:%M:%S", time.Now()))
					hub.Messages <- fmt.Sprintf("run action <span title=\"%s\">%v</span> on file %v at %v\n", strings.Replace(fmt.Sprintf("%s", a.Config), "\"", "&quot;", -1), a.Do, file, xtime.Fmt("%Y-%m-%d %H:%M:%S", time.Now()))
				}()
				config := action.EmptyConfig()
				json.Unmarshal(a.Config, &config)
				err = action.Perform(config, file)
				break
			}
		}

		if err == nil {
			do(a.Onsuccess, file)
		} else {
			log.Printf("Error on \"%v\" for file \"%v\"", err, file)
			do(a.Onfailure, file)
		}
	}
}

var done chan bool
var stop bool
var wg sync.WaitGroup

// FileLastModified returns the Time a file was last modified.
func FileLastModified(filename string) (*time.Time, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	statInfo, _ := f.Stat()
	modTime := statInfo.ModTime()
	return &modTime, nil
}

func (p *program) run() {
	var config []Config
	str, _ := file.Read(os.Args[1])

	err := json.Unmarshal([]byte(str), &config)

	if err != nil {
		log.Println(err)
	}

	watcher := make(map[string]*fsnotify.Watcher)

	i := 0

	for _, conf := range config {
		fmt.Println("load config ...")

		go func(configuration Config) {
			var timer *time.Ticker
			var eventCache *cache.Cache

			fmt.Println("start loop ...")

			if conf.Debounce {
				eventCache = cache.New(3*time.Second, 1*time.Second)
			}

			switch conf.Trigger {
			case "fsevent":
				watcher[conf.Folder+fmt.Sprintf(":%v", i)], err = fsnotify.NewWatcher()
				defer watcher[conf.Folder+fmt.Sprintf(":%v", i)].Close()
				fmt.Printf("Config FS-Watcher: %#v\n", watcher)
			case "ticker":
				timer = time.NewTicker(time.Millisecond * time.Duration(conf.Ticker))
			}

			wg.Add(1)
			defer func() {
				wg.Done()
			}()

			for !stop {
				fmt.Println("waiting for events")
				switch configuration.Trigger {
				case "fsevent":
					select {
					case event := <-watcher[configuration.Folder+fmt.Sprintf(":%v", i)].Events:
						fmt.Println("Event detected!")
						_, filename := filepath.Split(event.Name)
						match, _ := regex.MatchString(filename, configuration.Match)
						if event.Op != fsnotify.Remove && match {
							cv := fmt.Sprintf("%v:%v:%v", configuration.Folder, i, event.Name)
							if configuration.Debounce && eventCache.Get(cv) == nil {
								eventCache.Set(cv, true)
								do(configuration.Action, event.Name)
							}
						}
					case err := <-watcher[configuration.Folder+fmt.Sprintf(":%v", i)].Errors:
						log.Println("error:", err)
					}
				case "http":
					fmt.Println("coming soon.")
				case "ticker":
					select {
					case _ = <-timer.C:
						Folder := xtime.Fmt(configuration.Folder, time.Now())
						TickerFiles, _ := ioutil.ReadDir(Folder)

						for _, f := range TickerFiles {
							match, _ := regex.MatchString(f.Name(), configuration.Match)
							if match {
								if configuration.OnlyNew {
									lastmod, _ := FileLastModified(f.Name())
									if lastmod.Unix() < time.Now().Add(time.Millisecond*time.Duration(conf.Ticker)*-1).Unix() {
										continue
									}
								}

								fmt.Println(f.Name())
								cv := fmt.Sprintf("%v:%v:%v", Folder, i, f.Name())
								if configuration.Debounce && eventCache.Get(cv) == nil {
									eventCache.Set(cv, true)
									go do(configuration.Action, Folder+f.Name())
								}
							}
						}
					}
				default:
					fmt.Println("no or incompatible trigger configured")
				}
			}
		}(conf)
	}
	done <- true
}

var logger service.Logger
var hub *gwv.Connections

type program struct{}

func (p *program) Start(s service.Service) error {
	done = make(chan bool)
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	stop = true
	<-done
	wg.Wait()
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "fsagent",
		DisplayName: "FileSystem Agent",
		Description: "this service can monitor a folder and do configurable things on filesystem triggers.",
	}

	prg := &program{}
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

	go func() {
		hub = startWebGui()
	}()

	fmt.Println("run ...")
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

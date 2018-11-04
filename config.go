package fsagent

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"simonwaldherr.de/go/golibs/cache"
	"simonwaldherr.de/go/golibs/regex"
	"simonwaldherr.de/go/golibs/xtime"
)

// Config represents an element of the application configuration.
type Config struct {
	Folder   string `json:"folder"`
	Trigger  string `json:"trigger"`
	Ticker   int    `json:"ticker"`
	Match    string `json:"match"`
	Action   Action `json:"action"`
	OnlyNew  bool   `json:"onlynew"`
	Verbose  bool   `json:"verbose"`
	Debounce bool   `json:"debounce"`
}

func runConfig(wg *sync.WaitGroup, conf Config, i int, stop chan struct{}, watcher map[string]*fsnotify.Watcher) {
	var timer *time.Ticker
	var eventCache *cache.Cache

	fmt.Println("start loop ...")

	if conf.Debounce {
		eventCache = cache.New(3*time.Second, 1*time.Second)
	}

	folderidx := fmt.Sprintf("%s:%v", conf.Folder, i)

	switch conf.Trigger {
	case "fsevent":
		watcher[folderidx], _ = fsnotify.NewWatcher()
		defer watcher[folderidx].Close()
		fmt.Printf("Config FS-Watcher: %#v\n", watcher)
	case "ticker":
		timer = time.NewTicker(time.Millisecond * time.Duration(conf.Ticker))
	}

	wg.Add(1)
	defer wg.Done()

	for {
		fmt.Println("waiting for events")
		switch conf.Trigger {
		case "fsevent":
			handleFsEvent(conf, eventCache, i, watcher)
		case "http":
			fmt.Println("coming soon.")
		case "ticker":
			handleTicker(conf, eventCache, i, timer)
		default:
			fmt.Println("no or incompatible trigger configured")
			return
		}

		select {
		case <-stop:
			return
		default:
		}
	}

}

func handleTicker(conf Config, eventCache *cache.Cache, i int, timer *time.Ticker) {
	select {
	case _ = <-timer.C:
		Folder := xtime.Fmt(conf.Folder, time.Now())
		TickerFiles, _ := ioutil.ReadDir(Folder)

		for _, f := range TickerFiles {
			match, _ := regex.MatchString(f.Name(), conf.Match)
			if match {
				if conf.OnlyNew {
					lastmod, _ := FileLastModified(f.Name())
					if lastmod.Unix() < time.Now().Add(time.Millisecond*time.Duration(conf.Ticker)*-1).Unix() {
						continue
					}
				}

				fmt.Println(f.Name())
				cv := fmt.Sprintf("%v:%v:%v", Folder, i, f.Name())
				if conf.Debounce && eventCache.Get(cv) == nil {
					eventCache.Set(cv, true)
					go do(conf.Action, Folder+f.Name())
				}
			}
		}
	}

}

func handleFsEvent(conf Config, eventCache *cache.Cache, i int, watcher map[string]*fsnotify.Watcher) {
	select {
	case event := <-watcher[conf.Folder+fmt.Sprintf(":%v", i)].Events:
		fmt.Println("Event detected!")
		_, filename := filepath.Split(event.Name)
		match, _ := regex.MatchString(filename, conf.Match)
		if event.Op != fsnotify.Remove && match {
			cv := fmt.Sprintf("%v:%v:%v", conf.Folder, i, event.Name)
			if conf.Debounce && eventCache.Get(cv) == nil {
				eventCache.Set(cv, true)
				do(conf.Action, event.Name)
			}
		}
	case err := <-watcher[conf.Folder+fmt.Sprintf(":%v", i)].Errors:
		log.Println("error:", err)
	}
}

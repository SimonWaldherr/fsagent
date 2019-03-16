package fsagent

import (
	"encoding/json"
	"log"
	"time"

	"github.com/SimonWaldherr/golibs/cachedfile"
	gfile "github.com/SimonWaldherr/golibs/file"
	"simonwaldherr.de/go/fsagent/modules"
)

func init() {
	cachedfile.Init(15*time.Minute, 1*time.Minute)
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
				config := action.EmptyConfig()
				if gfile.IsFile(string(a.Config)) {
					str, _ := cachedfile.Read(string(a.Config))
					json.Unmarshal([]byte(str), &config)
				} else {
					json.Unmarshal(a.Config, &config)
				}

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

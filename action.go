package fsagent

import (
	"encoding/json"
	"log"

	"simonwaldherr.de/go/fsagent/modules"
)

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

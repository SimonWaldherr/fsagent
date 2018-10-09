package modules

import (
	"time"
)

type sleepConfig struct {
	Time int `json:"time"`
}

type Sleep struct{}

func (Sleep) Name() string {
	return "sleep"
}

func (Sleep) EmptyConfig() interface{} {
	return &sleepConfig{}
}

func (Sleep) Perform(config interface{}, fileName string) error {
	c := config.(sleepConfig)

	time.Sleep(time.Millisecond * time.Duration(c.Time))

	return nil
}

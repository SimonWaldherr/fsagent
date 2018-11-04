package modules

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"simonwaldherr.de/go/golibs/file"
)

type httpConfig struct {
	Path string `json:"path"`
}

type HTTP struct{}

func (HTTP) Name() string {
	return "http"
}

func (HTTP) EmptyConfig() interface{} {
	return &httpConfig{}
}

func (HTTP) Perform(config interface{}, fileName string) error {
	c := config.(*httpConfig)

	client := &http.Client{}
	str, _ := file.Read(fileName)
	body := bytes.NewBufferString(str)
	clength := strconv.Itoa(len(str))
	r, _ := http.NewRequest("POST", c.Path, body)
	r.Header.Add("User-Agent", "FS-Agent")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", clength)

	rsp, err := client.Do(r)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("the remote did not return a HTTP 200 response: %#v", rsp)
}

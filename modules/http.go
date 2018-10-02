package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"simonwaldherr.de/go/golibs/file"
	"strconv"
)

type httpConfig struct {
	Path string `json:"path"`
}

func HttpPostRequest(configName, fileName string) error {
	var c httpConfig
	str, _ := file.Read(configName)
	err := json.Unmarshal([]byte(str), &c)

	client := &http.Client{}
	str, _ = file.Read(fileName)
	body := bytes.NewBufferString(str)
	clength := strconv.Itoa(len(str))
	r, _ := http.NewRequest("POST", c.Path, body)
	r.Header.Add("User-Agent", "FS-Agent")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", clength)

	rsp, err := client.Do(r)

	if err != nil {
		return err
	} else {
		defer func() {
			rsp.Body.Close()
		}()
		if rsp.StatusCode == 200 {
			return err
		} else {
			return fmt.Errorf("The remote end did not return a HTTP 200 (OK) response:%#v\n", rsp)
		}
	}
}

package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type dagConfig struct {
	URL             string            `json:"url"`
	DagID           string            `json:"dagId"`
	Headers         map[string]string `json:"headers"`
	Timeout         int               `json:"timeout"`
	MaxPayloadBytes int               `json:"maxPayloadBytes"`
}

type DAG struct{}

func (DAG) Name() string {
	return "dag"
}

func (DAG) EmptyConfig() interface{} {
	return &dagConfig{}
}

func (DAG) Perform(config interface{}, fileName string) error {
	c := config.(*dagConfig)
	if c.URL == "" {
		return fmt.Errorf("url missing")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(map[string]interface{}{
		"dag_id": c.DagID,
		"conf": map[string]string{
			"filename": fileName,
			"content":  string(content),
		},
	})
	if err != nil {
		return err
	}
	if c.MaxPayloadBytes > 0 && len(payload) > c.MaxPayloadBytes {
		return fmt.Errorf("payload exceeds maxPayloadBytes (%v > %v)", len(payload), c.MaxPayloadBytes)
	}

	timeout := 60
	if c.Timeout > 0 {
		timeout = c.Timeout
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		b, _ := io.ReadAll(rsp.Body)
		return fmt.Errorf("dag request to %q failed with status %v: %s", c.URL, rsp.StatusCode, string(b))
	}

	return nil
}

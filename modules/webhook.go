package modules

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type webhookConfig struct {
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	ContentType string            `json:"contentType"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Timeout     int               `json:"timeout"`
}

type Webhook struct{}

func (Webhook) Name() string {
	return "webhook"
}

func (Webhook) EmptyConfig() interface{} {
	return &webhookConfig{}
}

func (Webhook) Perform(config interface{}, fileName string) error {
	c := config.(*webhookConfig)
	if c.URL == "" {
		return fmt.Errorf("url missing")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	method := c.Method
	if method == "" {
		method = "POST"
	}
	contentType := c.ContentType
	if contentType == "" {
		contentType = "text/plain"
	}

	body := string(content)
	if c.Body != "" {
		body = strings.Replace(c.Body, "$filecontent", string(content), -1)
		body = strings.Replace(body, "$filename", fileName, -1)
	}

	timeout := 60
	if c.Timeout > 0 {
		timeout = c.Timeout
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest(method, c.URL, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return fmt.Errorf("webhook request failed with status %v", rsp.StatusCode)
	}

	return nil
}

package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type llmConfig struct {
	URL          string `json:"url"`
	Model        string `json:"model"`
	Prompt       string `json:"prompt"`
	SystemPrompt string `json:"systemPrompt"`
	APIKey       string `json:"apiKey"`
	Output       string `json:"output"`
	Timeout      int    `json:"timeout"`
}

type OpenAI struct{}
type LMStudio struct{}
type Ollama struct{}

func (OpenAI) Name() string {
	return "openai"
}

func (LMStudio) Name() string {
	return "lmstudio"
}

func (Ollama) Name() string {
	return "ollama"
}

func (OpenAI) EmptyConfig() interface{} {
	return &llmConfig{}
}

func (LMStudio) EmptyConfig() interface{} {
	return &llmConfig{}
}

func (Ollama) EmptyConfig() interface{} {
	return &llmConfig{}
}

func llmOutputFile(fileName, output string) string {
	if output == "" {
		return fmt.Sprintf("%v.ai.txt", fileName)
	}
	return strings.Replace(output, "$file", fileName, -1)
}

func callOpenAICompatible(url string, c *llmConfig, fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	type llmMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	type llmRequest struct {
		Model    string       `json:"model"`
		Messages []llmMessage `json:"messages"`
	}

	type llmResponse struct {
		Choices []struct {
			Message llmMessage `json:"message"`
		} `json:"choices"`
	}

	messages := []llmMessage{}
	if c.SystemPrompt != "" {
		messages = append(messages, llmMessage{
			Role:    "system",
			Content: c.SystemPrompt,
		})
	}
	messages = append(messages, llmMessage{
		Role:    "user",
		Content: strings.TrimSpace(c.Prompt + "\n\n" + string(fileContent)),
	})

	body, err := json.Marshal(llmRequest{
		Model:    c.Model,
		Messages: messages,
	})
	if err != nil {
		return err
	}

	timeout := 60
	if c.Timeout > 0 {
		timeout = c.Timeout
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return fmt.Errorf("llm request failed with status %v", rsp.StatusCode)
	}

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var parsed llmResponse
	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}
	if len(parsed.Choices) == 0 || parsed.Choices[0].Message.Content == "" {
		return fmt.Errorf("llm response was empty")
	}

	return os.WriteFile(llmOutputFile(fileName, c.Output), []byte(parsed.Choices[0].Message.Content), 0600)
}

func (OpenAI) Perform(config interface{}, fileName string) error {
	c := config.(*llmConfig)
	url := c.URL
	if url == "" {
		url = "https://api.openai.com/v1/chat/completions"
	}
	return callOpenAICompatible(url, c, fileName)
}

func (LMStudio) Perform(config interface{}, fileName string) error {
	c := config.(*llmConfig)
	url := c.URL
	if url == "" {
		url = "http://localhost:1234/v1/chat/completions"
	}
	return callOpenAICompatible(url, c, fileName)
}

func (Ollama) Perform(config interface{}, fileName string) error {
	c := config.(*llmConfig)
	url := c.URL
	if url == "" {
		url = "http://localhost:11434/api/generate"
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		System string `json:"system,omitempty"`
		Stream bool   `json:"stream"`
	}{
		Model:  c.Model,
		Prompt: strings.TrimSpace(c.Prompt + "\n\n" + string(fileContent)),
		System: c.SystemPrompt,
		Stream: false,
	})
	if err != nil {
		return err
	}

	timeout := 60
	if c.Timeout > 0 {
		timeout = c.Timeout
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return fmt.Errorf("ollama request failed with status %v", rsp.StatusCode)
	}

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var parsed struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}
	if parsed.Response == "" {
		return fmt.Errorf("ollama response was empty")
	}

	return os.WriteFile(llmOutputFile(fileName, c.Output), []byte(parsed.Response), 0600)
}

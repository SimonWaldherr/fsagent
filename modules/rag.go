package modules

import (
	"os"
	"strings"
)

type ragConfig struct {
	URL          string   `json:"url"`
	Model        string   `json:"model"`
	Prompt       string   `json:"prompt"`
	SystemPrompt string   `json:"systemPrompt"`
	APIKey       string   `json:"apiKey"`
	Output       string   `json:"output"`
	Timeout      int      `json:"timeout"`
	ContextFiles []string `json:"contextFiles"`
}

type RAG struct{}

func (RAG) Name() string {
	return "rag"
}

func (RAG) EmptyConfig() interface{} {
	return &ragConfig{}
}

func (RAG) Perform(config interface{}, fileName string) error {
	c := config.(*ragConfig)
	url := c.URL
	if url == "" {
		url = "http://localhost:1234/v1/chat/completions"
	}

	contextParts := make([]string, 0, len(c.ContextFiles))
	for _, contextFile := range c.ContextFiles {
		contextContent, err := os.ReadFile(contextFile)
		if err == nil {
			contextParts = append(contextParts, string(contextContent))
		}
	}

	mergedPrompt := c.Prompt
	if len(contextParts) > 0 {
		mergedPrompt = strings.TrimSpace(mergedPrompt + "\n\nContext:\n" + strings.Join(contextParts, "\n\n"))
	}

	llmConf := &llmConfig{
		Model:        c.Model,
		Prompt:       mergedPrompt,
		SystemPrompt: c.SystemPrompt,
		APIKey:       c.APIKey,
		Output:       c.Output,
		Timeout:      c.Timeout,
	}

	return callOpenAICompatible(url, llmConf, fileName)
}

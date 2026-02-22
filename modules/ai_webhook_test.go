package modules

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLLMOutputFile(t *testing.T) {
	if got := llmOutputFile("a.txt", ""); got != "a.txt.ai.txt" {
		t.Fatalf("unexpected default output path: %v", got)
	}
	if got := llmOutputFile("a.txt", "$file.out"); got != "a.txt.out" {
		t.Fatalf("unexpected templated output path: %v", got)
	}
}

func TestCallOpenAICompatible(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	if err := ioutil.WriteFile(in, []byte("hello"), 0600); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"ok"}}]}`))
	}))
	defer server.Close()

	c := &llmConfig{
		Model:  "test-model",
		Prompt: "do this",
		Output: "$file.result",
	}

	if err := callOpenAICompatible(server.URL, c, in); err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadFile(in + ".result")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "ok" {
		t.Fatalf("unexpected output: %s", string(out))
	}
}

func TestWebhookPerform(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	if err := os.WriteFile(in, []byte("payload"), 0600); err != nil {
		t.Fatal(err)
	}

	var received string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		received = string(b)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	w := Webhook{}
	err := w.Perform(&webhookConfig{
		URL:  server.URL,
		Body: `{"name":"$filename","content":"$filecontent"}`,
	}, in)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(received, "payload") || !strings.Contains(received, "in.txt") {
		t.Fatalf("unexpected webhook body: %v", received)
	}
}

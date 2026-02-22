package modules

import (
	"encoding/json"
	"io"
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
	if err := os.WriteFile(in, []byte("hello"), 0600); err != nil {
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

	out, err := os.ReadFile(in + ".result")
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
		b, _ := io.ReadAll(r.Body)
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

func TestRAGPerformIncludesContext(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	context := dir + "/context.txt"

	if err := os.WriteFile(in, []byte("question"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(context, []byte("retrieved facts"), 0600); err != nil {
		t.Fatal(err)
	}

	var requestBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		requestBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"choices":[{"message":{"content":"rag-result"}}]}`))
	}))
	defer server.Close()

	rag := RAG{}
	err := rag.Perform(&ragConfig{
		URL:          server.URL,
		Model:        "model",
		Prompt:       "answer",
		ContextFiles: []string{context},
		Output:       "$file.rag.out",
	}, in)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(requestBody, "retrieved facts") {
		t.Fatalf("context not included in request body: %s", requestBody)
	}
	out, err := os.ReadFile(in + ".rag.out")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out)) != "rag-result" {
		t.Fatalf("unexpected rag output: %s", out)
	}
}

func TestDAGPerform(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	if err := os.WriteFile(in, []byte("dag-payload"), 0600); err != nil {
		t.Fatal(err)
	}

	var payload map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &payload)
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	dag := DAG{}
	err := dag.Perform(&dagConfig{
		URL:   server.URL,
		DagID: "example-dag",
	}, in)
	if err != nil {
		t.Fatal(err)
	}

	if payload["dag_id"] != "example-dag" {
		t.Fatalf("unexpected dag id: %v", payload["dag_id"])
	}
	conf, ok := payload["conf"].(map[string]interface{})
	if !ok || conf["filename"] != in || conf["content"] != "dag-payload" {
		t.Fatalf("unexpected conf payload: %#v", payload["conf"])
	}
}

func TestRAGPerformMissingContextFile(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	if err := os.WriteFile(in, []byte("question"), 0600); err != nil {
		t.Fatal(err)
	}

	rag := RAG{}
	err := rag.Perform(&ragConfig{
		URL:          "http://localhost:1234/v1/chat/completions",
		Model:        "model",
		Prompt:       "answer",
		ContextFiles: []string{dir + "/missing.txt"},
	}, in)
	if err == nil {
		t.Fatal("expected error for missing context file")
	}
}

func TestDAGPerformMaxPayloadBytes(t *testing.T) {
	dir := t.TempDir()
	in := dir + "/in.txt"
	if err := os.WriteFile(in, []byte("too-large"), 0600); err != nil {
		t.Fatal(err)
	}

	dag := DAG{}
	err := dag.Perform(&dagConfig{
		URL:             "http://example.invalid",
		DagID:           "example-dag",
		MaxPayloadBytes: 3,
	}, in)
	if err == nil || !strings.Contains(err.Error(), "maxPayloadBytes") {
		t.Fatalf("expected maxPayloadBytes error, got: %v", err)
	}
}

package fsagent

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

var safeID = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
var safeHex = regexp.MustCompile(`^[0-9a-f]+$`)

type WorkflowStore interface {
	Save(WorkflowSpec) error
	Get(id string) (WorkflowSpec, error)
	List() ([]string, error)
}

type SQLiteStore struct {
	Path string
}

func (s SQLiteStore) ensure() error {
	if s.Path == "" {
		return fmt.Errorf("sqlite path is empty")
	}
	_, err := exec.LookPath("sqlite3")
	if err != nil {
		return fmt.Errorf("sqlite3 binary not found: %v", err)
	}
	_, err = exec.Command("sqlite3", s.Path, `CREATE TABLE IF NOT EXISTS workflows (id TEXT PRIMARY KEY, name TEXT, graph_hex TEXT NOT NULL);`).CombinedOutput()
	return err
}

func (s SQLiteStore) Save(spec WorkflowSpec) error {
	if err := s.ensure(); err != nil {
		return err
	}
	if !safeID.MatchString(spec.ID) {
		return fmt.Errorf("workflow id contains unsafe characters")
	}
	graph, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	graphHex := hex.EncodeToString(graph)
	if !safeHex.MatchString(graphHex) {
		return fmt.Errorf("invalid graph encoding")
	}
	sql := fmt.Sprintf(
		"INSERT OR REPLACE INTO workflows (id, name, graph_hex) VALUES (%s,%s,%s);",
		sqlQuote(spec.ID), sqlQuote(spec.Name), sqlQuote(graphHex),
	)
	_, err = exec.Command("sqlite3", s.Path, sql).CombinedOutput()
	return err
}

func (s SQLiteStore) Get(id string) (WorkflowSpec, error) {
	if err := s.ensure(); err != nil {
		return WorkflowSpec{}, err
	}
	if !safeID.MatchString(id) {
		return WorkflowSpec{}, fmt.Errorf("workflow id contains unsafe characters")
	}
	out, err := exec.Command("sqlite3", s.Path, "-line", fmt.Sprintf("SELECT graph_hex FROM workflows WHERE id = %s LIMIT 1;", sqlQuote(id))).CombinedOutput()
	if err != nil {
		return WorkflowSpec{}, err
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return WorkflowSpec{}, fmt.Errorf("workflow %q not found", id)
	}
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return WorkflowSpec{}, fmt.Errorf("invalid sqlite response")
	}
	graphHex := strings.TrimSpace(parts[1])
	if !safeHex.MatchString(graphHex) {
		return WorkflowSpec{}, fmt.Errorf("invalid graph encoding")
	}
	b, err := hex.DecodeString(graphHex)
	if err != nil {
		return WorkflowSpec{}, err
	}
	var spec WorkflowSpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return WorkflowSpec{}, err
	}
	return spec, nil
}

func (s SQLiteStore) List() ([]string, error) {
	if err := s.ensure(); err != nil {
		return nil, err
	}
	out, err := exec.Command("sqlite3", s.Path, "SELECT id FROM workflows ORDER BY id ASC;").CombinedOutput()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	list := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			list = append(list, line)
		}
	}
	return list, nil
}

func sqlQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

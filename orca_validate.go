package fsagent

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func ValidateWorkflow(spec WorkflowSpec) error {
	if spec.ID == "" {
		return fmt.Errorf("workflow id is required")
	}
	if len(spec.Nodes) == 0 {
		return fmt.Errorf("workflow must contain at least one node")
	}

	seen := map[string]WorkflowNode{}
	for _, n := range spec.Nodes {
		if n.ID == "" {
			return fmt.Errorf("node id is required")
		}
		if n.Tool == "" {
			return fmt.Errorf("tool is required for node %q", n.ID)
		}
		if _, ok := seen[n.ID]; ok {
			return fmt.Errorf("duplicate node id %q", n.ID)
		}
		if err := validateToolPolicy(n); err != nil {
			return fmt.Errorf("node %q invalid: %v", n.ID, err)
		}
		seen[n.ID] = n
	}

	for _, n := range spec.Nodes {
		for _, d := range n.DependsOn {
			if _, ok := seen[d]; !ok {
				return fmt.Errorf("node %q depends on unknown node %q", n.ID, d)
			}
		}
	}

	if _, err := topologicalSort(spec.Nodes); err != nil {
		return err
	}

	return nil
}

func validateToolPolicy(n WorkflowNode) error {
	switch strings.ToLower(n.Tool) {
	case "http", "webhook":
		var c struct {
			Path string `json:"path"`
			URL  string `json:"url"`
		}
		if err := json.Unmarshal(n.Config, &c); err != nil {
			return fmt.Errorf("invalid config JSON: %v", err)
		}
		u := c.URL
		if u == "" {
			u = c.Path
		}
		if u != "" {
			pu, err := url.Parse(u)
			if err != nil {
				return fmt.Errorf("invalid url: %v", err)
			}
			if pu.Scheme != "https" && !strings.HasPrefix(u, "http://localhost") && !strings.HasPrefix(u, "http://127.0.0.1") {
				return fmt.Errorf("only https and local http endpoints are allowed")
			}
		}
	case "sql", "database", "query":
		var c struct {
			Query string `json:"query"`
			SQL   string `json:"sql"`
		}
		if err := json.Unmarshal(n.Config, &c); err != nil {
			return fmt.Errorf("invalid config JSON: %v", err)
		}
		q := strings.TrimSpace(strings.ToLower(c.Query + " " + c.SQL))
		if q != "" && !strings.HasPrefix(q, "select ") {
			return fmt.Errorf("only read-only SQL statements are allowed")
		}
		if strings.Contains(q, ";") || strings.Contains(q, "--") || strings.Contains(q, "/*") {
			return fmt.Errorf("multi-statement and commented SQL is not allowed")
		}
	case "move", "copy", "delete", "decompress", "compress":
		var cfg map[string]interface{}
		if err := json.Unmarshal(n.Config, &cfg); err != nil {
			return fmt.Errorf("invalid config JSON: %v", err)
		}
		if _, ok := cfg["name"]; !ok && n.Tool != "delete" {
			return fmt.Errorf("structured write tools require config.name")
		}
	}
	return nil
}

func topologicalSort(nodes []WorkflowNode) ([]WorkflowNode, error) {
	graph := map[string][]string{}
	in := map[string]int{}
	idx := map[string]WorkflowNode{}

	for _, n := range nodes {
		in[n.ID] = in[n.ID]
		idx[n.ID] = n
	}
	for _, n := range nodes {
		for _, d := range n.DependsOn {
			graph[d] = append(graph[d], n.ID)
			in[n.ID]++
		}
	}

	queue := make([]string, 0)
	for id, deg := range in {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	out := make([]WorkflowNode, 0, len(nodes))
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		out = append(out, idx[id])
		for _, next := range graph[id] {
			in[next]--
			if in[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if len(out) != len(nodes) {
		return nil, fmt.Errorf("workflow contains a cycle")
	}
	return out, nil
}

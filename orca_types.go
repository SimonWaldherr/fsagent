package fsagent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type WorkflowSpec struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Nodes []WorkflowNode `json:"nodes"`
}

type WorkflowNode struct {
	ID        string          `json:"id"`
	Tool      string          `json:"tool"`
	Config    json.RawMessage `json:"config"`
	DependsOn []string        `json:"dependsOn"`
}

type MachineAction struct {
	Do        string          `json:"do"`
	Config    json.RawMessage `json:"config"`
	Onsuccess []MachineAction `json:"onSuccess"`
	Onfailure []MachineAction `json:"onFailure"`
}

func (w WorkflowSpec) ToMermaid() string {
	var b strings.Builder
	b.WriteString("graph TD\n")
	for _, n := range w.Nodes {
		fmt.Fprintf(&b, "  %s[%s]\n", safeMermaidID(n.ID), n.Tool)
		for _, d := range n.DependsOn {
			fmt.Fprintf(&b, "  %s --> %s\n", safeMermaidID(d), safeMermaidID(n.ID))
		}
	}
	return b.String()
}

func safeMermaidID(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func (w WorkflowSpec) ToMachineActions() ([]MachineAction, error) {
	order, err := topologicalSort(w.Nodes)
	if err != nil {
		return nil, err
	}
	out := make([]MachineAction, 0, len(order))
	for _, n := range order {
		out = append(out, MachineAction{
			Do:     n.Tool,
			Config: n.Config,
		})
	}
	return out, nil
}

func MachineActionsToWorkflow(id, name string, act []MachineAction) WorkflowSpec {
	nodes := make([]WorkflowNode, 0, len(act))
	for i, a := range act {
		nodeID := fmt.Sprintf("n%v", i+1)
		node := WorkflowNode{
			ID:     nodeID,
			Tool:   a.Do,
			Config: a.Config,
		}
		if i > 0 {
			node.DependsOn = []string{fmt.Sprintf("n%v", i)}
		}
		nodes = append(nodes, node)
	}
	return WorkflowSpec{
		ID:    id,
		Name:  name,
		Nodes: nodes,
	}
}

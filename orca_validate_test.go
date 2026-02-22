package fsagent

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestValidateWorkflowCycle(t *testing.T) {
	spec := WorkflowSpec{
		ID: "wf1",
		Nodes: []WorkflowNode{
			{ID: "a", Tool: "sleep", Config: json.RawMessage(`{"time":1}`), DependsOn: []string{"b"}},
			{ID: "b", Tool: "sleep", Config: json.RawMessage(`{"time":1}`), DependsOn: []string{"a"}},
		},
	}
	if err := ValidateWorkflow(spec); err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestValidateWorkflowPolicyHTTP(t *testing.T) {
	spec := WorkflowSpec{
		ID: "wf2",
		Nodes: []WorkflowNode{
			{ID: "a", Tool: "webhook", Config: json.RawMessage(`{"url":"http://example.com"}`)},
		},
	}
	if err := ValidateWorkflow(spec); err == nil {
		t.Fatalf("expected unsafe http policy error")
	}
}

func TestWorkflowExports(t *testing.T) {
	spec := WorkflowSpec{
		ID: "wf3",
		Nodes: []WorkflowNode{
			{ID: "a", Tool: "sleep", Config: json.RawMessage(`{"time":1}`)},
			{ID: "b", Tool: "move", Config: json.RawMessage(`{"name":"x"}`), DependsOn: []string{"a"}},
		},
	}

	mermaid := spec.ToMermaid()
	if !strings.Contains(mermaid, "a --> b") {
		t.Fatalf("unexpected mermaid: %s", mermaid)
	}

	actions, err := spec.ToMachineActions()
	if err != nil {
		t.Fatal(err)
	}
	roundTrip := MachineActionsToWorkflow("wf3", "name", actions)
	if len(roundTrip.Nodes) != 2 {
		t.Fatalf("unexpected round-trip nodes: %#v", roundTrip.Nodes)
	}
}

package fsagent

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

//go:embed cmd/fsagent/web/orca.html
var orcaAssets embed.FS

var orcaOnce sync.Once

func registerOrcaHandlers(conf Config) {
	orcaOnce.Do(func() {
		dbPath := conf.WorkflowDB
		if dbPath == "" {
			dbPath = "orca.sqlite"
		}
		store := SQLiteStore{Path: dbPath}

		http.HandleFunc("/orca", func(w http.ResponseWriter, r *http.Request) {
			b, err := orcaAssets.ReadFile("cmd/fsagent/web/orca.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			if _, err := w.Write(b); err != nil {
				log.Printf("could not write /orca response: %v", err)
			}
		})

		http.HandleFunc("/api/orca/workflows", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				ids, err := store.List()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(map[string]interface{}{"ids": ids}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case http.MethodPost:
				var spec WorkflowSpec
				if err := json.NewDecoder(r.Body).Decode(&spec); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err := ValidateWorkflow(spec); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				if err := store.Save(spec); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusCreated)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})

		http.HandleFunc("/api/orca/convert/machine-to-graph", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			var req struct {
				ID      string          `json:"id"`
				Name    string          `json:"name"`
				Actions []MachineAction `json:"actions"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			spec := MachineActionsToWorkflow(req.ID, req.Name, req.Actions)
			if err := ValidateWorkflow(spec); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(spec); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		http.HandleFunc("/api/orca/workflows/", func(w http.ResponseWriter, r *http.Request) {
			id := strings.TrimPrefix(r.URL.Path, "/api/orca/workflows/")
			if id == "" {
				http.Error(w, "workflow id missing", http.StatusBadRequest)
				return
			}
			spec, err := store.Get(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			switch r.URL.Query().Get("format") {
			case "machine":
				machine, err := spec.ToMachineActions()
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(machine); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "mermaid":
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				if _, err := fmt.Fprint(w, spec.ToMermaid()); err != nil {
					log.Printf("could not write mermaid response: %v", err)
					return
				}
			default:
				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(spec); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
	})
}

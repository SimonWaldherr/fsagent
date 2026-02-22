package fsagent

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func serveInboundHTTP(conf Config, fnChannel chan<- string, source string) {
	path := conf.Folder
	if path == "" {
		path = "/" + source
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tmp, err := os.CreateTemp("", fmt.Sprintf("fsagent-%s-*.txt", source))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tmp.Close()
		if _, err := tmp.Write(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fnChannel <- tmp.Name()
		fmt.Fprint(w, strings.ToUpper(source)+" event accepted")
	})

	if err := http.ListenAndServe(conf.Port, nil); err != nil {
		log.Printf("inbound %s server stopped: %v", source, err)
	}
}

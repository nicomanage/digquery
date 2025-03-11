package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type DigResult struct {
	Answer []string `json:"answer"`
	Error  string   `json:"error,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	if domain == "" {
		http.Error(w, "Missing 'domain' parameter", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("dig", "+short", domain, "A")
	out, err := cmd.Output()
	if err != nil {
		//Handle errors robustly.  Consider logging, more specific error messages, etc.
		fmt.Fprintf(w, "{\"error\":\"%s\"}", err.Error())
		return
	}

	result := DigResult{Answer: strings.Split(string(out), "\n"), Error: ""}
	//Clean up empty lines in the result
	var cleanedAnswer []string
	for _, line := range result.Answer {
		if line != "" {
			cleanedAnswer = append(cleanedAnswer, line)
		}
	}
	result.Answer = cleanedAnswer

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

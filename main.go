package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/istsh/markdown-viewer/lexer"
	"github.com/istsh/markdown-viewer/parser"
)

type Input struct {
	Markdown string `json:"markdown"`
}

func main() {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(http.StatusOK)
		}
	})
	http.HandleFunc("/parse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}

			input := &Input{}
			if err := json.Unmarshal(body, input); err != nil {
				panic(err)
			}

			inputBytes := []byte(input.Markdown)

			l := lexer.New(inputBytes)
			p := parser.New(l)
			result := p.Parse()

			w.Header().Set("Content-Type", "application/json")
			w.Write(result)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}

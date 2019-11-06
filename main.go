package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/istsh/markdown-viewer/token"

	"github.com/istsh/markdown-viewer/lexer"
)

type Input struct {
	Markdown string `json:"markdown"`
}

// See: https://qiita.com/tbpgr/items/989c6badefff69377da7
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

			var res []byte

			l := lexer.New(inputBytes)
			for {
				tok := l.NextToken()
				if tok.Type == token.EOF {
					break
				}

				if tok.Type == token.LINE_FEED_CODE {
					res = append(res, '\n')
					//fmt.Println()
				} else if tok.Type == token.STRING {
					res = append(res, tok.Literal...)
					//fmt.Printf("%s", tok.Literal)
				} else {
					res = append(res, []byte(fmt.Sprintf("%q", tok.Type))...)
					//fmt.Printf("%q", tok.Type)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}

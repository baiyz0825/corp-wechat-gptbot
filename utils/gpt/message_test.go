package gpt

import (
	"fmt"
	"testing"
)

func TestTransMarkdown(t *testing.T) {
	str := "\n\npackage main\n\nimport (\n    \"fmt\"\n    \"net/http\"\n)\n\nfunc handler(w http.ResponseWriter, r *http.Request) {\n    fmt.Fprintf(w, \"Hi there, I love %s!\", r.URL.Path[1:])\n}\n\nfunc main() {\n    http.HandleFunc(\"/\", handler)\n    http.ListenAndServe(\":8080\", nil)\n}"

	fmt.Printf("```\n%s\n```", str)
}

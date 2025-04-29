package main

import (
	"github.com/robmerrell/exarch/internal/search"
)

// ../../../queries/func_def.scm
// var funcDefQuery string

const funcDefQuery = `(
(call target: (identifier) @keyword
  (arguments
    [(call target: (identifier))
     (identifier)] @func_name)) @identifier
(#match? @keyword "^(def|defp)$")
(#match? @func_name "follow")
)`

func main() {
	file := "/home/rob/projects/gh/elixir-phoenix-realworld-example-app/lib/real_world/accounts/users.ex"
	// funcinput := "follow"
	input := &search.SearchInput{
		SearchType:  search.SearchTypeStr,
		SearchTerms: "",
		// Function:    &funcinput,
	}

	search.SearchFile(file, input)
}

// func buildQuery(input *search.SearchInput) error {
// 	// if we are searching in functions then use the func def query

// 	// parse the template
// 	tpl, err := template.New("func_def").Parse(funcDefQuery)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

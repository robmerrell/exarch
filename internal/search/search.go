package search

import (
	"context"
	_ "embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

// This is a single purpose searcher with no plans for supporting other languages, so
// just define elixir at the package level.
var lang = elixir.GetLanguage()

// alias filename to string to make the search return type clearer.
type filename string

// SearchType and the "enum" below are used to constrain the search to a specific
// node type.
type SearchType int

const (
	SearchTypeAlias SearchType = iota
	SearchTypeStr
	SearchTypeAtom
)

// SearchInput holds all of the input necessary to perform a search. The only
// optional input is Function.
type SearchInput struct {
	SearchTerms string
	SearchType  SearchType

	Function *string // Search within specific matching functions
}

// Match represents a match found in the elixir source code.
type Match struct {
	Row      uint32 // The file row the matched node begins on
	Contents string // The contents of the matched node
}

// working
const funcDefQuery = `(
(call target: (identifier) @keyword
  (arguments
    [(call target: (identifier))
     (identifier)] @func_name)) @identifier
(#match? @keyword "^(def|defp)$")
(#match? @func_name "follow")
)`

//go:embed queries/alias_search.scm
var aliasSearchQuery string

//go:embed queries/string_search.scm
var strSearchQuery string

//go:embed queries/atom_search.scm
var atomSearchQuery string

func Search() error {
	// get all files to search
	filepath.WalkDir("/home/rob/projects/gh/elixir-phoenix-realworld-example-app", func(_path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// check for elixir files
		if entry.Type().IsRegular() {
			ext := filepath.Ext(entry.Name())
			if ext == ".ex" || ext == ".exs" {
				// read the files and check if they match any searches

				// put the matches in our list

				// if in-mod is given see if the file has a matching module
				// if in-fn is given see if the file has a matching function
				// fmt.Println(entry.Name())
			}

		}

		return nil
	})

	// go through each file and apply the filters first. If they are empty move on

	// find the nodes
	return nil
}

// searches in a specific file
func searchFile(file string, input *SearchInput) ([]Match, error) {
	// setup the parser
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	// read the file
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// get the root node to start searching from
	root, err := sitter.ParseCtx(context.Background(), contents, lang)
	if err != nil {
		return nil, err
	}

	// if searches are constrained by a function we perform the search on each node
	// of the found function. Otherwise perform the search from the root node.
	if input.Function == nil {
		return performSearch(root, contents, input)
	} else {
		// get all of the function nodes and perform a search on each

		return nil, err
	}

	// query for specific functions
	// if input.Function != nil {
	// 	query, err := sitter.NewQuery([]byte(funcDefQuery), elixir.GetLanguage())
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	cursor := sitter.NewQueryCursor()
	// 	cursor.Exec(query, root)

	// 	// results
	// 	for {
	// 		// get the match and break out if we're done matching
	// 		match, ok := cursor.NextMatch()
	// 		if !ok {
	// 			break
	// 		}

	// 		match = cursor.FilterPredicates(match, contents)
	// 		for _, capture := range match.Captures {
	// 			if capture.Node.Type() == "call" {

	// 				// look at the first child to tell if we're at the top level def/defn
	// 				if first := capture.Node.Child(0); first != nil {
	// 					if first.Content(contents) == "def" || first.Content(contents) == "defp" {
	// 						fmt.Println(capture.Node.String())
	// 						fmt.Println(capture.Node.Content(contents))
	// 						fmt.Println("=====")
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}

	// }
}

func performSearch(node *sitter.Node, contents []byte, input *SearchInput) ([]Match, error) {
	var matches []Match

	var tsQuery string
	var processFunc func(capture sitter.QueryCapture, contents []byte) *Match

	switch input.SearchType {
	case SearchTypeAlias:
		tsQuery = fmt.Sprintf(aliasSearchQuery, input.SearchTerms)
		processFunc = processStrCapture
	case SearchTypeAtom:
		tsQuery = fmt.Sprintf(atomSearchQuery, input.SearchTerms)
		processFunc = processStrCapture
	case SearchTypeStr:
		tsQuery = fmt.Sprintf(strSearchQuery, input.SearchTerms)
		processFunc = processStrCapture
	default:
		return matches, fmt.Errorf("Invalid search type: %d", input.SearchType)
	}

	query, err := sitter.NewQuery([]byte(tsQuery), lang)
	if err != nil {
		return matches, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, node)

	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		match = cursor.FilterPredicates(match, contents)
		for _, capture := range match.Captures {
			// run the process function on the capture and ignore any that return nil
			if processedMatch := processFunc(capture, contents); processedMatch != nil {
				matches = append(matches, *processedMatch)
			}
		}
	}

	return matches, nil
}

// process the captures for string searches and filter out module and doc comments.
func processStrCapture(capture sitter.QueryCapture, contents []byte) *Match {
	// follow the node up 2 parents. If the grandparent node is a doc or moduledoc
	// identifier then don't return a match.
	isDoc := false

	if parent := capture.Node.Parent(); parent != nil {
		if grandparent := parent.Parent(); grandparent != nil {
			target := grandparent.ChildByFieldName("target")
			if target != nil && (target.Content(contents) == "doc" || target.Content(contents) == "moduledoc") {
				isDoc = true
			}
		}
	}

	if isDoc {
		return nil
	} else {
		return &Match{
			Row:      capture.Node.StartPoint().Row + 1,
			Contents: capture.Node.Content(contents),
		}
	}

}

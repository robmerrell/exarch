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

// alias filename to string to make the search return type clearer.
type filename string

// SearchType and the "enum" below are used to constrain the search to a specific
// node type.
type SearchType int

const (
	SearchTypeStr SearchType = iota
	SearchTypeDoc
	SearchTypeFnCall
)

// SearchInput holds all of the input necessary to perform a search. The only
// optional input is Function.
type SearchInput struct {
	SearchTerms string
	SearchType  SearchType
	Dir         string
}

// Match represents a match found in the elixir source code.
type Match struct {
	Row      uint32 // The file row the matched node begins on
	Contents string // The contents of the matched node
}

type ResultsFormatter interface {
	Format() string
}

// Search performs a search and prints results to stdout
func Search(input *SearchInput) error {
	// get all files to search
	return filepath.WalkDir(input.Dir, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// check for elixir files
		if entry.Type().IsRegular() {
			ext := filepath.Ext(entry.Name())
			if ext == ".ex" || ext == ".exs" {
				res, err := searchFile(path, input)
				if err != nil {
					return err
				}

				relFile, err := filepath.Rel(input.Dir, path)
				if err != nil {
					return err
				}

				if len(res) > 0 {
					fmt.Println(relFile)
					for _, match := range res {
						fmt.Println(match)
					}
					fmt.Println("")
				}
			}
		}

		return nil
	})
}

func searchFile(file string, input *SearchInput) ([]string, error) {
	lang := elixir.GetLanguage()

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

	var searchResults []ResultsFormatter
	var searchErr error

	switch input.SearchType {
	case SearchTypeStr:
	case SearchTypeFnCall:
		searchResults, searchErr = searchFnCalls(root, contents, input)
	default:
		return nil, fmt.Errorf("Invalid search type: %d", input.SearchType)
	}

	if searchErr != nil {
		return nil, err
	}

	output := []string{}
	for _, match := range searchResults {
		output = append(output, match.Format())
	}

	return output, nil
}

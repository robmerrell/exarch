package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/robmerrell/exarch/internal/search"
	"github.com/urfave/cli/v3"
)

const desc = `Searches recursively from the current directory.

SEARCH_MODE can be one of:
1. fncall - Search for function calls. This searches partial matches,
            and is able to handle aliases if given a fully qualified
            function. Eg. TestApp.Users.process can search for
            Users.process if TestApp.Users has been aliased.
2. str - Search inside of strings.
3. doc - search inside of documentation.`

func main() {
	var searchMode string
	var searchTerms string

	cmd := &cli.Command{
		Name:        "exarch",
		Usage:       "Semantic Elixir Search",
		ArgsUsage:   "SEARCH_MODE SEARCH",
		Description: desc,
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:        "search_mode",
				Destination: &searchMode,
			},
			&cli.StringArg{
				Name:        "search_terms",
				Destination: &searchTerms,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// make sure the search mode is valid
			var searchType search.SearchType
			switch searchMode {
			case "fncall":
				searchType = search.SearchTypeFnCall
			case "str":
				searchType = search.SearchTypeStr
			case "doc":
				searchType = search.SearchTypeDoc
			default:
				return cli.Exit("Invalid SEARCH_MODE, use --help for instructions", 1)
			}

			if searchTerms == "" {
				return cli.Exit("Can't use empty search terms, use --help for instructions", 1)
			}

			input, err := buildInput(searchType, searchTerms)
			if err != nil {
				return cli.Exit(fmt.Sprintf("Input Error: %v", err), 1)
			}

			search.Search(input)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func buildInput(searchType search.SearchType, searchTerms string) (*search.SearchInput, error) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	input := &search.SearchInput{
		SearchType:  searchType,
		SearchTerms: searchTerms,
		Dir:         dir,
	}

	return input, nil
}

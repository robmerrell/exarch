package search

import (
	_ "embed"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

//go:embed queries/remote_fn_call.scm
var remoteFnCallQuery string

func parseAliases() {}

// Generate a list of all remote functions like Module.fn_call()
func parseRemoteCalls(root *sitter.Node) ([]*sitter.Node, error) {
	query, err := sitter.NewQuery([]byte(remoteFnCallQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	matches := []*sitter.Node{}
	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		// match = cursor.FilterPredicates(match, contents)
		for _, capture := range match.Captures {
			matches = append(matches, capture.Node)
		}
	}

	return matches, nil
}

func parseLocalCalls() {}

func search() {

}

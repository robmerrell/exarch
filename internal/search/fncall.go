package search

import (
	_ "embed"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

//go:embed queries/alias_search.scm
var aliasQuery string

//go:embed queries/remote_fn_call.scm
var remoteFnCallQuery string

// Generate a list of all module aliases. Any aliases that are group together in a tuple
// like Module.{Sub1, Sub2} are separated into multiple entries.
func parseAliases(root *sitter.Node, contents []byte) ([]string, error) {
	query, err := sitter.NewQuery([]byte(aliasQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	aliases := []string{}
	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		match = cursor.FilterPredicates(match, contents)
		for _, capture := range match.Captures {
			if capture.Node.Type() == "arguments" {
				if child := capture.Node.Child(0); child != nil {
					switch child.Type() {
					// single alias
					case "alias":
						aliases = append(aliases, capture.Node.Content(contents))
					// multiple aliases
					case "dot":
						modulePrefix := child.ChildByFieldName("left").Content(contents)

						aliasNodes := child.ChildByFieldName("right")
						for i := range int(aliasNodes.ChildCount()) {
							if aliasNodes.Child(i).Type() == "alias" {
								alias := fmt.Sprintf("%s.%s", modulePrefix, aliasNodes.Child(i).Content(contents))
								aliases = append(aliases, alias)
							}

						}
					}
				}

			}
		}
	}

	return aliases, nil

}

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

		for _, capture := range match.Captures {
			matches = append(matches, capture.Node)
		}
	}

	return matches, nil
}

func parseLocalCalls() {}

func search() {

}

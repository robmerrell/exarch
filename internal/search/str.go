package search

import (
	_ "embed"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

type Str struct {
	Line     uint32
	Contents string
}

func (s Str) Format() string {
	return fmt.Sprintf("%d:%s", s.Line, s.Contents)
}

//go:embed queries/string_search.scm
var strSearchQuery string

func searchStr(root *sitter.Node, contents []byte, input *SearchInput) ([]ResultsFormatter, error) {
	query, err := sitter.NewQuery([]byte(strSearchQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	matches := []ResultsFormatter{}
	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		match = cursor.FilterPredicates(match, contents)
		for _, capture := range match.Captures {
			lines := strings.Split(capture.Node.Content(contents), "\n")
			for i, line := range lines {
				if !isDoc(capture.Node, contents) && strings.Contains(line, input.SearchTerms) {
					matches = append(matches, Str{
						Contents: line,
						Line:     capture.Node.StartPoint().Row + uint32(i),
					})
				}
			}
		}
	}

	return matches, nil
}

func searchDoc(root *sitter.Node, contents []byte, input *SearchInput) ([]ResultsFormatter, error) {
	query, err := sitter.NewQuery([]byte(strSearchQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	matches := []ResultsFormatter{}
	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		match = cursor.FilterPredicates(match, contents)
		for _, capture := range match.Captures {
			lines := strings.Split(capture.Node.Content(contents), "\n")
			for i, line := range lines {
				if isDoc(capture.Node, contents) && strings.Contains(line, input.SearchTerms) {
					matches = append(matches, Str{
						Contents: line,
						Line:     capture.Node.StartPoint().Row + uint32(i),
					})
				}
			}
		}
	}

	return matches, nil
}

// checks if the string is part of a documentation block
func isDoc(node *sitter.Node, contents []byte) bool {
	// follow the node up 2 parents. If the grandparent node is a doc or moduledoc
	// identifier then don't return a match.
	isDoc := false

	if parent := node.Parent(); parent != nil {
		if grandparent := parent.Parent(); grandparent != nil {
			target := grandparent.ChildByFieldName("target")
			if target != nil && (target.Content(contents) == "doc" || target.Content(contents) == "moduledoc") {
				isDoc = true
			}
		}
	}

	return isDoc
}

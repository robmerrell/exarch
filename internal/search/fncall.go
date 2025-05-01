package search

import (
	_ "embed"
	"fmt"
	"slices"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

type FnCall struct {
	ModulePath string
	Name       string
	Line       uint32
	Contents   string
}

func (f FnCall) Format() string {
	return fmt.Sprintf("%d:%s", f.Line, f.Contents)
}

type Alias struct {
	ModulePath string
	As         string
	Line       uint32
	Contents   string
}

//go:embed queries/alias_search.scm
var aliasQuery string

//go:embed queries/remote_fn_call.scm
var remoteFnCallQuery string

// Generate a list of all module aliases. Any aliases that are group together in a tuple
// like Module.{Sub1, Sub2} are separated into multiple entries.
func parseAliases(root *sitter.Node, contents []byte) ([]Alias, error) {
	query, err := sitter.NewQuery([]byte(aliasQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	aliases := []Alias{}
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
						aliases = append(aliases, singleAlias(capture.Node, contents))
					// multiple aliases
					case "dot":
						aliases = append(aliases, multipleAliases(child, contents)...)
					}
				}
			}
		}
	}

	return aliases, nil
}

func singleAlias(node *sitter.Node, contents []byte) Alias {
	// TODO: has an 'as' key
	modulePath := node.Content(contents)
	aliasAs := modulePath[strings.LastIndex(modulePath, ".")+1:]

	return Alias{
		ModulePath: modulePath,
		As:         aliasAs,
		Contents:   node.Parent().Content(contents),
		Line:       node.StartPoint().Row + 1,
	}
}

func multipleAliases(node *sitter.Node, contents []byte) []Alias {
	modulePrefix := node.ChildByFieldName("left").Content(contents)
	aliasNodes := node.ChildByFieldName("right")

	aliases := []Alias{}
	for i := range int(aliasNodes.ChildCount()) {
		if aliasNodes.Child(i).Type() == "alias" {
			as := aliasNodes.Child(i).Content(contents)
			aliases = append(aliases, Alias{
				ModulePath: fmt.Sprintf("%s.%s", modulePrefix, as),
				As:         as,
				Contents:   node.Parent().Parent().Content(contents),
				Line:       node.StartPoint().Row + 1,
			})
		}
	}

	return aliases
}

// Generate a list of all remote functions like Module.fn_call()
func parseRemoteCalls(root *sitter.Node, contents []byte, aliases []Alias) ([]FnCall, error) {
	query, err := sitter.NewQuery([]byte(remoteFnCallQuery), elixir.GetLanguage())
	if err != nil {
		return nil, err
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	functions := []FnCall{}
	for {
		// get the match and break out if we're done matching
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			if child := capture.Node.ChildByFieldName("target"); child != nil {
				modulePrefix := child.ChildByFieldName("left").Content(contents)
				fnName := child.ChildByFieldName("right").Content(contents)
				modulePath := findFullModulePath(modulePrefix, aliases)

				functions = append(functions, FnCall{
					ModulePath: modulePath,
					Name:       fnName,
					Contents:   capture.Node.Content(contents),
					Line:       capture.Node.StartPoint().Row + 1,
				})
			}
		}
	}

	return functions, nil
}

// find the entire module path just given a prefix and a list of aliases in the module. Do this by
// seeing if the prefix matches and As: field in the aliases. If not match from right to left.
func findFullModulePath(modulePrefix string, aliases []Alias) string {
	// matching as
	for _, alias := range aliases {
		if alias.As == modulePrefix {
			return alias.ModulePath
		}
	}

	// matching right to left
	for _, alias := range aliases {
		splitAliasModulePath := strings.Split(alias.ModulePath, ".")
		slices.Reverse(splitAliasModulePath)
		splitModulePrefix := strings.Split(modulePrefix, ".")
		slices.Reverse(splitModulePrefix)

		if len(splitAliasModulePath) >= len(splitModulePrefix) {
			matches := true
			for i := range len(splitModulePrefix) {
				matches = matches && (splitAliasModulePath[i] == splitModulePrefix[i])
			}

			if matches {
				return alias.ModulePath
			}
		}
	}

	return modulePrefix
}

func parseLocalCalls() {}

func searchFnCalls(root *sitter.Node, contents []byte, input *SearchInput) ([]ResultsFormatter, error) {
	aliases, err := parseAliases(root, contents)
	if err != nil {
		return nil, err
	}

	fnCalls, err := parseRemoteCalls(root, contents, aliases)
	if err != nil {
		return nil, err
	}

	matching := []ResultsFormatter{}
	for _, fn := range fnCalls {
		fullFnCall := fmt.Sprintf("%s.%s", fn.ModulePath, fn.Name)
		if strings.Contains(fullFnCall, input.SearchTerms) {
			matching = append(matching, fn)
		}
	}

	return matching, nil
}

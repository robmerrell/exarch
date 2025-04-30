package search

import (
	"context"
	"os"
	"reflect"
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/elixir"
)

func readTestFile(t *testing.T) (*sitter.Node, []byte) {
	lang = elixir.GetLanguage()

	// setup the parser
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	// read the file
	contents, err := os.ReadFile("testdata/users.ex")
	if err != nil {
		t.Errorf("Unable to read file, %v", err)
	}

	// get the root node to start searching from
	root, err := sitter.ParseCtx(context.Background(), contents, lang)
	if err != nil {
		t.Errorf("Unable to parse, %v", err)
	}

	return root, contents
}

func TestParseRemoteCalls(t *testing.T) {
	root, contents := readTestFile(t)
	nodes, err := parseRemoteCalls(root)
	if err != nil {
		t.Errorf("parse remote calls failed: %v", err)
	}

	fns := []string{}
	for _, node := range nodes {
		fns = append(fns, node.Content(contents))
	}

	expected := []string{
		"Repo.get!(User, id)",
		"Repo.get_by(User, username: username)",
		"TestApp.FilterChain.process()",
		"User.changeset(attrs)",
		"Repo.update()",
	}

	if !reflect.DeepEqual(fns, expected) {
		t.Errorf("got %v want %v", fns, expected)
	}
}

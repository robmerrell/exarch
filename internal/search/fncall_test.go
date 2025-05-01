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
	lang := elixir.GetLanguage()

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

func TestParseAliases(t *testing.T) {
	root, contents := readTestFile(t)
	aliases, err := parseAliases(root, contents)
	if err != nil {
		t.Errorf("parse aliases failed: %v", err)
	}

	expected := []Alias{
		{ModulePath: "TestApp.Repo", As: "Repo", Contents: "alias TestApp.Repo", Line: 6},
		{ModulePath: "TestApp.Accounts.User", As: "User", Contents: "alias TestApp.Accounts.{User, Admin}", Line: 7},
		{ModulePath: "TestApp.Accounts.Admin", As: "Admin", Contents: "alias TestApp.Accounts.{User, Admin}", Line: 7},
	}

	if !reflect.DeepEqual(aliases, expected) {
		t.Errorf("got %v want %v", aliases, expected)
	}
}

func TestParseRemoteCalls(t *testing.T) {
	root, contents := readTestFile(t)

	aliases, err := parseAliases(root, contents)
	if err != nil {
		t.Errorf("parse aliases failed: %v", err)
	}

	fnCalls, err := parseRemoteCalls(root, contents, aliases)
	if err != nil {
		t.Errorf("parse remote calls failed: %v", err)
	}

	expected := []FnCall{
		{ModulePath: "TestApp.Repo", Name: "get!", Contents: "Repo.get!(User, id)", Line: 12},
		{ModulePath: "TestApp.Repo", Name: "get_by", Contents: "Repo.get_by(User, username: username)", Line: 15},
		{ModulePath: "TestApp.FilterChain", Name: "process", Contents: "TestApp.FilterChain.process()", Line: 20},
		{ModulePath: "TestApp.Accounts.User", Name: "changeset", Contents: "User.changeset(attrs)", Line: 21},
		{ModulePath: "TestApp.Repo", Name: "update", Contents: "Repo.update()", Line: 22},
	}

	if !reflect.DeepEqual(fnCalls, expected) {
		t.Errorf("got %v want %v", fnCalls, expected)
	}
}

func TestFindFullModulePath(t *testing.T) {
	root, contents := readTestFile(t)

	aliases, err := parseAliases(root, contents)
	if err != nil {
		t.Errorf("parse aliases failed: %v", err)
	}

	if module := findFullModulePath("TestApp.Accounts.User", aliases); module != "TestApp.Accounts.User" {
		t.Errorf("got %v want %v", module, "TestApp.Accounts.User")
	}

	if module := findFullModulePath("Accounts.User", aliases); module != "TestApp.Accounts.User" {
		t.Errorf("got %v want %v", module, "TestApp.Accounts.User")
	}

	if module := findFullModulePath("User", aliases); module != "TestApp.Accounts.User" {
		t.Errorf("got %v want %v", module, "TestApp.Accounts.User")
	}

	if module := findFullModulePath("No", aliases); module != "No" {
		t.Errorf("got %v want %v", module, "No")
	}
}

func TestSearchFnCallsFullyQualifiedName(t *testing.T) {
	root, contents := readTestFile(t)

	input := &SearchInput{
		SearchType:  SearchTypeFnCall,
		SearchTerms: "TestApp.Repo.update",
	}

	fnCalls, _ := searchFnCalls(root, contents, input)
	expected := []ResultsFormatter{
		FnCall{ModulePath: "TestApp.Repo", Name: "update", Contents: "Repo.update()", Line: 22},
	}

	if !reflect.DeepEqual(fnCalls, expected) {
		t.Errorf("got %v want %v", fnCalls, expected)
	}
}

func TestSearchFnCallsPartialName(t *testing.T) {
	root, contents := readTestFile(t)

	input := &SearchInput{
		SearchType:  SearchTypeFnCall,
		SearchTerms: "FilterChain.process",
	}

	fnCalls, _ := searchFnCalls(root, contents, input)
	expected := []ResultsFormatter{
		FnCall{ModulePath: "TestApp.FilterChain", Name: "process", Contents: "TestApp.FilterChain.process()", Line: 20},
	}

	if !reflect.DeepEqual(fnCalls, expected) {
		t.Errorf("got %v want %v", fnCalls, expected)
	}
}

func TestSearchFnCallsNoModule(t *testing.T) {
	root, contents := readTestFile(t)

	input := &SearchInput{
		SearchType:  SearchTypeFnCall,
		SearchTerms: "process",
	}

	fnCalls, _ := searchFnCalls(root, contents, input)
	expected := []ResultsFormatter{
		FnCall{ModulePath: "TestApp.FilterChain", Name: "process", Contents: "TestApp.FilterChain.process()", Line: 20},
	}

	if !reflect.DeepEqual(fnCalls, expected) {
		t.Errorf("got %v want %v", fnCalls, expected)
	}
}

func TestSearchFnCallsOnlyModule(t *testing.T) {
	root, contents := readTestFile(t)

	input := &SearchInput{
		SearchType:  SearchTypeFnCall,
		SearchTerms: "Repo",
	}

	fnCalls, _ := searchFnCalls(root, contents, input)
	expected := []ResultsFormatter{
		FnCall{ModulePath: "TestApp.Repo", Name: "get!", Contents: "Repo.get!(User, id)", Line: 12},
		FnCall{ModulePath: "TestApp.Repo", Name: "get_by", Contents: "Repo.get_by(User, username: username)", Line: 15},
		FnCall{ModulePath: "TestApp.Repo", Name: "update", Contents: "Repo.update()", Line: 22},
	}

	if !reflect.DeepEqual(fnCalls, expected) {
		t.Errorf("got %v want %v", fnCalls, expected)
	}
}

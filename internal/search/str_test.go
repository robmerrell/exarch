package search

import (
	"reflect"
	"testing"
)

func TestSearchStr(t *testing.T) {
	root, contents := readTestFile(t)
	input := &SearchInput{
		SearchType:  SearchTypeStr,
		SearchTerms: "string",
	}

	matches, err := searchStr(root, contents, input)
	if err != nil {
		t.Errorf("search string failed: %v", err)
	}

	expected := []ResultsFormatter{
		Str{Contents: "\"string one\"", Line: 41},
		Str{Contents: "\"string two\"", Line: 42},
		Str{Contents: "\"string three\"", Line: 43},
	}

	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("got %v want %v", matches, expected)
	}
}

func TestSearchDoc(t *testing.T) {
	root, contents := readTestFile(t)
	input := &SearchInput{
		SearchType:  SearchTypeDoc,
		SearchTerms: "user",
	}

	matches, err := searchDoc(root, contents, input)
	if err != nil {
		t.Errorf("search doc failed: %v", err)
	}

	expected := []ResultsFormatter{
		Str{Contents: "  Get a user by id", Line: 10},
		Str{Contents: "  Returns a hello message for the user as a string", Line: 26},
	}

	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("got %v want %v", matches, expected)
	}
}

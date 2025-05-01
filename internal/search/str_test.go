package search

import (
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

	expected := []Str{
		{Contents: "\"string one\"", Line: 41},
		{Contents: "\"string two\"", Line: 42},
		{Contents: "\"string three\"", Line: 43},
	}

	if len(matches) != 3 {
		t.Errorf("too many matches")
	}

	for i := range matches {
		if matches[i] != expected[i] {
			t.Errorf("got %v want %v", matches[i], expected[i])
		}
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

	expected := []Str{
		{Contents: "  Get a user by id", Line: 10},
		{Contents: "  Returns a hello message for the user as a string", Line: 26},
	}

	if len(matches) != 2 {
		t.Errorf("too many matches")
	}

	for i := range matches {
		if matches[i] != expected[i] {
			t.Errorf("got %v want %v", matches[i], expected[i])
		}
	}
}

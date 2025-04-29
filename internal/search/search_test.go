package search

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSearchWithString(t *testing.T) {
	input := &SearchInput{
		SearchType:  SearchTypeStr,
		SearchTerms: "string",
	}

	matches, err := searchFile("testdata/users.ex", input)
	if err != nil {
		t.Errorf("searching file failed: %+v", err)
	}

	expected := []Match{
		{Row: 40, Contents: `"string one"`},
		{Row: 41, Contents: `"string two"`},
		{Row: 42, Contents: `"string three"`},
	}

	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("got %v want %v", matches, expected)
	}
}

func TestSearchWithAtom(t *testing.T) {
	input := &SearchInput{
		SearchType:  SearchTypeAtom,
		SearchTerms: "ber",
	}

	matches, err := searchFile("testdata/users.ex", input)
	if err != nil {
		t.Errorf("searching file failed: %+v", err)
	}

	expected := []Match{
		{Row: 30, Contents: ":blueberry"},
		{Row: 31, Contents: ":strawberry"},
	}

	if !reflect.DeepEqual(matches, expected) {
		t.Errorf("got %v want %v", matches, expected)
	}
}

func TestSearchWithAlias(t *testing.T) {
	input := &SearchInput{
		SearchType:  SearchTypeAlias,
		SearchTerms: "",
	}

	matches, err := searchFile("testdata/users.ex", input)
	if err != nil {
		t.Errorf("searching file failed: %+v", err)
	}

	fmt.Println(matches)

	// expected := []Match{
	// 	{Row: 30, Contents: ":blueberry"},
	// 	{Row: 31, Contents: ":strawberry"},
	// }

	// if !reflect.DeepEqual(matches, expected) {
	// 	t.Errorf("got %v want %v", matches, expected)
	// }
}

// func TestSearchWithStringIgnoresComments(t *testing.T) {
// 	input := &SearchInput{
// 		SearchType:  SearchTypeStr,
// 		SearchTerms: "user",
// 	}

// 	matches, err := searchFile("testdata/users.ex", input)
// 	if err != nil {
// 		t.Errorf("searching file failed: %+v", err)
// 	}

// 	fmt.Println(matches)
// }

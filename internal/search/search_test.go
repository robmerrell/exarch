package search

import (
	"reflect"
	"testing"
)

func TestSearchFile(t *testing.T) {
	res, err := searchFile("testdata/users.ex", &SearchInput{
		SearchTerms: "process",
		SearchType:  SearchTypeFnCall,
		Dir:         "",
	})
	if err != nil {
		t.Errorf("failed searchFile: %v", err)
	}

	expected := []string{"20:TestApp.FilterChain.process()"}

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("got %v want %v", res, expected)
	}
}

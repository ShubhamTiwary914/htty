package htty_test

import (
	utils "htty/utils"
	"os"
	"testing"
	"strings"
)

func TestSortedInsertReadAndSearch(tt *testing.T) {
	// temp file
	f, err := os.CreateTemp("", "words-*.txt")
	if err != nil {
		tt.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	f.Close()

	path := f.Name()

	// insert words
	input := []string{
		"Apple",
		"apple",
		"Hel",
		"Hell",
		"Hello",
		"HelloWorld",
		"random",
		"Word",
	}

	for _, w := range input {
		if err := utils.Insert_SortedFile(w, path); err != nil {
			tt.Fatalf("insert failed: %v", err)
		}
	}

	// read and verify sorted order (case-insensitive primary)
	words, err := utils.ReadTextLines_intoList(path)
	if err != nil {
		tt.Fatalf("read failed: %v", err)
	}

	for i := 1; i < len(words); i++ {
		prev := words[i-1]
		curr := words[i]

		if !(strings.ToLower(prev) < strings.ToLower(curr) ||
			(strings.EqualFold(prev, curr) && prev <= curr)) {
			tt.Errorf("not sorted at index %d: %s, %s", i, prev, curr)
		}
	}

	// prefix search
	res, err := utils.PrefixClosestSearch("hello", path)
	if err != nil {
		tt.Fatalf("search failed: %v", err)
	}

	expected := []string{"Hello", "HelloWorld"}

	if len(res) != len(expected) {
		tt.Errorf("expected %d results, got %d", len(expected), len(res))
	}

	for i := range expected {
		if res[i] != expected[i] {
			tt.Errorf("expected %s, got %s", expected[i], res[i])
		}
	}
}

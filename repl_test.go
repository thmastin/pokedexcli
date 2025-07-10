package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Hello WORLD",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   HELLO   world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "   ",
			expected: []string{},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual: %v and expected: %v do not match", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Input: %v does not match Expected: %v", word, expectedWord)
			}
		}
	}
}

func TestGetFirstWord(t *testing.T) {
	cases := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "normal input",
			input:    []string{"hello", "world"},
			expected: "hello",
		},
		{
			name:     "single word",
			input:    []string{"pokeman"},
			expected: "pokeman",
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := GetFirstWord(c.input)
			if actual != c.expected {
				t.Errorf("Expected: %v, Got: %v", c.expected, actual)
			}
		})
	}

}

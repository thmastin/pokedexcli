package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty spaces before, middle, and after",
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			name:     "first word capitalized second word all caps",
			input:    "Hello WORLD",
			expected: []string{"hello", "world"},
		},
		{
			name:     "empty spaces and capitilization",
			input:    "   HELLO   world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{},
		},
		{
			name:     "all spaces input",
			input:    "   ",
			expected: []string{},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
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
		})
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
			actual := getFirstWord(c.input)
			if actual != c.expected {
				t.Errorf("Expected: %v, Got: %v", c.expected, actual)
			}
		})
	}

}

func TestDisplayOutput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "hello",
			expected: "Your command was: hello\n",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "Please enter a command\n",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := displayOutput(c.input)
			if actual != c.expected {
				t.Errorf("Expected: %v, Got: %v", c.expected, actual)
			}
		})
	}
}

func TestProcessCommand(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "one word input no caps and no spaces",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "one word input with caps",
			input:    "HELLO",
			expected: "hello",
		},
		{
			name:     "Multi word input with caps",
			input:    "This is a test command a user may submit",
			expected: "this",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Input all spaces",
			input:    "     ",
			expected: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := processCommand(c.input)
			if actual != c.expected {
				t.Errorf("Expected: %v, Actual: %v", c.expected, actual)
			}
		})
	}
}

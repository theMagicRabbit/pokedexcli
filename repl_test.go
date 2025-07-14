package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input	string
		expected []string
	}{
		{
			input: " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input: "",
			expected: []string{""},
		},
		{
			input: "this little light of mine ",
			expected: []string{"this", "little", "light", "of", "mine"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length of actual output: %d was not expected length: %d", len(actual), len(c.expected))
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: '%s' does not equal actual: '%s'", expectedWord, word)
			}
		}
	}
}


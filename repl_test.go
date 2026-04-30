package main

import (
	"reflect"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		input := c.input
		expected := c.expected
		actual := cleanInput(input)

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %#v but got: %#v", expected, actual)
		}
	}
}

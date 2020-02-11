package main

import (
	"reflect"
	"testing"
)

func TestTrimAndRemoveBlanks(t *testing.T) {

	cases := []struct {
		items      []string
		trimString string
		expected   []string
	}{
		{
			items:      []string{" one ", " two "},
			trimString: " ",
			expected:   []string{"one", "two"},
		},
		{
			items:      []string{"///one///", "///two///", ""},
			trimString: "/",
			expected:   []string{"one", "two"},
		},
		{
			items:      []string{"one", "two"},
			trimString: "---",
			expected:   []string{"one", "two"},
		},
		{
			items:      []string{".-one;,.", "...two;-."},
			trimString: ".-;,",
			expected:   []string{"one", "two"},
		},
		{
			items:      []string{"----------", "two"},
			trimString: "-",
			expected:   []string{"two"},
		},
		{
			items:      []string{"", ""},
			trimString: "\n",
			expected:   nil,
		},
	}

	for i, cs := range cases {
		got := trimAndRemoveBlanks(cs.items, cs.trimString)
		if !reflect.DeepEqual(got, cs.expected) {
			t.Errorf("Test No. %d failed. Want %v, Got %v", i, cs.expected, got)
		}
	}
}

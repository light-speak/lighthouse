package utils

import "testing"

func TestUcFirst(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"", ""},
	}

	for _, test := range tests {
		result := UcFirst(test.input)
		if result != test.expected {
			t.Errorf("UcFirst(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestLcFirst(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"Hello", "hello"},
		{"World", "world"},
		{"", ""},
	}

	for _, test := range tests {
		result := LcFirst(test.input)
		if result != test.expected {
			t.Errorf("LcFirst(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"HelloWorld123", "hello_world123"},
		{"", ""},
	}

	for _, test := range tests {
		result := SnakeCase(test.input)
		if result != test.expected {
			t.Errorf("SnakeCase(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello_world_123", "helloWorld123"},
		{"", ""},
	}

	for _, test := range tests {
		result := CamelCase(test.input)
		if result != test.expected {
			t.Errorf("CamelCase(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

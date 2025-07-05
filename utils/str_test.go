package utils

import "testing"

func TestUcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"world", "World"},
		{"", ""},
		{"h", "H"},
		{"123", "123"},
		{"Hello", "Hello"},
	}

	for _, test := range tests {
		t.Run("Test UcFirst "+test.input, func(t *testing.T) {
			result := UcFirst(test.input)
			if result != test.expected {
				t.Errorf("UcFirst(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestLcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"World", "world"},
		{"", ""},
		{"h", "h"},
		{"123", "123"},
		{"hello", "hello"},
	}

	for _, test := range tests {
		t.Run("Test LcFirst "+test.input, func(t *testing.T) {
			result := LcFirst(test.input)
			if result != test.expected {
				t.Errorf("LcFirst(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"helloWorld", "hello_world"},
		{"HelloWorld", "hello_world"},
		{"hello_world", "hello_world"},
		{"HelloWorld123", "hello_world123"},
		{"", ""},
		{"ABC", "abc"},
		{"abcDEF", "abc_def"},
		{"UserID", "user_id"},
		{"APIResponse", "api_response"},
		{"iOS", "i_os"},
		{"iPhone", "i_phone"},
		{"ID123", "id123"},
	}

	for _, test := range tests {
		t.Run("Test SnakeCase "+test.input, func(t *testing.T) {
			result := SnakeCase(test.input)
			if result != test.expected {
				t.Errorf("SnakeCase(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello_world_123", "helloWorld123"},
		{"", ""},
		{"hello", "hello"},
		{"user_id", "userId"},
		{"api_response", "apiResponse"},
		{"i_os", "iOs"},
		{"i_phone", "iPhone"},
		{"id_123", "id123"},
		{"hello__world", "helloWorld"},
		{"_hello_world", "helloWorld"},
		{"hello_world_", "helloWorld"},
		{"userId", "userId"},
		{"tagId", "tagId"},
		{"BackId", "backId"},
		{"back_id", "backId"},
		{"Back_Id", "backId"},
	}

	for _, test := range tests {
		t.Run("Test CamelCase "+test.input, func(t *testing.T) {
			result := CamelCase(test.input)
			if result != test.expected {
				t.Errorf("CamelCase(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestStrPtr(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"", ""},
		{"123", "123"},
	}

	for _, test := range tests {
		t.Run("Test StrPtr "+test.input, func(t *testing.T) {
			result := StrPtr(test.input)
			if *result != test.expected {
				t.Errorf("StrPtr(%q) = %q; expected %q", test.input, *result, test.expected)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Irregular forms
		{"man", "men"},
		{"woman", "women"},
		{"child", "children"},
		{"tooth", "teeth"},
		{"foot", "feet"},
		{"mouse", "mice"},
		{"person", "people"},

		// Words ending in s, x, z, ch, sh
		{"bus", "buses"},
		{"box", "boxes"},
		{"buzz", "buzzes"},
		{"church", "churches"},
		{"dish", "dishes"},

		// Words ending in y
		{"baby", "babies"}, // consonant + y
		{"boy", "boys"},    // vowel + y
		{"day", "days"},
		{"key", "keys"},

		// Regular words
		{"book", "books"},
		{"pen", "pens"},
		{"dog", "dogs"},
		{"cat", "cats"},
	}

	for _, test := range tests {
		t.Run("Test Pluralize "+test.input, func(t *testing.T) {
			result := Pluralize(test.input)
			if result != test.expected {
				t.Errorf("Pluralize(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

func TestIsVowel(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{'a', true},
		{'e', true},
		{'i', true},
		{'o', true},
		{'u', true},
		{'b', false},
		{'c', false},
		{'z', false},
		{'x', false},
		{'y', false},
	}

	for _, test := range tests {
		t.Run(string(test.input), func(t *testing.T) {
			result := isVowel(test.input)
			if result != test.expected {
				t.Errorf("isVowel(%q) = %v; expected %v", string(test.input), result, test.expected)
			}
		})
	}
}

func TestIsInternalType(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"__type", true},
		{"__schema", true},
		{"_normal", false},
		{"normal", false},
		{"", false},
		{"_", false},
		{"__", true},
		{"___", true},
	}

	for _, test := range tests {
		t.Run("Test IsInternalType "+test.input, func(t *testing.T) {
			result := IsInternalType(test.input)
			if result != test.expected {
				t.Errorf("IsInternalType(%q) = %v; expected %v", test.input, result, test.expected)
			}
		})
	}
}

func TestCamelCaseWithSpecial(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_id", "UserID"},
		{"ip_location", "IPLocation"},
		{"api_url_path", "APIURLPath"},
		{"get_api_url", "GetAPIURL"},
		{"html_id_list", "HTMLIDList"},
		{"oauth_token", "OAuthToken"},
		{"user_api_key", "UserAPIKey"},
		{"cpu_usage", "CPUUsage"},
		{"get_user_id", "GetUserID"},
		{"set_api_url", "SetAPIURL"},
		{"id", "ID"},
		{"user_ip", "UserIP"},
		{"ip", "IP"},
		{"api", "API"},
		{"normal_word", "NormalWord"},
		{"mixed_ip_id_api", "MixedIPIDAPI"},
		{"_id", "ID"},
		{"id_", "ID"},
		{"__id", "ID"},
		{"id__", "ID"},
		{"_ip_", "IP"},
		{"__api__", "API"},
		{"userId", "UserID"},
		{"ipLocation", "IPLocation"},
	}

	for _, test := range tests {
		t.Run("Test CamelCaseWithSpecial "+test.input, func(t *testing.T) {
			result := CamelCaseWithSpecial(test.input)
			if result != test.expected {
				t.Errorf("CamelCaseWithSpecial(%q) = %q; expected %q", test.input, result, test.expected)
			}
		})
	}
}

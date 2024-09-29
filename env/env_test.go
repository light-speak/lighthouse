package env

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{
			key:          "APP_NAME",
			defaultValue: "default_value",
			expected:     "default_value",
		},
	}

	for _, test := range tests {
		result := GetEnv(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("GetEnv(%q, %q) = %q; want %q", test.key, test.defaultValue, result, test.expected)
		}
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		key          string
		defaultValue int
		expected     int
	}{
		{
			key:          "PORT",
			defaultValue: 8080,
			expected:     1111,
		},
	}

	os.Setenv("PORT", "1111")

	for _, test := range tests {
		result := GetEnvInt(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("GetEnvInt(%q, %d) = %d; want %d", test.key, test.defaultValue, result, test.expected)
		}
	}
}


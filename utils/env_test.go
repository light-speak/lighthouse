package utils

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_KEY", "test_value")
	if val := GetEnv("TEST_KEY"); val != "test_value" {
		t.Errorf("GetEnv(TEST_KEY) = %s; want test_value", val)
	}

	// Test with default value
	if val := GetEnv("NON_EXISTENT_KEY", "default"); val != "default" {
		t.Errorf("GetEnv(NON_EXISTENT_KEY, default) = %s; want default", val)
	}
}

func TestGetEnvInt64(t *testing.T) {
	// Test with valid int64
	os.Setenv("TEST_INT64", "123")
	if val := GetEnvInt64("TEST_INT64"); val != 123 {
		t.Errorf("GetEnvInt64(TEST_INT64) = %d; want 123", val)
	}

	// Test with invalid value
	os.Setenv("TEST_INT64_INVALID", "abc")
	if val := GetEnvInt64("TEST_INT64_INVALID", 456); val != 456 {
		t.Errorf("GetEnvInt64(TEST_INT64_INVALID, 456) = %d; want 456", val)
	}
}

func TestGetEnvBool(t *testing.T) {
	// Test with true value
	os.Setenv("TEST_BOOL", "true")
	if val := GetEnvBool("TEST_BOOL"); !val {
		t.Errorf("GetEnvBool(TEST_BOOL) = %v; want true", val)
	}

	// Test with invalid value
	os.Setenv("TEST_BOOL_INVALID", "invalid")
	if val := GetEnvBool("TEST_BOOL_INVALID", true); !val {
		t.Errorf("GetEnvBool(TEST_BOOL_INVALID, true) = %v; want true", val)
	}
}

func TestGetEnvArray(t *testing.T) {
	// Test with valid array
	os.Setenv("TEST_ARRAY", "a,b,c")
	arr := GetEnvArray("TEST_ARRAY", ",")
	if len(arr) != 3 || arr[0] != "a" || arr[1] != "b" || arr[2] != "c" {
		t.Errorf("GetEnvArray(TEST_ARRAY) = %v; want [a b c]", arr)
	}

	// Test with empty value
	os.Setenv("TEST_ARRAY_EMPTY", "")
	defaultArr := []string{"x", "y"}
	arr = GetEnvArray("TEST_ARRAY_EMPTY", ",", defaultArr)
	if len(arr) != 2 || arr[0] != "x" || arr[1] != "y" {
		t.Errorf("GetEnvArray(TEST_ARRAY_EMPTY) = %v; want [x y]", arr)
	}
}

func TestGetEnvFloat64(t *testing.T) {
	// Test with valid float64
	os.Setenv("TEST_FLOAT", "123.45")
	if val := GetEnvFloat64("TEST_FLOAT"); val != 123.45 {
		t.Errorf("GetEnvFloat64(TEST_FLOAT) = %f; want 123.45", val)
	}

	// Test with invalid value
	os.Setenv("TEST_FLOAT_INVALID", "abc")
	if val := GetEnvFloat64("TEST_FLOAT_INVALID", 456.78); val != 456.78 {
		t.Errorf("GetEnvFloat64(TEST_FLOAT_INVALID, 456.78) = %f; want 456.78", val)
	}
}

func TestCleanup(t *testing.T) {
	// Cleanup test environment variables
	os.Unsetenv("TEST_KEY")
	os.Unsetenv("TEST_INT64")
	os.Unsetenv("TEST_INT64_INVALID")
	os.Unsetenv("TEST_BOOL")
	os.Unsetenv("TEST_BOOL_INVALID")
	os.Unsetenv("TEST_ARRAY")
	os.Unsetenv("TEST_ARRAY_EMPTY")
	os.Unsetenv("TEST_FLOAT")
	os.Unsetenv("TEST_FLOAT_INVALID")
}

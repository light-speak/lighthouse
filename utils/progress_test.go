package utils

import (
	"testing"
	"time"
)

func TestSmoothProgress(t *testing.T) {
	SmoothProgress(0, 20, "test", time.Second, true)
	SmoothProgress(20, 40, "test", time.Second, true)
	SmoothProgress(40, 100, "test", time.Second, true)
}

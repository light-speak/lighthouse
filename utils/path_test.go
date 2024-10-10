package utils

import (
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	path, err:= GetProjectPath()
	if err != nil {
		t.Error(err)
	}
	t.Log(path)
}

func TestGetModPath(t *testing.T) {
	modPath, err := GetModPath(nil)
	if err != nil {
		t.Fatalf("GetModPath() failed: %v", err)
	}
	t.Logf("Mod path: %s", modPath)
}


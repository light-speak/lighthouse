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

func TestGetCallerPath(t *testing.T) {
	callerPath, err := GetCallerPath()
	if err != nil {
		t.Fatalf("GetCallerPath() failed: %v", err)
	}
	t.Logf("Caller path: %s", callerPath)
}

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetProjectPath(t *testing.T) {
	path, err := GetProjectPath()
	if err != nil {
		t.Errorf("GetProjectPath() failed: %v", err)
	}

	if path == "" {
		t.Error("GetProjectPath() returned empty path")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("GetProjectPath() returned non-existent path: %s", path)
	}
}

func TestGetGoModPath(t *testing.T) {
	gomodPath, err := GetGoModPath()
	if err != nil {
		t.Errorf("GetGoModPath() failed: %v", err)
		return
	}

	if gomodPath == "" {
		t.Error("GetGoModPath() returned empty path")
		return
	}

	// Verify the path ends with go.mod
	if filepath.Base(gomodPath) != "go.mod" {
		t.Errorf("GetGoModPath() = %v, expected path to end with go.mod", gomodPath)
	}

	// Verify the file exists
	if _, err := os.Stat(gomodPath); os.IsNotExist(err) {
		t.Errorf("GetGoModPath() returned non-existent path: %s", gomodPath)
	}
}

func TestGetModuleRoot(t *testing.T) {
	modRoot, err := GetModuleRoot()
	if err != nil {
		t.Errorf("GetModuleRoot() failed: %v", err)
		return
	}

	if modRoot == "" {
		t.Error("GetModuleRoot() returned empty path")
		return
	}

	// Verify go.mod exists in the module root
	gomodPath := filepath.Join(modRoot, "go.mod")
	if _, err := os.Stat(gomodPath); os.IsNotExist(err) {
		t.Errorf("GetModuleRoot() = %v, but go.mod not found there", modRoot)
	}
}

func TestGetModPath(t *testing.T) {
	// Since GetModPath now uses go env GOMOD, we test that it returns
	// the correct module path for the current module
	modPath, err := GetModPath(nil)
	if err != nil {
		t.Errorf("GetModPath() failed: %v", err)
		return
	}

	if modPath == "" {
		t.Error("GetModPath() returned empty module path")
		return
	}

	// Verify it returns the lighthouse module path
	expectedModPath := "github.com/light-speak/lighthouse"
	if modPath != expectedModPath {
		t.Errorf("GetModPath() = %v, want %v", modPath, expectedModPath)
	}
}

func TestGetPkgPath(t *testing.T) {
	// Get the module root to construct test paths
	modRoot, err := GetModuleRoot()
	if err != nil {
		t.Fatalf("GetModuleRoot() failed: %v", err)
	}

	tests := []struct {
		name     string
		filePath string
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid package path - utils directory",
			filePath: filepath.Join(modRoot, "utils", "path_test.go"),
			want:     "github.com/light-speak/lighthouse/utils",
			wantErr:  false,
		},
		{
			name:     "Valid package path - root directory",
			filePath: filepath.Join(modRoot, "go.mod"),
			want:     "github.com/light-speak/lighthouse",
			wantErr:  false,
		},
		{
			name:     "Empty file path",
			filePath: "",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkgPath, err := GetPkgPath("", tt.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPkgPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && pkgPath != tt.want {
				t.Errorf("GetPkgPath() = %v, want %v", pkgPath, tt.want)
			}
		})
	}
}

func TestGetGoPath(t *testing.T) {
	gopath := GetGoPath()

	// 验证返回的 GOPATH 不为空
	if gopath == "" {
		t.Error("GetGoPath() returned empty path")
	}

	// 验证路径是否存在
	if _, err := os.Stat(gopath); os.IsNotExist(err) {
		t.Errorf("GetGoPath() returned non-existent path: %s", gopath)
	}
}

func TestGetFilePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantFile string
		wantErr  bool
	}{
		{
			name:     "Valid file path",
			path:     "utils/path_test.go",
			wantFile: "path_test.go",
			wantErr:  false,
		},
		{
			name:     "Directory path",
			path:     "utils",
			wantFile: "",
			wantErr:  true,
		},
		{
			name:     "Empty path",
			path:     "",
			wantFile: "",
			wantErr:  true,
		},
		{
			name:     "File path with multiple directories",
			path:     "a/b/c/test.go",
			wantFile: "test.go",
			wantErr:  false,
		},
		{
			name:     "File path with extension",
			path:     "test.txt",
			wantFile: "test.txt",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := GetFilePath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFilePath(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if file == "" {
					t.Error("GetFilePath() returned empty file name")
					return
				}

				if file != tt.wantFile {
					t.Errorf("GetFilePath() = %q, want %q", file, tt.wantFile)
				}
			}
		})
	}
}

func TestGetFileDir(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name:    "Valid file path",
			path:    "utils/path_test.go",
			want:    "utils",
			wantErr: false,
		},
		{
			name:    "Directory path",
			path:    "utils",
			want:    ".",
			wantErr: false,
		},
		{
			name:    "Empty path",
			path:    "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := GetFileDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && dir != tt.want {
				t.Errorf("GetFileDir() = %v, want %v", dir, tt.want)
			}
		})
	}
}

func TestMkdirAll(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Create single directory",
			path:    filepath.Join(tempDir, "test1"),
			wantErr: false,
		},
		{
			name:    "Create nested directories",
			path:    filepath.Join(tempDir, "test2", "nested", "dir"),
			wantErr: false,
		},
		{
			name:    "Create existing directory",
			path:    tempDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := MkdirAll(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("MkdirAll() failed to create directory: %s", tt.path)
				}
			}
		})
	}
}

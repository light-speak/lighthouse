package utils

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
)

func AddImport(file string, importStr string, withUnderscore bool) error {
	// Read the file content
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// Parse the file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, content, parser.ParseComments)
	if err != nil {
		return err
	}

	// Check if import already exists
	exists := false
	for _, imp := range f.Imports {
		if imp.Path.Value == fmt.Sprintf(`"%s"`, importStr) {
			exists = true
			break
		}
	}

	if !exists {
		// Create new import spec
		newImport := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, importStr),
			},
		}

		if withUnderscore {
			newImport.Name = &ast.Ident{Name: "_"}
		}

		// Add to imports
		if f.Imports == nil {
			f.Imports = []*ast.ImportSpec{newImport}
		} else {
			f.Imports = append(f.Imports, newImport)
		}

		// Format and write back
		var buf bytes.Buffer
		err = format.Node(&buf, fset, f)
		if err != nil {
			return err
		}

		err = os.WriteFile(file, buf.Bytes(), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

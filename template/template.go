package template

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/light-speak/lighthouse/utils"
)

type Options struct {
	Path     string
	Template string
	FileName string
	Data     any
	Editable bool

	funcs template.FuncMap
}

// addFunc adds template functions to the template
func (o *Options) addFunc() {
	o.funcs = template.FuncMap{
		"ucFirst":   utils.UcFirst,
		"lcFirst":   utils.LcFirst,
		"camelCase": utils.CamelCase,
		"snakeCase": utils.SnakeCase,
		"funcStart": UserCodeStart,
		"funcEnd":   UserCodeEnd,
		"section":   UserCodeSection,
	}
}

func Render(options *Options) error {
	options.addFunc()

	// Step 1: Load existing content if the file exists
	existingContent := ""
	filePath := fmt.Sprintf("%s/%s.go", options.Path, options.FileName)
	if _, err := os.Stat(filePath); err == nil {
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", filePath, err)
		}
		existingContent = string(contentBytes)
	} else {
		utils.EnsureDir(options.Path)
	}

	// Step 2: Extract user code sections
	userCodeMap := extractUserCode(existingContent)

	// Step 3: Prepare template with functions
	tmpl, err := template.New("").Funcs(options.funcs).Parse(options.Template)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Step 4: Execute template
	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, options.Data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Step 5: Merge user code sections with rendered content
	finalContent := mergeUserCode(renderedContent.String(), userCodeMap)

	// Step 6: Detect required imports
	imports := detectImports(finalContent)
	importsStr := formatImport(imports)

	// Step 7: Add imports to the top of the file
	finalContent = fmt.Sprintf(
		"%s\npackage %s\n\n%s\n\n%s",
		getEditableSection(options),
		getPackageName(options.Path),
		importsStr,
		finalContent,
	)

	// Step 8: Write final content to file
	if err := os.WriteFile(filePath, []byte(finalContent), 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

var (
	funcRegex    = regexp.MustCompile(`(?m)^(\s*)// Func:(\w+) user code start\. Do not remove this comment\.(?:\n|.)*?(?m)^(\s*)// Func:\w+ user code end\. Do not remove this comment\.`)
	sectionRegex = regexp.MustCompile(`(?m)^(\s*)// Section: user code section start\. Do not remove this comment\.(?:\n|.)*?(?m)^(\s*)// Section: user code section end\. Do not remove this comment\.`)
)

// extractUserCode extracts user code sections from the existing content
// it will return a map with the key as the section name and the value as the section content
func extractUserCode(content string) map[string]string {
	userCodeMap := make(map[string]string)

	funcMatches := funcRegex.FindAllStringSubmatch(content, -1)
	for _, match := range funcMatches {
		startIndex := strings.Index(match[0], "user code start. Do not remove this comment.")
		endIndex := strings.LastIndex(match[0], "// Func:")
		if startIndex != -1 && endIndex != -1 {
			userCode := match[0][startIndex+len("user code start. Do not remove this comment.") : endIndex]
			funcName := regexp.MustCompile(`Func:(\w+)`).FindStringSubmatch(match[0])[1]
			userCodeMap[funcName] = strings.TrimSpace(userCode)
		}
	}

	sectionMatches := sectionRegex.FindAllStringSubmatch(content, -1)
	for i, match := range sectionMatches {
		startIndex := strings.Index(match[0], "user code section start. Do not remove this comment.")
		endIndex := strings.LastIndex(match[0], "// Section:")
		if startIndex != -1 && endIndex != -1 {
			userCode := match[0][startIndex+len("user code section start. Do not remove this comment.") : endIndex]
			userCodeMap[fmt.Sprintf("section_%d", i+1)] = strings.TrimSpace(userCode)
		}
	}
	return userCodeMap
}

// mergeUserCode merges user code sections with rendered content
func mergeUserCode(content string, userCodeMap map[string]string) string {
	// Replace function code blocks with user-provided code
	content = funcRegex.ReplaceAllStringFunc(content, func(match string) string {
		submatches := funcRegex.FindStringSubmatch(match)
		indent, actionName := submatches[1], submatches[2]
		if userCode, exists := userCodeMap[actionName]; exists {
			return fmt.Sprintf("%s// Func:%s user code start. Do not remove this comment.\n%s%s\n%s// Func:%s user code end. Do not remove this comment.", indent, actionName, indent, strings.TrimSpace(userCode), indent, actionName)
		}
		return match // If no user code exists, retain the original block
	})

	// Replace section code blocks with user-provided code
	sectionCount := 0
	content = sectionRegex.ReplaceAllStringFunc(content, func(match string) string {
		sectionCount++
		sectionKey := fmt.Sprintf("section_%d", sectionCount)
		submatches := sectionRegex.FindStringSubmatch(match)
		indent := submatches[1]
		if userCode, exists := userCodeMap[sectionKey]; exists {
			return fmt.Sprintf("%s// Section: user code section start. Do not remove this comment.\n%s%s\n%s// Section: user code section end. Do not remove this comment.", indent, indent, strings.TrimSpace(userCode), indent)
		}
		return match // If no user code exists, retain the original block
	})

	return content
}

// getPackageName gets the package name from the path
// for example, if the path is "cmd/cli/generate/cmd", the package name will be "cmd"
func getPackageName(path string) string {
	packageName := filepath.Base(path)
	packageName = strings.TrimSpace(packageName)
	return packageName
}

// getEditableSection gets the editable section from the options
// if the options.Editable is true, it will return the editable section
// otherwise, it will return the non-editable section
func getEditableSection(options *Options) string {
	if options.Editable {
		return "// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF."
	}
	return "// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT."
}

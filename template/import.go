package template

import (
	"fmt"
	"regexp"
	"strings"
)

type Import struct {
	Path  string
	Alias string
}

var importRegexMap = map[string]Import{
	`time\.`: {
		Path: "time",
	},
	`fmt\.`: {
		Path: "fmt",
	},
	`os\.`: {
		Path: "os",
	},
	`path\.`: {
		Path: "path",
	},
	`strings\.`: {
		Path: "strings",
	},
	`strconv\.`: {
		Path: "strconv",
	},
	`regexp\.`: {
		Path: "regexp",
	},
	`\bio\b`: {
		Path: "io",
	},
	`bufio\.`: {
		Path: "bufio",
	},
	`exec\.`: {
		Path: "os/exec",
	},
	`template\.`: {
		Path: "github.com/light-speak/lighthouse/template",
	},
	`command\.`: {
		Path: "github.com/light-speak/lighthouse/command",
	},
	`log\.`: {
		Path: "github.com/light-speak/lighthouse/log",
	},
	`model\.`: {
		Path: "github.com/light-speak/lighthouse/graphql/model",
	},
	`ast\.`: {
		Path: "github.com/light-speak/lighthouse/graphql/ast",
	},
	`gorm\.`: {
		Path: "gorm.io/gorm",
	},
	`context\.`: {
		Path: "context",
	},
	`sync\.`: {
		Path: "sync",
	},
	`utils\.`: {
		Path: "github.com/light-speak/lighthouse/utils",
	},
	`excute\.`: {
		Path: "github.com/light-speak/lighthouse/graphql/excute",
	},
}

// AddImportRegex add a new import regex and path to the importRegexMap
// it will auto import the package in the template
// if alias is not empty, it will use the alias as the import alias
// like import command "github.com/light-speak/lighthouse/command"
func AddImportRegex(regex string, path string, alias string) error {
	_, err := regexp.Compile(regex)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %v", err)
	}
	importRegexMap[regex] = Import{
		Path:  path,
		Alias: alias,
	}
	return nil
}

// FormatImport format imports to string
// if imports is empty, return ""
// if imports has only one import, use single import like "import %s \"%s\""
// if imports has more than one import, use multi import like "import (\n%s\n)"
func formatImport(imports []*Import) string {
	switch len(imports) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(`import %s "%s"`, imports[0].Alias, imports[0].Path)
	default:
		var sb strings.Builder
		sb.WriteString("import (\n")
		for _, imp := range imports {
			if imp.Alias != "" {
				sb.WriteString(fmt.Sprintf(`  %s "%s"`+"\n", imp.Alias, imp.Path))
			} else {
				sb.WriteString(fmt.Sprintf(`  "%s"`+"\n", imp.Path))
			}
		}
		sb.WriteString(")")
		return sb.String()
	}
}

// detectImports detect imports from content
// if content used fmt.Println, it will add "fmt" to imports
// default regex defined in importRegexMap
// other regex will be added by AddImportRegex
func detectImports(content string) []*Import {
	imports := []*Import{}
	for regex, imp := range importRegexMap {
		matched, _ := regexp.MatchString(regex, content)
		if matched {
			imports = append(imports, &imp)
		}
	}
	return imports
}

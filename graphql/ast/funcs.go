package ast

import (
	"fmt"
	"strings"

	"github.com/light-speak/lighthouse/utils"
)

var excludeFieldName = map[string]struct{}{
	"id":        {},
	"createdAt": {},
	"updatedAt": {},
	"deletedAt": {},
}

func Fields(fields map[string]*Field) string {
	var lines []string
	for _, field := range fields {
		if _, ok := excludeFieldName[field.Name]; ok {
			continue
		}
		line := fmt.Sprintf("  %s %s %s", utils.UcFirst(field.Name), field.Type.GetGoType(false), genTag(field))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func genTag(field *Field) string {
	tags := map[string][]string{
		"json": {field.Name},
	}

	// Collect directives
	for _, directive := range GetDirective("tag", field.Directives) {
		if arg := directive.GetArg("name"); arg != nil {
			tags[arg.Name] = append(tags[arg.Name], arg.Value.(string))
		}
	}

	if directive := GetDirective("index", field.Directives); len(directive) == 1 {
		if arg := directive[0].GetArg("name"); arg != nil {
			tags["gorm"] = append(tags["gorm"], fmt.Sprintf("index:%s", arg.Value.(string)))
		} else {
			tags["gorm"] = append(tags["gorm"], "index")
		}
	}

	if directive := GetDirective("unique", field.Directives); len(directive) == 1 {
		tags["gorm"] = append(tags["gorm"], "unique")
	}

	// Build the tag string using strings.Builder
	var builder strings.Builder
	builder.WriteString("`") // Start with a single backtick

	for key, value := range tags {
		builder.WriteString(fmt.Sprintf("%s:\"%s\" ", key, strings.Join(value, ";")))
	}

	builder.WriteString("`") // End with a single backtick
	return builder.String()
}

func Model(typeNode *ObjectNode) string {
	dbName := ""
	trackName := ""
	var builder strings.Builder

	if directive := GetDirective("model", typeNode.Directives); len(directive) == 1 {
		if arg := directive[0].GetArg("name"); arg != nil {
			dbName = arg.Value.(string)
		}
		trackName = "model.Model"
	}
	if directive := GetDirective("softDeleteModel", typeNode.Directives); len(directive) == 1 {
		if arg := directive[0].GetArg("name"); arg != nil {
			dbName = arg.Value.(string)
		}
		trackName = "model.ModelSoftDelete"
	}

	builder.WriteString(fmt.Sprintf("type %s struct {\n", typeNode.GetName()))
	if trackName != "" {
		builder.WriteString(fmt.Sprintf("  %s\n", trackName))
	}
	builder.WriteString(Fields(typeNode.Fields))
	builder.WriteString("\n}\n")

	builder.WriteString(fmt.Sprintf("\nfunc (*%s) IsModel() bool { return true }", typeNode.GetName()))

	if dbName != "" {
		builder.WriteString(fmt.Sprintf("\nfunc (%s) TableName() string { return \"%s\" }", typeNode.GetName(), dbName))
	}
	return builder.String()
}

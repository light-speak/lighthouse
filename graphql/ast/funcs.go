package ast

import (
	"fmt"
	"strings"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

var excludeFieldName = map[string]struct{}{
	"id":         {},
	"created_at": {},
	"updated_at": {},
	"deleted_at": {},
	"__typename": {},
}

func PrefixModels(typeName string) string {
	if strings.HasPrefix(typeName, "*") {
		if strings.HasPrefix(strings.TrimPrefix(typeName, "*"), "[]") {
			return "*[]*models." + strings.TrimPrefix(strings.TrimPrefix(typeName, "*"), "[]")
		}
		return "*models." + strings.TrimPrefix(typeName, "*")
	}
	if strings.HasPrefix(typeName, "[]") {
		return "[]*models." + strings.TrimPrefix(typeName, "[]")
	}
	return "*models." + typeName
}

func Fields(fields map[string]*Field) string {
	var lines []string
	for _, field := range fields {
		if _, ok := excludeFieldName[field.Name]; ok {
			continue
		}
		line := fmt.Sprintf("  %s %s %s", utils.UcFirst(utils.CamelCase(field.Name)), field.Type.GetGoType(false), genTag(field))
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func indexDirective(tags map[string][]string, directive *Directive) error {
	if arg := directive.GetArg("name"); arg != nil {
		tags["gorm"] = append(tags["gorm"], fmt.Sprintf("index:%s", arg.Value.(string)))
	} else {
		tags["gorm"] = append(tags["gorm"], "index")
	}
	return nil
}

func uniqueDirective(tags map[string][]string, directive *Directive) error {
	tags["gorm"] = append(tags["gorm"], "unique")
	return nil
}

func tagDirective(tags map[string][]string, directive *Directive) error {
	if arg := directive.GetArg("name"); arg != nil {
		tags[arg.Name] = append(tags[arg.Name], arg.Value.(string))
	}
	return nil
}

func defaultDirective(tags map[string][]string, directive *Directive) error {
	if arg := directive.GetArg("value"); arg != nil {
		tags["gorm"] = append(tags["gorm"], fmt.Sprintf("default:%s", arg.Value.(string)))
	}
	return nil
}

func typeDirective(tags map[string][]string, directive *Directive) error {
	if arg := directive.GetArg("name"); arg != nil {
		tags["gorm"] = append(tags["gorm"], fmt.Sprintf("type:%s", arg.Value.(string)))
	}
	return nil
}

var directiveFns = map[string]func(map[string][]string, *Directive) error{
	"index":   indexDirective,
	"unique":  uniqueDirective,
	"tag":     tagDirective,
	"default": defaultDirective,
	"type":    typeDirective,
}

func genTag(field *Field) string {
	tags := map[string][]string{
		"json": {field.Name},
	}
	hasType := false

	for _, directive := range field.Directives {
		if fn, ok := directiveFns[directive.Name]; ok {
			if directive.Name == "type" {
				hasType = true
			}
			if err := fn(tags, directive); err != nil {
				log.Error().Err(err).Msgf("failed to apply directive %s to field %s", directive.Name, field.Name)
				return ""
			}
		}
	}
	if !hasType && field.Type.GetRealType().Name == "String" {
		tags["gorm"] = append(tags["gorm"], fmt.Sprintf("type:varchar(%s)", "255"))
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

func BuildRelation(field *Field) string {
	return fmt.Sprintf("{Name: \"%s\", RelationType: ast.%s, ForeignKey: \"%s\", Reference: \"%s\"}", field.Relation.Name, field.Relation.RelationType, field.Relation.ForeignKey, field.Relation.Reference)
}

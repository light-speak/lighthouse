package ast

// var excludeFieldName = map[string]struct{}{
// 	"id":        {},
// 	"createdAt": {},
// 	"updatedAt": {},
// 	"deletedAt": {},
// }

// func Fields(fields []*FieldNode) string {
// 	var lines []string
// 	for _, field := range fields {
// 		if _, ok := excludeFieldName[field.Name]; ok {
// 			continue
// 		}
// 		//TODO: 这里需要处理下 Go Type => TypeRef
// 		// line := fmt.Sprintf("  %s %s %s", utils.UcFirst(field.Name), field.Type.GoType(), genTag(field))
// 		// lines = append(lines, line)
// 	}
// 	return strings.Join(lines, "\n")
// }

// func genTag(field *FieldNode) string {
// 	tags := map[string][]string{
// 		"json": {field.Name},
// 	}

// 	// Collect directives
// 	for _, directive := range field.GetDirectivesByName("tag") {
// 		if arg := directive.GetArg("name"); arg != nil {
// 			tags[arg.Name] = append(tags[arg.Name], arg.Value.Value.(*StringValue).Value)
// 		}
// 	}

// 	if directive := field.GetDirective("index"); directive != nil {
// 		if arg := directive.GetArg("name"); arg != nil {
// 			tags["gorm"] = append(tags["gorm"], fmt.Sprintf("index:%s", arg.Value.Value.(*StringValue).Value))
// 		} else {
// 			tags["gorm"] = append(tags["gorm"], "index")
// 		}
// 	}

// 	if directive := field.GetDirective("unique"); directive != nil {
// 		tags["gorm"] = append(tags["gorm"], "unique")
// 	}

// 	// Build the tag string using strings.Builder
// 	var builder strings.Builder
// 	builder.WriteString("`") // Start with a single backtick

// 	for key, value := range tags {
// 		builder.WriteString(fmt.Sprintf("%s:\"%s\" ", key, strings.Join(value, ";")))
// 	}

// 	builder.WriteString("`") // End with a single backtick
// 	return builder.String()
// }

// func Model(typeNode *TypeNode) string {
// 	dbName := ""
// 	trackName := ""
// 	var builder strings.Builder

// 	if directive := typeNode.GetDirective("model"); directive != nil {
// 		if arg := directive.GetArg("name"); arg != nil {
// 			dbName = arg.Value.Value.(*StringValue).Value
// 		}
// 		trackName = "model.Model"
// 	}
// 	if directive := typeNode.GetDirective("softDeleteModel"); directive != nil {
// 		if arg := directive.GetArg("name"); arg != nil {
// 			dbName = arg.Value.Value.(*StringValue).Value
// 		}
// 		trackName = "model.ModelSoftDelete"
// 	}

// 	builder.WriteString(fmt.Sprintf("type %s struct {\n", typeNode.GetName()))
// 	if trackName != "" {
// 		builder.WriteString(fmt.Sprintf("  %s\n", trackName))
// 	}
// 	builder.WriteString(Fields(typeNode.Fields))
// 	builder.WriteString("\n}\n")

// 	builder.WriteString(fmt.Sprintf("\nfunc (*%s) IsModel() bool { return true }", typeNode.GetName()))

// 	if dbName != "" {
// 		builder.WriteString(fmt.Sprintf("\nfunc (%s) TableName() string { return \"%s\" }", typeNode.GetName(), dbName))
// 	}
// 	return builder.String()
// }

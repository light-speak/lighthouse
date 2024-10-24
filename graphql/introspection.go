package graphql

import (
	"fmt"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/utils"
)

var frontendDirectiveLocations = []string{
	"QUERY",
	"MUTATION",
	"SUBSCRIPTION",
	"FIELD",
	"FRAGMENT_DEFINITION",
	"FRAGMENT_SPREAD",
	"INLINE_FRAGMENT",
	"VARIABLE_DEFINITION",
}

func isInternalType(name string) bool {
	return len(name) >= 2 && name[:2] == "__"
}

// resolveSchemaFields resolves the fields of the __schema query.
func ResolveSchemaFields(qp *parser.QueryParser, field *ast.Field) (interface{}, error) {
	res := make(map[string]interface{})
	for _, cField := range field.Children {
		var err error
		if cField.IsFragment || cField.IsUnion {
			for _, fragmentField := range cField.Children {
				res[fragmentField.Name], err = resolveSchemaField(qp, fragmentField)
				if err != nil {
					return nil, err
				}
			}
		} else {
			res[cField.Name], err = resolveSchemaField(qp, cField)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

// resolveSchemaField resolve the specific field in the __schema query
func resolveSchemaField(qp *parser.QueryParser, field *ast.Field) (interface{}, error) {
	switch field.Name {
	case "queryType", "mutationType", "subscriptionType":
		return resolveOperationType(qp, field.Name)
	case "types":
		return resolveAllTypes(qp, field), nil
	case "directives":
		return resolveAllDirectives(qp, field), nil
	default:
		return nil, nil
	}
}

// resolveOperationType parse the operation type
func resolveOperationType(qp *parser.QueryParser, typeName string) (interface{}, error) {
	opType := utils.UcFirst(typeName[:len(typeName)-4])
	obj := qp.Parser.NodeStore.Objects[opType]
	if obj == nil || len(obj.Fields) == 0 {
		return nil, nil
	}
	return map[string]interface{}{"name": opType}, nil
}

// resolveTypeByName resolves the __type query by type name.
func ResolveTypeByName(qp *parser.QueryParser, field *ast.Field) (interface{}, error) {
	typeName, ok := field.Args["name"].Value.(string)
	if !ok {
		return nil, nil
	}

	node := qp.Parser.NodeStore.Nodes[typeName]
	if node == nil {
		return nil, nil
	}
	return resolveTypeFields(field, node)
}

// resolveAllTypes resolves the "types" field by returning all types except internal ones.
func resolveAllTypes(qp *parser.QueryParser, field *ast.Field) []interface{} {
	var types []interface{}
	for _, node := range qp.Parser.NodeStore.Nodes {
		if !isInternalType(node.GetName()) {
			typeRes, _ := resolveTypeFields(field, node)
			types = append(types, typeRes)
		}
	}
	return types
}

// resolveAllDirectives resolves the "directives" field by returning all relevant directives.
func resolveAllDirectives(qp *parser.QueryParser, field *ast.Field) []interface{} {
	var directives []interface{}
	for _, directive := range qp.Parser.NodeStore.Directives {
		if shouldIncludeDirective(directive) {
			directiveRes, _ := resolveDirectiveFields(field, directive)
			directives = append(directives, directiveRes)
		}
	}
	return directives
}

// shouldIncludeDirective checks if a directive should be included in the response.
func shouldIncludeDirective(directive *ast.DirectiveDefinition) bool {
	for _, location := range directive.Locations {
		if utils.Contains(frontendDirectiveLocations, string(location)) {
			return true
		}
	}
	return false
}

// resolveTypeFields resolves all fields for a given type node.
func resolveTypeFields(field *ast.Field, node ast.Node) (interface{}, error) {
	res := make(map[string]interface{})
	for _, cField := range field.Children {
		var err error
		if cField.IsFragment || cField.IsUnion {
			for _, fragmentField := range cField.Children {
				res[fragmentField.Name], err = resolveTypeField(fragmentField, node)
				if err != nil {
					return nil, err
				}
			}
		} else {
			res[cField.Name], err = resolveTypeField(cField, node)
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

// resolveTypeField resolves a specific field for a given type node.
func resolveTypeField(field *ast.Field, node ast.Node) (interface{}, error) {
	switch field.Name {
	case "name", "description", "kind":
		return getBasicTypeInfo(node, field.Name)
	case "fields", "inputFields", "interfaces", "possibleTypes", "enumValues":
		return resolveComplexTypeField(field, node)
	}
	return nil, nil
}

// getBasicTypeInfo returns basic information about a type.
func getBasicTypeInfo(node ast.Node, infoType string) (interface{}, error) {
	switch infoType {
	case "name":
		return node.GetName(), nil
	case "description":
		return node.GetDescription(), nil
	case "kind":
		return node.GetKind(), nil
	}
	return nil, fmt.Errorf("unknown basic info type: %s", infoType)
}

// resolveComplexTypeField handles more complex type fields.
func resolveComplexTypeField(field *ast.Field, node ast.Node) (interface{}, error) {
	var res []interface{}

	switch field.Name {
	case "fields":
		if node.GetKind() == ast.KindObject || node.GetKind() == ast.KindInterface {
			res = resolveFields(field, node.GetFields())
		}
	case "inputFields":
		if node.GetKind() == ast.KindInputObject {
			res = resolveInputFields(node.(*ast.InputObjectNode).Fields)
		}
	case "interfaces":
		if node.GetKind() == ast.KindObject {
			res = resolveInterfaces(field, node.(*ast.ObjectNode).Interfaces)
		}
	case "possibleTypes":
		if node.GetKind() == ast.KindInterface || node.GetKind() == ast.KindUnion {
			res = resolvePossibleTypes(field, node)
		}
	case "enumValues":
		if node.GetKind() == ast.KindEnum {
			res = resolveEnumValues(node.(*ast.EnumNode).EnumValues)
		}
	}

	if len(res) == 0 {
		return []interface{}{}, nil
	}
	return res, nil
}

// resolveTypeRef creates a type reference, handling nested types like NonNull and List.
func resolveTypeRef(field *ast.Field, typeRef *ast.TypeRef) interface{} {
	if typeRef == nil {
		return nil
	}
	result := make(map[string]interface{})

	for _, cField := range field.Children {
		if cField.IsFragment || cField.IsUnion {

			for _, fragmentField := range cField.Children {
				result[fragmentField.Name] = resolveTypeRefField(fragmentField, typeRef)
			}
		} else {
			result[cField.Name] = resolveTypeRefField(cField, typeRef)
		}
	}
	return result
}

func resolveTypeRefField(field *ast.Field, typeRef *ast.TypeRef) interface{} {
	switch field.Name {
	case "name":
		return typeRef.Name
	case "kind":
		return typeRef.Kind
	case "ofType":
		return resolveTypeRef(field, typeRef.OfType)
	}
	return nil
}

// resolveFields resolves fields for an object or interface type.
func resolveFields(field *ast.Field, fields map[string]*ast.Field) []interface{} {
	var result []interface{}
	for _, f := range fields {
		if !isInternalType(f.Name) {
			fieldRes, _ := resolveFieldFields(field, f)
			result = append(result, fieldRes)
		}
	}
	return result
}

// resolveFieldFields resolves fields for a specific field in a type.
func resolveFieldFields(field *ast.Field, nodeField *ast.Field) (interface{}, error) {
	res := make(map[string]interface{})
	for _, cField := range field.Children {
		if cField.IsFragment || cField.IsUnion {
			for _, fragmentField := range cField.Children {
				switch fragmentField.Name {
				case "name":
					res[fragmentField.Name] = nodeField.Name
				case "description":
					res[fragmentField.Name] = nodeField.Description
				case "args":
					res[fragmentField.Name] = resolveArguments(fragmentField, nodeField.Args)
				case "type":
					res[fragmentField.Name] = resolveTypeRef(fragmentField, nodeField.Type)
				case "isDeprecated":
					res[fragmentField.Name] = nodeField.IsDeprecated
				case "deprecationReason":
					res[fragmentField.Name] = nodeField.DeprecationReason
				}
			}
		} else {
			switch cField.Name {
			case "name":
				res[cField.Name] = nodeField.Name
			case "description":
				res[cField.Name] = nodeField.Description
			case "args":
				res[cField.Name] = resolveArguments(cField, nodeField.Args)
			case "type":
				res[cField.Name] = resolveTypeRef(cField, nodeField.Type)
			case "isDeprecated":
				res[cField.Name] = nodeField.IsDeprecated
			case "deprecationReason":
				res[cField.Name] = nodeField.DeprecationReason
			}
		}
	}
	return res, nil
}

// resolveArguments resolves arguments for a field or directive.
func resolveArguments(field *ast.Field, args map[string]*ast.Argument) []interface{} {
	var result []interface{}

	for _, arg := range args {
		argRes := make(map[string]interface{})
		for _, cField := range field.Children {
			if cField.IsFragment || cField.IsUnion {
				for _, fragmentField := range cField.Children {
					argRes[fragmentField.Name], _ = resolveArgumentField(fragmentField, arg)
				}
			} else {
				argRes[cField.Name], _ = resolveArgumentField(cField, arg)
			}
		}
		result = append(result, argRes)
	}
	if len(result) == 0 {
		return []interface{}{}
	}
	return result
}

func resolveArgumentField(field *ast.Field, arg *ast.Argument) (interface{}, error) {
	switch field.Name {
	case "name":
		return arg.Name, nil
	case "description":
		return arg.Description, nil
	case "type":
		return resolveTypeRef(field, arg.Type), nil
	case "defaultValue":
		return arg.GetDefaultValue(), nil
	}
	return nil, nil
}

// resolveInputFields resolves input fields for an input object type.
func resolveInputFields(inputFields map[string]*ast.Field) []interface{} {
	var result []interface{}
	for _, inputField := range inputFields {
		if !isInternalType(inputField.Name) {
			inputFieldRes := map[string]interface{}{
				"name":        inputField.Name,
				"description": inputField.Description,
				"type":        resolveTypeRef(inputField, inputField.Type),
			}
			result = append(result, inputFieldRes)
		}
	}
	if len(result) == 0 {
		return []interface{}{}
	}
	return result
}

// resolveInterfaces resolves interfaces implemented by an object type.
func resolveInterfaces(field *ast.Field, interfaces map[string]*ast.InterfaceNode) []interface{} {
	var result []interface{}
	for _, iface := range interfaces {
		result = append(result, resolveTypeRef(field, &ast.TypeRef{
			Kind: iface.Kind,
			Name: iface.Name,
		}))
	}
	return result
}

// resolvePossibleTypes resolves possible types for an interface or union.
func resolvePossibleTypes(field *ast.Field, node ast.Node) []interface{} {
	var result []interface{}
	switch n := node.(type) {
	case *ast.InterfaceNode:
		for _, objectNode := range n.PossibleTypes {
			result = append(result, resolveTypeRef(field, &ast.TypeRef{
				Kind: objectNode.Kind,
				Name: objectNode.Name,
			}))
		}
	case *ast.UnionNode:
		for _, objectNode := range n.PossibleTypes {
			result = append(result, resolveTypeRef(field, &ast.TypeRef{
				Kind: objectNode.Kind,
				Name: objectNode.Name,
			}))
		}
	}
	return result
}

// resolveEnumValues resolves enum values for an enum type.
func resolveEnumValues(enumValues map[string]*ast.EnumValue) []interface{} {
	var result []interface{}
	for _, enumValue := range enumValues {
		if !isInternalType(enumValue.Name) {
			valueRes := map[string]interface{}{
				"name":              enumValue.Name,
				"description":       enumValue.Description,
				"isDeprecated":      enumValue.IsDeprecated,
				"deprecationReason": enumValue.DeprecationReason,
			}
			result = append(result, valueRes)
		}
	}
	return result
}

// resolveDirectiveFields resolves fields for a directive.
func resolveDirectiveFields(field *ast.Field, directive *ast.DirectiveDefinition) (interface{}, error) {
	res := make(map[string]interface{})
	for _, cField := range field.Children {
		if cField.IsFragment || cField.IsUnion {
			for _, fragmentField := range cField.Children {
				switch fragmentField.Name {
				case "name":
					res[fragmentField.Name] = directive.Name
				case "description":
					res[fragmentField.Name] = directive.GetDescription()
				case "args":
					res[fragmentField.Name] = resolveArguments(fragmentField, directive.Args)
				case "locations":
					res[fragmentField.Name] = directive.Locations
				}
			}
		} else {
			switch cField.Name {
			case "name":
				res[cField.Name] = directive.Name
			case "description":
				res[cField.Name] = directive.GetDescription()
			case "args":
				res[cField.Name] = resolveArguments(cField, directive.Args)
			case "locations":
				res[cField.Name] = directive.Locations
			}
		}
	}
	return res, nil
}

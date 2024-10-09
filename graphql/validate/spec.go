package validate

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

// validateDirectiveDefinition validates a directive definition node
// 1. check if the directive locations are valid
// 2. check if the directive arguments are valid
func validateDirectiveDefinition(node ast.Node) error {
	directiveDefinition, ok := node.(*ast.DirectiveDefinitionNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not a directive definition",
		}
	}

	for _, loc := range directiveDefinition.Locations {
		if !ast.IsValidLocation(loc) {
			return &errors.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("invalid location: %s", loc),
			}
		}
	}

	return validateArguments(node)
}

// validateScalar validates a scalar node
func validateScalar(node ast.Node) error {
	if _, ok := node.(*ast.ScalarNode); !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not a scalar",
		}
	}
	return nil
}

func validateUnion(node ast.Node) error {
	union, ok := node.(*ast.UnionNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not a union",
		}
	}

	union.TypeNodes = make([]ast.Node, 0, len(union.Types))
	for _, t := range union.Types {
		if node := getValueTypeNode(t); node == nil {
			return &errors.ValidateError{
				Node:    union,
				Message: fmt.Sprintf("type %s not found", t),
			}
		} else {
			union.TypeNodes = append(union.TypeNodes, node)
		}
	}
	return nil
}

func validateEnum(node ast.Node) error {
	enum, ok := node.(*ast.EnumNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not an enum value",
		}
	}

	hasValue := false
	for _, v := range enum.Values {
		if enumDirective := v.GetDirective("enum"); enumDirective != nil {
			ev, _ := ast.ExtractValue(enumDirective.GetArg("value").Value.Value)
			v.Value = int8(ev.(int64))
			hasValue = true
		} else if hasValue {
			return &errors.ValidateError{
				Node:    node,
				Message: "all enum values must have @enum(value: <int>) or none",
			}
		}
	}

	return nil
}

func validateInterface(node ast.Node) error {
	if _, ok := node.(*ast.InterfaceNode); !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not an interface",
		}
	}
	return nil
}

func validateInput(node ast.Node) error {
	input, ok := node.(*ast.InputNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not an input",
		}
	}

	for _, field := range input.Fields {
		err := validateFieldType(field.Type)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateFragment(node ast.Node) error {
	fragment, ok := node.(*ast.FragmentNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not a fragment",
		}
	}

	typeNode := getValueTypeNode(fragment.On)
	if typeNode == nil {
		return &errors.ValidateError{
			Node:    node,
			Message: "fragment on must be a type",
		}
	}
	fragment.Type = typeNode

	return nil
}

func validateType(node ast.Node) error {
	typeNode, ok := node.(*ast.TypeNode)
	if !ok {
		return &errors.ValidateError{
			Node:    node,
			Message: "node is not a type",
		}
	}

	implementTypes := make([]ast.Node, 0, len(typeNode.Implements))
	implementFields := make([]*ast.FieldNode, 0, len(typeNode.Implements))

	for _, interfaceName := range typeNode.Implements {
		interfaceNode := getValueTypeNode(interfaceName)
		if interfaceNode == nil {
			return &errors.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("interface %s not found", interfaceName),
			}
		}
		implementTypes = append(implementTypes, interfaceNode)
		mergedFields, err := mergeFields(implementFields, interfaceNode.GetFields())
		if err != nil {
			return err
		}
		implementFields = mergedFields
	}

	for _, field := range typeNode.GetFields() {
		err := validateFieldType(field.Type)
		if err != nil {
			return err
		}
		// remove the field from the implement fields
		implemented := false
		implementFields, implemented = removeCompatibleField(implementFields, field)
		if !implemented && field.Type.Level == 1 {
			paginate := field.GetDirective("paginate")
			if paginate != nil {
				respType := addPaginationResponseType(field.Type.ElemType)
				field.Type = &ast.FieldType{
					Name:      respType.GetName(),
					Type:      respType,
					IsNonNull: true,
				}
			}
		}
	}

	if len(implementFields) > 0 {
		return &errors.ValidateError{
			Node:    node,
			Message: fmt.Sprintf("field %s not implemented in type %s", implementFields[0].Name, typeNode.GetName()),
		}
	}
	typeNode.ImplementTypes = implementTypes

	return nil
}

// mergeFields
func mergeFields(implementFields []*ast.FieldNode, newFields []*ast.FieldNode) ([]*ast.FieldNode, error) {
	fieldMap := make(map[string]*ast.FieldNode)

	for _, field := range implementFields {
		fieldMap[field.Name] = field
	}

	for _, newField := range newFields {
		if _, exists := fieldMap[newField.Name]; exists {
			return nil, fmt.Errorf("duplicate field: %s", newField.Name)
		}
		fieldMap[newField.Name] = newField
	}

	mergedFields := make([]*ast.FieldNode, 0, len(fieldMap))
	for _, field := range fieldMap {
		mergedFields = append(mergedFields, field)
	}

	return mergedFields, nil
}

// removeField removes a field from a list of fields
func removeCompatibleField(fields []*ast.FieldNode, field *ast.FieldNode) ([]*ast.FieldNode, bool) {
	for i, f := range fields {
		if areTypesCompatible(field.Type, f.Type) {
			return append(fields[:i], fields[i+1:]...), true
		}
	}
	return fields, false
}

func validateField(node ast.Node) error {
	return nil
}

func validateArguments(node ast.Node) error {
	for _, arg := range node.GetArgs() {
		if err := validateFieldType(arg.Type); err != nil {
			return err
		}
	}
	return nil
}

// areTypesCompatible checks if two field types are compatible
func areTypesCompatible(typeA, typeB *ast.FieldType) bool {
	for typeA.IsList && typeB.IsList {
		typeA = typeA.ElemType
		typeB = typeB.ElemType
	}
	return typeA.Name == typeB.Name && typeA.TypeCategory == typeB.TypeCategory && typeA.IsList == typeB.IsList
}

func validateFieldType(fieldType *ast.FieldType) error {
	levels := 0

	elemType := fieldType
	for elemType.IsList {
		levels++
		elemType = elemType.ElemType
	}
	typeName := elemType.Name
	nodeType := getValueTypeNode(typeName)
	if nodeType == nil {
		return fmt.Errorf("validate field type: %s not found", typeName)
	}
	elemType.Type = nodeType
	fieldType.Level = levels

	return nil
}

func validateDirectives(node ast.Node) error {
	for _, directive := range node.GetDirectives() {
		if directiveDefinition := getDirectiveDefinition(directive.GetName()); directiveDefinition == nil {
			return &errors.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("directive %s not found", directive.GetName()),
			}
		}
	}
	return nil
}

func getDirectiveDefinition(name string) *ast.DirectiveDefinitionNode {
	return p.DirectiveMap[name]
}

func getValueTypeNode(name string) ast.Node {
	if node, exists := p.TypeMap[name]; exists {
		return node
	}
	if node, exists := p.UnionMap[name]; exists {
		return node
	}
	if node, exists := p.InterfaceMap[name]; exists {
		return node
	}
	if node, exists := p.InputMap[name]; exists {
		return node
	}
	if node, exists := p.EnumMap[name]; exists {
		return node
	}
	if node, exists := p.ScalarMap[name]; exists {
		return node
	}
	return nil
}

func addPaginationResponseType(fieldType *ast.FieldType) ast.Node {
	elemName := fieldType.Name
	return p.AddType(fmt.Sprintf("%sPaginateResponse", elemName), &ast.TypeNode{
		BaseNode: ast.BaseNode{
			Name:        fmt.Sprintf("%sPaginateResponse", elemName),
			Description: fmt.Sprintf("The %sPaginateResponse type represents a paginated list of %s.", elemName, elemName),
		},
		Fields: []*ast.FieldNode{
			{
				BaseNode: ast.BaseNode{
					Name: "data",
				},
				Type: &ast.FieldType{
					Name:      elemName,
					IsNonNull: true,
					IsList:    true,
					ElemType:  fieldType,
				},
			},
			{
				BaseNode: ast.BaseNode{
					Name: "paginateInfo",
				},
				Type: &ast.FieldType{
					Name:      "PaginateInfo",
					Type:      p.TypeMap["PaginateInfo"],
					IsNonNull: true,
				},
			},
		},
	}, false)
}
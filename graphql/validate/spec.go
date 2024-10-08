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
		if err := validateFieldType(field.Type); err != nil {
			return err
		}
	}

	return nil
}

func validateFragment(node ast.Node) error {
	// TODO: fragment validation
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

	fieldMap := make(map[string]*ast.FieldNode, len(typeNode.GetFields()))
	for _, field := range typeNode.GetFields() {
		if err := validateFieldType(field.Type); err != nil {
			return err
		}
		fieldMap[field.Name] = field
	}

	for _, interfaceName := range typeNode.Implements {
		interfaceNode := getValueTypeNode(interfaceName)
		if interfaceNode == nil {
			return &errors.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("interface %s not found", interfaceName),
			}
		}

		for _, interfaceField := range interfaceNode.GetFields() {
			typeField, exists := fieldMap[interfaceField.Name]
			if !exists {
				return &errors.ValidateError{
					Node:    node,
					Message: fmt.Sprintf("field %s from interface %s not implemented in type %s", interfaceField.Name, interfaceName, typeNode.GetName()),
				}
			}

			if !areTypesCompatible(typeField.Type, interfaceField.Type) {
				return &errors.ValidateError{
					Node:    node,
					Message: fmt.Sprintf("field %s type mismatch: expected %s but got %s", interfaceField.Name, interfaceField.Type.Name, typeField.Type.Name),
				}
			}
		}
	}

	return nil
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
	elemType := fieldType
	for elemType.IsList {
		elemType = elemType.ElemType
	}
	typeName := elemType.Name

	nodeType := getValueTypeNode(typeName)
	if nodeType == nil {
		return fmt.Errorf("validate field type: %s not found", typeName)
	}
	elemType.Type = nodeType
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

package validate

import (
	"fmt"

	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/err"
)

func validateDirectiveDefinition(node ast.Node) error {
	// log.Warn().Msgf("DirectiveDefinitionNode Count: %d", len(p.DirectiveMap))
	directiveDefinition, ok := node.(*ast.DirectiveDefinitionNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a directive definition",
		}
	}
	// log.Info().Msgf("Directive locations: %v", directiveDefinition.Locations)
	// validate locations
	for _, loc := range directiveDefinition.Locations {
		if !ast.IsValidLocation(loc) {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("invalid location: %s", loc),
			}
		}
	}

	return nil
}

func validateScalar(node ast.Node) error {
	_, ok := node.(*ast.ScalarNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a scalar",
		}
	}

	return nil
}

func validateUnion(node ast.Node) error {
	union, ok := node.(*ast.UnionNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a union",
		}
	}
	nodes := make([]ast.Node, 0)
	for _, t := range union.Types {
		node := getValueTypeNode(t)
		if node == nil {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("type %s not found", t),
			}
		}
		nodes = append(nodes, node)
	}
	union.TypeNodes = nodes
	return nil
}

func validateEnum(node ast.Node) error {
	_, ok := node.(*ast.EnumNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a enum value",
		}
	}

	return nil
}

func validateInterface(node ast.Node) error {
	_, ok := node.(*ast.InterfaceNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a interface",
		}
	}

	return nil
}

func validateInput(node ast.Node) error {
	input, ok := node.(*ast.InputNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a input",
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

func validateEnumValue(node ast.Node) error {

	return nil
}

func validateFragment(node ast.Node) error {
	nodeType := node.GetType()
	if nodeType != ast.NodeTypeFragment {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a fragment",
		}
	}
	fragment, ok := node.(*ast.FragmentNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a fragment",
		}
	}

	parentNode := getValueTypeNode(fragment.On)
	if parentNode == nil {
		return &err.ValidateError{
			Node:    node,
			Message: fmt.Sprintf("type %s not found", fragment.On),
		}
	}

	fields := parentNode.GetFields()
	fieldMap := make(map[string]bool)
	for _, field := range fields {
		fieldMap[field.Name] = true
	}

	for _, field := range fragment.GetFields() {
		if _, exists := fieldMap[field.Name]; !exists {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("field %s does not exist in parent type %s", field.Name, parentNode.GetName()),
			}
		}
	}

	return nil
}

func validateType(node ast.Node) error {
	typeNode, ok := node.(*ast.TypeNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a type",
		}
	}
	for _, field := range typeNode.GetFields() {
		err := validateFieldType(field.Type)
		if err != nil {
			return err
		}
	}

	fieldMap := make(map[string]*ast.FieldNode)
	for _, field := range typeNode.GetFields() {
		fieldMap[field.Name] = field
	}

	for _, interfaceName := range typeNode.Implements {
		interfaceNode := getValueTypeNode(interfaceName)
		if interfaceNode == nil {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("interface %s not found", interfaceName),
			}
		}

		// 检查接口中的字段是否在 typeNode 中存在，并且类型是否匹配
		for _, interfaceField := range interfaceNode.GetFields() {
			if typeField, exists := fieldMap[interfaceField.Name]; exists {
				// 检查字段类型是否匹配
				if !areTypesCompatible(typeField.Type, interfaceField.Type) {
					return &err.ValidateError{
						Node: node,
						Message: fmt.Sprintf("field %s type mismatch: expected %s but got %s",
							interfaceField.Name, interfaceField.Type.Name, typeField.Type.Name),
					}
				}
			} else {
				return &err.ValidateError{
					Node: node,
					Message: fmt.Sprintf("field %s from interface %s not implemented in type %s",
						interfaceField.Name, interfaceName, typeNode.GetName()),
				}
			}
		}

	}

	return nil
}

func validateField(node ast.Node) error {
	log.Info().Msgf("field: %s", node.GetName())
	log.Info().Msgf("field: %s", node.(*ast.FieldNode).Name)
	// fields := node.GetFields()
	// for _, field := range fields {
	// 	log.Info().Msgf("field: %s", field.GetName())
	// }

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

// 检查类型是否兼容的辅助函数
func areTypesCompatible(typeA, typeB *ast.FieldType) bool {
	// todo: 检查非空
	return typeA.Name == typeB.Name && typeA.TypeCategory == typeB.TypeCategory && typeA.IsList == typeB.IsList
}

func validateFieldType(fieldType *ast.FieldType) error {
	elemType := fieldType
	typeName := ""
	if elemType.IsList {
		for elemType.IsList {
			elemType = elemType.ElemType
			if !elemType.IsList {
				typeName = elemType.Name
				break
			}
		}
	} else {
		typeName = elemType.Name
	}
	nodeType := getValueTypeNode(typeName)
	if nodeType == nil {
		return fmt.Errorf("validate field type:  %s not found", typeName)
	}
	elemType.Type = nodeType
	return nil
}

func validateDirectives(node ast.Node) error {
	// for _, directive := range node.GetDirectives() {
	// 	log.Debug().Msgf("directive: %s", directive.GetName())
	// }
	return nil
}

func getValueTypeNode(name string) ast.Node {

	typeNode, exists := p.TypeMap[name]
	if exists {
		return typeNode
	}

	unionNode, exists := p.UnionMap[name]
	if exists {
		return unionNode
	}

	interfaceNode, exists := p.InterfaceMap[name]
	if exists {
		return interfaceNode
	}

	inputNode, exists := p.InputMap[name]
	if exists {
		return inputNode
	}

	enumNode, exists := p.EnumMap[name]
	if exists {
		return enumNode
	}

	scalarNode, exists := p.ScalarMap[name]
	if exists {
		return scalarNode
	}

	return nil
}

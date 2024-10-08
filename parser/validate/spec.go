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
	scalar, ok := node.(*ast.ScalarNode)
	if !ok {
		return &err.ValidateError{
			Node:    node,
			Message: "node is not a scalar",
		}
	}
	log.Debug().Msgf("scalar: %s", scalar.GetName())

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
	return nil
}

func validateInterface(node ast.Node) error {
	return nil
}

func validateInput(node ast.Node) error {
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

	fragmentFields := fragment.Fields
	for _, field := range fragmentFields {
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
	log.Debug().Msgf("type: %s", typeNode.GetName())
	for _, field := range typeNode.GetFields() {
		typeName := field.Type.Name
		typeNode := getValueTypeNode(typeName) // String Int , user: User , 
		if typeNode == nil {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("type %s not found", typeName),
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
		elemType := arg.Type
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
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("type %s not found", typeName),
			}
		}
		elemType.Type = nodeType
	}
	return nil
}

func validateDirectives(node ast.Node) error {
	for _, directive := range node.GetDirectives() {
		log.Debug().Msgf("directive: %s", directive.GetName())
	}
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

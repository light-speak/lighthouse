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

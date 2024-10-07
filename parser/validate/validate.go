package validate

import (
	"fmt"

	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/err"
	"github.com/light-speak/lighthouse/parser/parser"
)

var p *parser.Parser

func Validate(node ast.Node, parser *parser.Parser) error {
	// log.Debug().Str("node", node.GetName()).Str("type", string(node.GetType())).Msg("validate node")
	p = parser
	// Create a map of node types to validation functions
	validators := map[ast.NodeType]func(node ast.Node) error{
		ast.NodeTypeDirectiveDefinition: validateDirectiveDefinition,
		ast.NodeTypeDirective:           validateDirective,
		ast.NodeTypeScalar:              validateScalar,
		ast.NodeTypeUnion:               validateUnion,
		ast.NodeTypeEnum:                validateEnum,
		ast.NodeTypeInterface:           validateInterface,
		ast.NodeTypeInput:               validateInput,
		ast.NodeTypeEnumValue:           validateEnumValue,
		ast.NodeTypeFragment:            validateFragment,
		ast.NodeTypeField:               validateField,
	}

	// log.Info().Msgf("scalars: %v", p.ScalarMap)

	// Get the validation function based on the node type
	if validateFunc, exists := validators[node.GetType()]; exists {
		// validate arguments
		err := validateArguments(node)
		if err != nil {
			return err
		}
		// validate directives
		err = validateDirectives(node)
		if err != nil {
			return err
		}
		return validateFunc(node)
	}

	return nil
}

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

func validateDirective(node ast.Node) error {
	return nil
}

func validateScalar(node ast.Node) error {
	return nil
}

func validateUnion(node ast.Node) error {
	// log.Warn().Msgf("Union Node Count: %d", len(p.UnionMap))
	// union, ok := node.(*ast.UnionNode)
	// if !ok {
	// 	return &err.ValidateError{
	// 		Node:    node,
	// 		Message: "node is not a union",
	// 	}
	// }
	// log.Warn().Msgf("Union Node Types: %v", union.Types)
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
		nodeType := GetValueTypeNode(typeName)
		if nodeType == nil {
			return &err.ValidateError{
				Node:    node,
				Message: fmt.Sprintf("type %s not found", typeName),
			}
		}
		elemType.Type = nodeType
		// log.Info().Str("type", typeName).
		// 	Str("elemType", elemType.Name).
		// 	Str("node", node.GetName()).
		// 	Str("nodeType", string(nodeType.GetType())).
		// 	Msgf("type %s found", typeName)
	}
	return nil
}

func validateDirectives(node ast.Node) error {
	return nil
}

func GetValueTypeNode(name string) ast.Node {
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

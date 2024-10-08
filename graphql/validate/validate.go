package validate

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
)

var p *parser.Parser

func Validate(node ast.Node, parser *parser.Parser) error {
	// log.Debug().Str("node", node.GetName()).Str("type", string(node.GetNodeType())).Msg("validate node")
	p = parser
	// Create a map of node types to validation functions
	validators := map[ast.NodeType]func(node ast.Node) error{
		ast.NodeTypeDirectiveDefinition: validateDirectiveDefinition,
		ast.NodeTypeScalar:              validateScalar,
		ast.NodeTypeUnion:               validateUnion,
		ast.NodeTypeEnum:                validateEnum,
		ast.NodeTypeInterface:           validateInterface,
		ast.NodeTypeInput:               validateInput,
		ast.NodeTypeFragment:            validateFragment,
		ast.NodeTypeField:               validateField,
		ast.NodeTypeType:                validateType,
	}

	// log.Info().Msgf("scalars: %v", p.ScalarMap)

	// Get the validation function based on the node type
	if validateFunc, exists := validators[node.GetNodeType()]; exists {
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

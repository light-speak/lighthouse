package ast

import "github.com/light-speak/lighthouse/parser/value"

type ScalarNode struct {
	Name        string
	Description string
	Scalar      ScalarType
	Directives  []*DirectiveNode
}

type ScalarType interface {
	ParseValue(v string) (interface{}, error)
	Serialize(v interface{}) (string, error)
	ParseLiteral(v value.Value) (interface{}, error)
}

func (s *ScalarNode) GetName() string {
	return s.Name
}

func (s *ScalarNode) GetType() NodeType {
	return NodeTypeScalar
}

func (s *ScalarNode) GetDescription() string {
	return s.Description
}

func (s *ScalarNode) IsDeprecated() (bool, string) {
	return false, ""
}

func (s *ScalarNode) GetField(name string) *FieldNode {
	return nil
}

func (s *ScalarNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range s.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

func (s *ScalarNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (s *ScalarNode) GetParent() Node {
	return nil
}

func (s *ScalarNode) GetDirectives() []*DirectiveNode {
	return s.Directives
}

func (s *ScalarNode) GetArgs() []*ArgumentNode {
	return nil
}

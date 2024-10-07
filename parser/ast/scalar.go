package ast

type ScalarNode struct {
	Name        string
	Description string
	Scalar      ScalarType
}

type ScalarType interface {
	ParseValue(value string) (interface{}, error)
	Serialize(value interface{}) (string, error)
	ParseLiteral(value Value) (interface{}, error)
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

func (s *ScalarNode) GetImplements() []string {
	return []string{}
}

func (s *ScalarNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (s *ScalarNode) GetDirectives() []DirectiveNode {
	return []DirectiveNode{}
}

func (s *ScalarNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (s *ScalarNode) IsDeprecated() bool {
	return false
}

func (s *ScalarNode) GetDeprecationReason() string {
	return ""
}

func (s *ScalarNode) IsNonNull() bool {
	return true
}

func (s *ScalarNode) IsList() bool {
	return false
}

func (s *ScalarNode) HasField(name string) bool {
	return false
}

func (s *ScalarNode) HasDirective(name string) bool {
	return false
}

func (s *ScalarNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (s *ScalarNode) GetParent() Node {
	return nil
}

package ast

type UnionNode struct {
	Name        string
	Types       []string
	EntityTypes []*TypeNode
	Description string
}

func (u *UnionNode) GetName() string {
	return u.Name
}

func (u *UnionNode) GetType() NodeType {
	return NodeTypeUnion
}

func (u *UnionNode) GetDescription() string {
	return u.Description
}

func (u *UnionNode) GetImplements() []string {
	return []string{}
}

func (u *UnionNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (u *UnionNode) GetDirectives() []DirectiveNode {
	return []DirectiveNode{}
}

func (u *UnionNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (u *UnionNode) IsDeprecated() bool {
	return false
}

func (u *UnionNode) GetDeprecationReason() string {
	return ""
}

func (u *UnionNode) IsNonNull() bool {
	return true
}

func (u *UnionNode) IsList() bool {
	return false
}

func (u *UnionNode) GetElemType() *FieldType {
	return nil
}

func (u *UnionNode) GetDefaultValue() string {
	return ""
}

func (u *UnionNode) HasField(name string) bool {
	return false
}

func (u *UnionNode) HasDirective(name string) bool {
	return false
}

func (u *UnionNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (u *UnionNode) GetParent() Node {
	return nil
}

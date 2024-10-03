package ast

type UnionNode struct {
	Name        string
	Types       []string
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
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) IsDeprecated() bool {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) GetDeprecationReason() string {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) IsNonNull() bool {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) IsList() bool {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) GetElemType() *FieldType {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) GetDefaultValue() string {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) HasField(name string) bool {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) HasDirective(name string) bool {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) GetDirective(name string) *DirectiveNode {
	panic("not implemented") // TODO: Implement
}

func (u *UnionNode) GetParent() Node {
	panic("not implemented") // TODO: Implement
}


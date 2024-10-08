package ast

type UnionNode struct {
	Name        string
	Types       []string
	TypeNodes   []Node
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

func (u *UnionNode) IsDeprecated() (bool, string) {
	return false, ""
}

func (u *UnionNode) GetField(name string) *FieldNode {
	return nil
}

func (u *UnionNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (u *UnionNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (u *UnionNode) GetParent() Node {
	return nil
}

func (u *UnionNode) GetDirectives() []*DirectiveNode {
	return nil
}

func (u *UnionNode) GetArgs() []*ArgumentNode {
	return nil
}

package ast

type ArgumentNode struct {
	Name         string
	Type         *FieldType
	Value        *ArgumentValue
	DefaultValue *ArgumentValue
	Description  string
	Directives   []DirectiveNode
	Parent       Node
}

func (a *ArgumentNode) GetType() NodeType {
	return NodeTypeArgument
}

func (a *ArgumentNode) GetImplements() []string {
	return []string{}
}

func (a *ArgumentNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (a *ArgumentNode) GetDirectives() []DirectiveNode {
	return a.Directives
}

func (a *ArgumentNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (a *ArgumentNode) IsDeprecated() bool {
	return a.HasDirective("deprecated")
}

func (a *ArgumentNode) GetDeprecationReason() string {
	return ""
}

func (a *ArgumentNode) IsNonNull() bool {
	return a.Type.IsNonNull
}

func (a *ArgumentNode) IsList() bool {
	return a.Type.IsList
}

func (a *ArgumentNode) HasField(name string) bool {
	return false
}

func (a *ArgumentNode) HasDirective(name string) bool {
	for _, directive := range a.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (a *ArgumentNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range a.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (a *ArgumentNode) GetParent() Node {
	return a.Parent
}

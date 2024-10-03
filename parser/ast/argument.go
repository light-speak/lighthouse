package ast

type ArgumentNode struct {
	Name         string
	Type         *FieldType
	Value        string
	Description  string
	Directives   []DirectiveNode
	DefaultValue string
	Parent       Node
}

func (a *ArgumentNode) GetName() string {
	return a.Name
}

func (a *ArgumentNode) GetType() NodeType {
	return NodeTypeArgument
}

func (a *ArgumentNode) GetDescription() string {
	return a.Description
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
	return a.GetDirective("deprecated").Args[0].Value
}

func (a *ArgumentNode) IsNonNull() bool {
	return a.Type.IsNonNull
}

func (a *ArgumentNode) IsList() bool {
	return a.Type.IsList
}

func (a *ArgumentNode) GetElemType() *FieldType {
	return a.Type.ElemType
}

func (a *ArgumentNode) GetDefaultValue() string {
	return a.DefaultValue
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

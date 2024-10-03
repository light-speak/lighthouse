package ast

type FieldNode struct {
	Name        string
	Type        *FieldType
	Description string
	Args        []ArgumentNode
	Directives  []DirectiveNode
	Parent      Node
}

func (f *FieldNode) GetName() string {
	return f.Name
}

func (f *FieldNode) GetType() NodeType {
	return NodeTypeField
}

func (f *FieldNode) GetDescription() string {
	return f.Description
}

func (f *FieldNode) GetImplements() []string {
	return []string{}
}

func (f *FieldNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (f *FieldNode) GetDirectives() []DirectiveNode {
	return f.Directives
}

func (f *FieldNode) GetArgs() []ArgumentNode {
	return f.Args
}

func (f *FieldNode) IsDeprecated() bool {
	return f.HasDirective("deprecated")
}

func (f *FieldNode) GetDeprecationReason() string {
	return f.GetDirective("deprecated").Args[0].Value
}

func (f *FieldNode) IsNonNull() bool {
	return f.Type.IsNonNull
}

func (f *FieldNode) IsList() bool {
	return f.Type.IsList
}

func (f *FieldNode) GetElemType() *FieldType {
	return f.Type.ElemType
}

func (f *FieldNode) GetDefaultValue() string {
	return ""
}

func (f *FieldNode) HasField(name string) bool {
	return false
}

func (f *FieldNode) HasDirective(name string) bool {
	return f.GetDirective(name) != nil
}

func (f *FieldNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range f.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (f *FieldNode) GetParent() Node {
	return f.Parent
}

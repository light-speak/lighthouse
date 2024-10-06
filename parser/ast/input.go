package ast

type InputNode struct {
	Name        string
	Description string
	Fields      []FieldNode
	Directives  []DirectiveNode
}

func (i *InputNode) GetName() string {
	return i.Name
}

func (i *InputNode) GetType() NodeType {
	return NodeTypeInput
}

func (i *InputNode) GetDescription() string {
	return i.Description
}

func (i *InputNode) GetImplements() []string {
	return []string{}
}

func (i *InputNode) GetFields() []FieldNode {
	return i.Fields
}

func (i *InputNode) GetDirectives() []DirectiveNode {
	return i.Directives
}

func (i *InputNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (i *InputNode) IsDeprecated() bool {
	return i.HasDirective("deprecated")
}

func (i *InputNode) GetDeprecationReason() string {
	return ""
}

func (i *InputNode) IsNonNull() bool {
	return true
}

func (i *InputNode) IsList() bool {
	return false
}

func (i *InputNode) GetElemType() *FieldType {
	return nil
}

func (i *InputNode) GetDefaultValue() string {
	return ""
}

func (i *InputNode) HasField(name string) bool {
	for _, field := range i.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func (i *InputNode) HasDirective(name string) bool {
	for _, directive := range i.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (i *InputNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range i.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (i *InputNode) GetParent() Node {
	return nil
}


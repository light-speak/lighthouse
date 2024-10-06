package ast

type EnumNode struct {
	Name        string
	Values      []EnumValueNode
	Description string
	Directives  []DirectiveNode
}

func (e *EnumNode) GetName() string {
	return e.Name
}

func (e *EnumNode) GetType() NodeType {
	return NodeTypeEnum
}

func (e *EnumNode) GetDescription() string {
	return e.Description
}

func (e *EnumNode) GetImplements() []string {
	return []string{}
}

func (e *EnumNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (e *EnumNode) GetDirectives() []DirectiveNode {
	return e.Directives
}

func (e *EnumNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (e *EnumNode) IsDeprecated() bool {
	return e.HasDirective("deprecated")
}

func (e *EnumNode) GetDeprecationReason() string {
	return ""
}

func (e *EnumNode) IsNonNull() bool {
	return true
}

func (e *EnumNode) IsList() bool {
	return false
}

func (e *EnumNode) GetElemType() *FieldType {
	return nil
}

func (e *EnumNode) GetDefaultValue() string {
	return ""
}

func (e *EnumNode) HasField(name string) bool {
	return false
}

func (e *EnumNode) HasDirective(name string) bool {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (e *EnumNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (e *EnumNode) GetParent() Node {
	return nil
}

package ast

type EnumValueNode struct {
	Name        string
	Description string
	Directives  []DirectiveNode
	Parent      Node
}

func (e *EnumValueNode) GetName() string {
	return e.Name
}

func (e *EnumValueNode) GetType() NodeType {
	return NodeTypeEnumValue
}

func (e *EnumValueNode) GetDescription() string {
	return e.Description
}

func (e *EnumValueNode) GetImplements() []string {
	return []string{}
}

func (e *EnumValueNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (e *EnumValueNode) GetDirectives() []DirectiveNode {
	return e.Directives
}

func (e *EnumValueNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (e *EnumValueNode) IsDeprecated() bool {
	return e.HasDirective("deprecated")
}

func (e *EnumValueNode) GetDeprecationReason() string {
	return ""
}

func (e *EnumValueNode) IsNonNull() bool {
	return true
}

func (e *EnumValueNode) IsList() bool {
	return false
}


func (e *EnumValueNode) HasField(name string) bool {
	return false
}

func (e *EnumValueNode) HasDirective(name string) bool {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (e *EnumValueNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range e.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (e *EnumValueNode) GetParent() Node {
	return e.Parent
}

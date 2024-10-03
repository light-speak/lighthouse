package ast

type TypeNode struct {
	Name          string
	Implements    []string
	Fields        []FieldNode
	Description   string
	OperationType OperationType
	Directives    []DirectiveNode
}

func (t *TypeNode) GetName() string {
	return t.Name
}

func (t *TypeNode) GetType() NodeType {
	return NodeTypeType
}

func (t *TypeNode) GetDescription() string {
	return t.Description
}

func (t *TypeNode) GetImplements() []string {
	return t.Implements
}

func (t *TypeNode) GetFields() []FieldNode {
	return t.Fields
}

func (t *TypeNode) GetDirectives() []DirectiveNode {
	return t.Directives
}

func (t *TypeNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (t *TypeNode) IsDeprecated() bool {
	return t.HasDirective("deprecated")
}

func (t *TypeNode) GetDeprecationReason() string {
	return t.GetDirective("deprecated").Args[0].Value
}

func (t *TypeNode) IsNonNull() bool {
	return true
}

func (t *TypeNode) IsList() bool {
	return false
}

func (t *TypeNode) GetElemType() *FieldType {
	return nil
}

func (t *TypeNode) GetDefaultValue() string {
	return ""
}

func (t *TypeNode) HasField(name string) bool {
	return false
}

func (t *TypeNode) HasDirective(name string) bool {
	for _, directive := range t.Directives {
		if directive.Name == name {
			return true
		}
	}
	return false
}

func (t *TypeNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range t.Directives {
		if directive.Name == name {
			return &directive
		}
	}
	return nil
}

func (t *TypeNode) GetParent() Node {
	return nil
}

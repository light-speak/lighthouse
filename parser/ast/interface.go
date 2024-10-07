package ast

type InterfaceNode struct {
	Name        string
	Fields      []FieldNode
	Description string
}

func (i *InterfaceNode) GetName() string {
	return i.Name
}

func (i *InterfaceNode) GetType() NodeType {
	return NodeTypeInterface
}

func (i *InterfaceNode) GetDescription() string {
	return i.Description
}

func (i *InterfaceNode) GetImplements() []string {
	return []string{}
}

func (i *InterfaceNode) GetFields() []FieldNode {
	return i.Fields
}

func (i *InterfaceNode) GetDirectives() []DirectiveNode {
	return []DirectiveNode{}
}

func (i *InterfaceNode) GetArgs() []ArgumentNode {
	return []ArgumentNode{}
}

func (i *InterfaceNode) IsDeprecated() bool {
	return false
}

func (i *InterfaceNode) GetDeprecationReason() string {
	return ""
}

func (i *InterfaceNode) IsNonNull() bool {
	return true
}

func (i *InterfaceNode) IsList() bool {
	return false
}

func (i *InterfaceNode) HasField(name string) bool {
	for _, field := range i.Fields {
		if field.Name == name {
			return true
		}
	}
	return false
}

func (i *InterfaceNode) HasDirective(name string) bool {
	return false
}

func (i *InterfaceNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (i *InterfaceNode) GetParent() Node {
	return nil
}

package ast

type InterfaceNode struct {
	Name        string
	Fields      []*FieldNode
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

func (i *InterfaceNode) IsDeprecated() (bool, string) {
	return false, ""
}

func (i *InterfaceNode) GetField(name string) *FieldNode {
	for _, field := range i.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (i *InterfaceNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (i *InterfaceNode) GetArg(name string) *ArgumentNode {
	return nil
}

func (i *InterfaceNode) GetParent() Node {
	return nil
}

func (i *InterfaceNode) GetDirectives() []*DirectiveNode {
	return nil
}

func (i *InterfaceNode) GetArgs() []*ArgumentNode {
	return nil
}
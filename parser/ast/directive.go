package ast

type DirectiveNode struct {
	Name                string
	Args                []*ArgumentNode
	DirectiveDefinition *DirectiveDefinitionNode
	Parent              Node
}

func (d *DirectiveNode) GetName() string {
	return d.Name
}

func (d *DirectiveNode) GetType() NodeType {
	return NodeTypeDirective
}

func (d *DirectiveNode) GetDescription() string {
	return d.DirectiveDefinition.Description
}

func (d *DirectiveNode) IsDeprecated() (bool, string) {
	return false, ""
}

func (d *DirectiveNode) GetField(name string) *FieldNode {
	return nil
}

func (d *DirectiveNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (d *DirectiveNode) GetArg(name string) *ArgumentNode {
	for _, arg := range d.Args {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}

func (d *DirectiveNode) GetParent() Node {
	return d.Parent
}

func (d *DirectiveNode) GetDirectives() []*DirectiveNode {
	return nil
}

func (d *DirectiveNode) GetArgs() []*ArgumentNode {
	return d.Args
}

func (d *DirectiveNode) GetFields() []*FieldNode {
	return nil
}

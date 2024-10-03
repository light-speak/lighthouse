package ast

type DirectiveNode struct {
	Name                string
	Args                []ArgumentNode
	DirectiveDefinition *DirectiveDefinitionNode
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

func (d *DirectiveNode) GetImplements() []string {
	return []string{}
}

func (d *DirectiveNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (d *DirectiveNode) GetDirectives() []DirectiveNode {
	return []DirectiveNode{}
}

func (d *DirectiveNode) GetArgs() []ArgumentNode {
	return d.Args
}

func (d *DirectiveNode) IsDeprecated() bool {
	return false
}

func (d *DirectiveNode) GetDeprecationReason() string {
	return ""
}

func (d *DirectiveNode) IsNonNull() bool {
	return true
}

func (d *DirectiveNode) IsList() bool {
	return false
}

func (d *DirectiveNode) GetElemType() *FieldType {
	return nil
}

func (d *DirectiveNode) GetDefaultValue() string {
	return ""
}

func (d *DirectiveNode) HasField(name string) bool {
	return false
}

func (d *DirectiveNode) HasDirective(name string) bool {
	return false
}

func (d *DirectiveNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (d *DirectiveNode) GetParent() Node {
	return nil
}

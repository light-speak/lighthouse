package ast

type DirectiveDefinitionNode struct {
	Name        string
	Args        []ArgumentNode
	Description string
	Locations   []DirectiveDefinitionNodeLocation
}

type DirectiveDefinitionNodeLocation string

const (
	DirectiveDefinitionNodeLocationField     DirectiveDefinitionNodeLocation = "FIELD_DEFINITION"
	DirectiveDefinitionNodeLocationArgument  DirectiveDefinitionNodeLocation = "ARGUMENT_DEFINITION"
	DirectiveDefinitionNodeLocationInterface DirectiveDefinitionNodeLocation = "INTERFACE"
	DirectiveDefinitionNodeLocationUnion     DirectiveDefinitionNodeLocation = "UNION"
	DirectiveDefinitionNodeLocationEnum      DirectiveDefinitionNodeLocation = "ENUM"
	DirectiveDefinitionNodeLocationInput     DirectiveDefinitionNodeLocation = "INPUT_OBJECT"
	DirectiveDefinitionNodeLocationScalar    DirectiveDefinitionNodeLocation = "SCALAR"
	DirectiveDefinitionNodeLocationObject    DirectiveDefinitionNodeLocation = "OBJECT"
)

func (d *DirectiveDefinitionNode) GetName() string {
	return d.Name
}

func (d *DirectiveDefinitionNode) GetType() NodeType {
	return NodeTypeDirectiveDefinition
}

func (d *DirectiveDefinitionNode) GetDescription() string {
	return d.Description
}

func (d *DirectiveDefinitionNode) GetImplements() []string {
	return []string{}
}

func (d *DirectiveDefinitionNode) GetFields() []FieldNode {
	return []FieldNode{}
}

func (d *DirectiveDefinitionNode) GetDirectives() []DirectiveNode {
	return []DirectiveNode{}
}

func (d *DirectiveDefinitionNode) GetArgs() []ArgumentNode {
	return d.Args
}

func (d *DirectiveDefinitionNode) IsDeprecated() bool {
	return false
}

func (d *DirectiveDefinitionNode) GetDeprecationReason() string {
	return ""
}

func (d *DirectiveDefinitionNode) IsNonNull() bool {
	return true
}

func (d *DirectiveDefinitionNode) IsList() bool {
	return false
}

func (d *DirectiveDefinitionNode) GetElemType() *FieldType {
	return nil
}

func (d *DirectiveDefinitionNode) GetDefaultValue() string {
	return ""
}

func (d *DirectiveDefinitionNode) HasField(name string) bool {
	return false
}

func (d *DirectiveDefinitionNode) HasDirective(name string) bool {
	return false
}

func (d *DirectiveDefinitionNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (d *DirectiveDefinitionNode) GetParent() Node {
	return nil
}

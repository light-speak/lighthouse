package ast

type DirectiveDefinitionNode struct {
	Name        string
	Args        []ArgumentNode
	Description string
	Locations   []Location
}

type Location string

const (
	LocationQuery              Location = `QUERY`
	LocationMutation           Location = `MUTATION`
	LocationSubscription       Location = `SUBSCRIPTION`
	LocationField              Location = `FIELD`
	LocationFragmentDefinition Location = `FRAGMENT_DEFINITION`
	LocationFragmentSpread     Location = `FRAGMENT_SPREAD`
	LocationInlineFragment     Location = `INLINE_FRAGMENT`

	LocationSchema               Location = `SCHEMA`
	LocationScalar               Location = `SCALAR`
	LocationObject               Location = `OBJECT`
	LocationFieldDefinition      Location = `FIELD_DEFINITION`
	LocationArgumentDefinition   Location = `ARGUMENT_DEFINITION`
	LocationInterface            Location = `INTERFACE`
	LocationUnion                Location = `UNION`
	LocationEnum                 Location = `ENUM`
	LocationEnumValue            Location = `ENUM_VALUE`
	LocationInputObject          Location = `INPUT_OBJECT`
	LocationInputFieldDefinition Location = `INPUT_FIELD_DEFINITION`
	LocationVariableDefinition   Location = `VARIABLE_DEFINITION`
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

package ast

type DirectiveNode struct {
	BaseNode
	Args                []*ArgumentNode
	DirectiveDefinition *DirectiveDefinitionNode
	Parent              Node
}

func (d *DirectiveNode) GetNodeType() NodeType    { return NodeTypeDirective }
func (d *DirectiveNode) GetArgs() []*ArgumentNode { return d.Args }
func (d *DirectiveNode) GetParent() Node          { return d.Parent }
func (d *DirectiveNode) GetArg(name string) *ArgumentNode {
	return findArg(d.Args, name)
}

type DirectiveDefinitionNode struct {
	BaseNode
	Args       []*ArgumentNode
	Locations  []Location
	Repeatable bool
}

func (d *DirectiveDefinitionNode) GetNodeType() NodeType    { return NodeTypeDirectiveDefinition }
func (d *DirectiveDefinitionNode) GetArgs() []*ArgumentNode { return d.Args }
func (d *DirectiveDefinitionNode) GetArg(name string) *ArgumentNode {
	return findArg(d.Args, name)
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

var validLocations = map[Location]struct{}{
	LocationQuery:                {},
	LocationMutation:             {},
	LocationSubscription:         {},
	LocationField:                {},
	LocationFragmentDefinition:   {},
	LocationFragmentSpread:       {},
	LocationInlineFragment:       {},
	LocationSchema:               {},
	LocationScalar:               {},
	LocationObject:               {},
	LocationFieldDefinition:      {},
	LocationArgumentDefinition:   {},
	LocationInterface:            {},
	LocationUnion:                {},
	LocationEnum:                 {},
	LocationEnumValue:            {},
	LocationInputObject:          {},
	LocationInputFieldDefinition: {},
	LocationVariableDefinition:   {},
}

func IsValidLocation(loc Location) bool {
	_, ok := validLocations[loc]
	return ok
}

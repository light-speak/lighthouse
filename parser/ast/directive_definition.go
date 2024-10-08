package ast

type DirectiveDefinitionNode struct {
	Name        string
	Args        []*ArgumentNode
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
	_, exists := validLocations[loc]
	return exists
}

func (d *DirectiveDefinitionNode) GetName() string {
	return d.Name
}

func (d *DirectiveDefinitionNode) GetType() NodeType {
	return NodeTypeDirectiveDefinition
}

func (d *DirectiveDefinitionNode) GetDescription() string {
	return d.Description
}

func (d *DirectiveDefinitionNode) IsDeprecated() (bool, string) {
	return false, ""
}

func (d *DirectiveDefinitionNode) GetField(name string) *FieldNode {
	return nil
}

func (d *DirectiveDefinitionNode) GetDirective(name string) *DirectiveNode {
	return nil
}

func (d *DirectiveDefinitionNode) GetArg(name string) *ArgumentNode {
	for _, arg := range d.Args {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}

func (d *DirectiveDefinitionNode) GetParent() Node {
	return nil
}

func (d *DirectiveDefinitionNode) GetDirectives() []*DirectiveNode {
	return nil
}

func (d *DirectiveDefinitionNode) GetArgs() []*ArgumentNode {
	return d.Args
}

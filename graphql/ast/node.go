package ast

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
)

// Node represents a GraphQL AST node.
type Node interface {
	// GetName returns the name of the node.
	GetName() string

	// GetKind returns the kind of the node.
	GetKind() Kind

	// GetDescription returns the description of the node, if available.
	// It may return an empty string if no description is provided.
	GetDescription() string

	// GetDirectivesByName returns a slice of directives with the specified name.
	// It may return an empty slice if no directives are found.
	GetDirectivesByName(name string) []*Directive

	// Validate validates the node.
	Validate(store *NodeStore) error

	// GetDirectives returns the directives of the node.
	GetDirectives() []*Directive
	GetFields() map[string]*Field
}

type Kind string

const (
	KindScalar      Kind = "SCALAR"
	KindObject      Kind = "OBJECT"
	KindInterface   Kind = "INTERFACE"
	KindUnion       Kind = "UNION"
	KindEnum        Kind = "ENUM"
	KindInputObject Kind = "INPUT_OBJECT"
	KindList        Kind = "LIST"
	KindNonNull     Kind = "NON_NULL"
)

func (k Kind) String() string { return string(k) }

type BaseNode struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Kind        Kind         `json:"kind"`
	Directives  []*Directive `json:"-"`
}

func (n *BaseNode) GetName() string        { return n.Name }
func (n *BaseNode) GetKind() Kind          { return n.Kind }
func (n *BaseNode) GetDescription() string { return n.Description }
func (n *BaseNode) GetDirectivesByName(name string) []*Directive {
	return GetDirective(name, n.Directives)
}
func (n *BaseNode) GetDirectives() []*Directive     { return n.Directives }
func (n *BaseNode) GetFields() map[string]*Field    { return nil }
func (n *BaseNode) Validate(store *NodeStore) error { return nil }

type ObjectNode struct {
	BaseNode
	Fields         map[string]*Field         `json:"fields"`
	InterfaceNames []string                  `json:"-"`
	Interface      map[string]*InterfaceNode `json:"-"`
}

func (o *ObjectNode) GetKind() Kind                { return KindObject }
func (o *ObjectNode) GetFields() map[string]*Field { return o.Fields }
func (o *ObjectNode) Validate(store *NodeStore) error {
	for _, field := range o.Fields {
		if err := field.Validate(store, o.Fields, o, LocationFieldDefinition); err != nil {
			return err
		}
	}
	if err := ValidateDirectives(o.GetName(), o.GetDirectives(), store, LocationObject); err != nil {
		return err
	}
	if len(o.InterfaceNames) > 0 {
		for _, interfaceName := range o.InterfaceNames {
			interfaceNode, ok := store.Interfaces[interfaceName]
			if !ok {
				return &errors.ValidateError{
					NodeName: o.GetName(),
					Message:  fmt.Sprintf("interface %s not found", interfaceName),
				}
			}
			if o.Interface == nil {
				o.Interface = make(map[string]*InterfaceNode)
			}
			o.Interface[interfaceName] = interfaceNode
		}
	}
	return nil
}

type InterfaceNode struct {
	BaseNode
	Fields map[string]*Field
}

func (o *InterfaceNode) GetKind() Kind                { return KindInterface }
func (o *InterfaceNode) GetFields() map[string]*Field { return o.Fields }
func (o *InterfaceNode) Validate(store *NodeStore) error {
	for _, field := range o.Fields {
		if err := field.Validate(store, o.Fields, o, LocationFieldDefinition); err != nil {
			return err
		}
	}
	if err := ValidateDirectives(o.GetName(), o.GetDirectives(), store, LocationInterface); err != nil {
		return err
	}
	return nil
}

type UnionNode struct {
	BaseNode
	TypeNames     map[string]string      `json:"-"`
	PossibleTypes map[string]*ObjectNode `json:"possibleTypes"`
}

func (u *UnionNode) GetKind() Kind { return KindUnion }
func (u *UnionNode) Validate(store *NodeStore) error {
	if err := ValidateDirectives(u.GetName(), u.GetDirectives(), store, LocationUnion); err != nil {
		return err
	}
	for _, typeName := range u.TypeNames {
		typeNode, ok := store.Objects[typeName]
		if !ok {
			return &errors.ValidateError{
				NodeName: u.GetName(),
				Message:  fmt.Sprintf("type %s not found", typeName),
			}
		}
		if u.PossibleTypes == nil {
			u.PossibleTypes = make(map[string]*ObjectNode)
		}
		u.PossibleTypes[typeName] = typeNode
	}
	return nil
}

type EnumNode struct {
	BaseNode
	EnumValues map[string]*EnumValue `json:"enumValues"`
}

func (e *EnumNode) GetKind() Kind { return KindEnum }
func (e *EnumNode) Validate(store *NodeStore) error {
	for _, enumValue := range e.EnumValues {
		if err := enumValue.Validate(store); err != nil {
			return err
		}
	}
	if err := ValidateDirectives(e.GetName(), e.GetDirectives(), store, LocationEnum); err != nil {
		return err
	}
	return nil
}

type EnumValue struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Directives  []*Directive `json:"-"`
	Value       int8         `json:"value"`
}

func (e *EnumValue) Validate(store *NodeStore) error {
	if err := ValidateDirectives(e.Name, e.Directives, store, LocationEnumValue); err != nil {
		return err
	}
	directives := GetDirective("enum", e.Directives)
	if len(directives) == 1 {
		directive := directives[0]
		value := directive.GetArg("value")
		if value != nil {
			if value.Value == nil {
				return &errors.ValidateError{
					NodeName: e.Name,
					Message:  "enum value must have value argument",
				}
			}
			e.Value = int8(value.Value.(int64))
		}
	}
	return nil
}

type InputObjectNode struct {
	BaseNode
	Fields map[string]*Field
}

func (i *InputObjectNode) GetKind() Kind                { return KindInputObject }
func (i *InputObjectNode) GetFields() map[string]*Field { return i.Fields }

func (i *InputObjectNode) Validate(store *NodeStore) error {
	for _, field := range i.Fields {
		if err := field.Validate(store, i.Fields, i, LocationInputFieldDefinition); err != nil {
			return err
		}
	}
	if err := ValidateDirectives(i.GetName(), i.GetDirectives(), store, LocationInputObject); err != nil {
		return err
	}
	return nil
}

type ScalarNode struct {
	BaseNode
	ScalarType ScalarType `json:"-"`
}

func (s *ScalarNode) GetKind() Kind { return KindScalar }
func (s *ScalarNode) Validate(store *NodeStore) error {
	if err := ValidateDirectives(s.GetName(), s.GetDirectives(), store, LocationScalar); err != nil {
		return err
	}
	s.ScalarType = store.ScalarTypes[s.GetName()]
	if s.ScalarType == nil {
		return &errors.ValidateError{
			NodeName: s.GetName(),
			Message:  fmt.Sprintf("scalar %s not found", s.GetName()),
		}
	}
	return nil
}

type ScalarType interface {
	ParseValue(v string) (interface{}, error)
	Serialize(v interface{}) (string, error)
	ParseLiteral(v Value) (interface{}, error)
	GoType() string
}

type Field struct {
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	Args              map[string]*Argument `json:"args"`
	Type              *TypeRef             `json:"type"`
	IsDeprecated      bool                 `json:"isDeprecated"`
	DeprecationReason string               `json:"deprecationReason"`

	Children   map[string]*Field `json:"-"`
	Directives []*Directive      `json:"-"`
	IsFragment bool              `json:"-"`
	Fragment   *FragmentNode     `json:"-"`
}

func (f *Field) Validate(store *NodeStore, objectFields map[string]*Field, objectNode Node, location Location) error {
	if f.IsFragment {
		location = LocationFragmentSpread
	}
	if err := ValidateDirectives(f.Name, f.Directives, store, location); err != nil {
		return err
	}
	directives := GetDirective("deprecated", f.Directives)
	if len(directives) == 1 {
		directive := directives[0]
		deprecationReason := directive.GetArg("reason")
		if deprecationReason != nil {
			f.IsDeprecated = true
			if deprecationReason.Value != nil {
				f.DeprecationReason = deprecationReason.Value.(string)
			} else if deprecationReason.DefaultValue != nil {
				f.DeprecationReason = deprecationReason.DefaultValue.(string)
			} else {
				f.DeprecationReason = "field is deprecated"
			}
		}
	}
	for _, arg := range f.Args {
		if err := arg.Validate(store); err != nil {
			return err
		}
	}
	if f.Type != nil {
		if err := f.Type.Validate(store); err != nil {
			return err
		}
	}
	if f.IsFragment {
		//TODO: validate fragment
	}
	if f.Children != nil {
		// the field is a fragment field or query field, so we need to validate the children
		// the children are the fields of the object type
		var obj Node
		if f.Type.Kind == KindList || f.Type.Kind == KindNonNull {
			obj = f.Type.OfType.TypeNode
		} else {
			obj = f.Type.TypeNode
		}
		if obj.GetKind() != KindObject {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  fmt.Sprintf("field %s must be of object type to have children", f.Name),
			}
		}
		for _, child := range f.Children {
			// the child is a field of the object type
			if err := child.Validate(store, obj.(*ObjectNode).Fields, obj, location); err != nil {
				return err
			}
		}
	}

	return nil
}

type TypeRef struct {
	Kind     Kind     `json:"kind"`
	Name     string   `json:"name"`
	OfType   *TypeRef `json:"ofType"`
	TypeNode Node     `json:"-"`
}

func (t *TypeRef) Validate(store *NodeStore) error {
	if t.Kind == KindNonNull {
		if t.OfType == nil {
			return &errors.ValidateError{
				NodeName: t.Name,
				Message:  "non-null type cannot be null",
			}
		}
		if err := t.OfType.Validate(store); err != nil {
			return err
		}
	} else if t.Kind == KindList {
		if t.OfType == nil {
			return &errors.ValidateError{
				NodeName: t.Name,
				Message:  "list type cannot be null",
			}
		}
		if err := t.OfType.Validate(store); err != nil {
			return err
		}
	} else {
		if t.Name == "" {
			return &errors.ValidateError{
				NodeName: t.Name,
				Message:  "named type cannot be empty",
			}
		}
		if node, ok := store.Nodes[t.Name]; !ok {
			return &errors.ValidateError{
				NodeName: t.Name,
				Message:  fmt.Sprintf("type %s not found", t.Name),
			}
		} else {
			t.Kind = node.GetKind()
			t.TypeNode = node
		}
	}
	return nil
}

func (t *TypeRef) ValidateValue(v interface{}) error {
	switch t.Kind {
	case KindScalar:
		return t.validateScalarValue(v)
	case KindEnum:
		return t.validateEnumValue(v)
	case KindObject:
		return t.validateObjectValue(v)
	case KindInputObject:
		return t.validateInputObjectValue(v)
	case KindList:
		return t.validateListValue(v)
	case KindNonNull:
		if v == nil {
			return fmt.Errorf("value cannot be nil for non-null type %s", t.Name)
		}
		return t.OfType.ValidateValue(v)
	default:
		return fmt.Errorf("unsupported type kind: %s", t.Kind)
	}
}

func (t *TypeRef) validateScalarValue(v interface{}) error {
	switch t.Name {
	case "Int":
		if _, ok := v.(int64); !ok {
			return fmt.Errorf("expected Int, got %T", v)
		}
	case "Float":
		if _, ok := v.(float64); !ok {
			return fmt.Errorf("expected Float, got %T", v)
		}
	case "String":
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected String, got %T", v)
		}
	case "Boolean":
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected Boolean, got %T", v)
		}
	case "ID":
		if _, ok := v.(string); !ok {
			if _, ok := v.(int); !ok {
				return fmt.Errorf("expected ID (String or Int), got %T", v)
			}
		}
	default:
		// For custom scalar types, we might need a more sophisticated validation
		return nil
	}
	return nil
}

func (t *TypeRef) validateEnumValue(v interface{}) error {
	strValue, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected enum value to be string, got %T", v)
	}

	enumNode, ok := t.TypeNode.(*EnumNode)
	if !ok {
		return fmt.Errorf("invalid enum type node")
	}

	for _, enumValue := range enumNode.EnumValues {
		if enumValue.Name == strValue {
			return nil
		}
	}

	return fmt.Errorf("invalid enum value: %s", strValue)
}

func (t *TypeRef) validateObjectValue(v interface{}) error {
	objValue, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected object value to be map[string]interface{}, got %T", v)
	}

	objNode, ok := t.TypeNode.(*ObjectNode)
	if !ok {
		return fmt.Errorf("invalid object type node")
	}

	for _, field := range objNode.Fields {
		fieldValue, exists := objValue[field.Name]
		if !exists && field.Type.Kind == KindNonNull {
			return fmt.Errorf("required field %s is missing", field.Name)
		}
		if exists {
			if err := field.Type.ValidateValue(fieldValue); err != nil {
				return fmt.Errorf("invalid value for field %s: %v", field.Name, err)
			}
		}
	}

	return nil
}

func (t *TypeRef) validateInputObjectValue(v interface{}) error {
	inputObjValue, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected input object value to be map[string]interface{}, got %T", v)
	}

	inputObjNode, ok := t.TypeNode.(*InputObjectNode)
	if !ok {
		return fmt.Errorf("invalid input object type node")
	}

	for _, field := range inputObjNode.Fields {
		fieldValue, exists := inputObjValue[field.Name]
		if !exists && field.Type.Kind == KindNonNull {
			return fmt.Errorf("required input field %s is missing", field.Name)
		}
		if exists {
			if err := field.Type.ValidateValue(fieldValue); err != nil {
				return fmt.Errorf("invalid value for input field %s: %v", field.Name, err)
			}
		}
	}

	return nil
}

func (t *TypeRef) validateListValue(v interface{}) error {
	list, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("expected list, got %T", v)
	}

	for i, item := range list {
		if err := t.OfType.ValidateValue(item); err != nil {
			return fmt.Errorf("invalid value for list item at index %d: %v", i, err)
		}
	}

	return nil
}

type Location string

const (
	LocationQuery                Location = `QUERY`        //TODO: validate
	LocationMutation             Location = `MUTATION`     //TODO: validate
	LocationSubscription         Location = `SUBSCRIPTION` //TODO: validate
	LocationField                Location = `FIELD`        //TODO: validate
	LocationFragmentDefinition   Location = `FRAGMENT_DEFINITION`
	LocationFragmentSpread       Location = `FRAGMENT_SPREAD`
	LocationInlineFragment       Location = `INLINE_FRAGMENT`
	LocationSchema               Location = `SCHEMA` //TODO: validate
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

type Directive struct {
	Name       string               `json:"name"`
	Args       map[string]*Argument `json:"args"`
	Definition *DirectiveDefinition `json:"-"`
}

func (d *Directive) GetArg(name string) *Argument {
	return d.Args[name]
}

func (d *Directive) Validate(store *NodeStore, location Location) error {
	d.Definition = store.Directives[d.Name]
	if d.Definition == nil {
		return &errors.ValidateError{
			NodeName: d.Name,
			Message:  fmt.Sprintf("directive %s not found", d.Name),
		}
	}
	match := false
	for _, loc := range d.Definition.Locations {
		if loc == location {
			match = true
			break
		}
	}
	if !match {
		return &errors.ValidateError{
			NodeName: d.Name,
			Message:  fmt.Sprintf("directive %s is not valid for location %s", d.Name, location),
		}
	}

	d.Definition.Directives = append(d.Definition.Directives, d)
	return nil
}

type DirectiveDefinition struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Args        map[string]*Argument `json:"args"`
	Locations   []Location           `json:"locations"`
	Repeatable  bool                 `json:"repeatable"`

	Directives []*Directive `json:"-"`
}

func (d *DirectiveDefinition) Validate(store *NodeStore) error {
	for _, arg := range d.Args {
		if err := arg.Validate(store); err != nil {
			return err
		}
	}
	if len(d.Locations) == 0 {
		return &errors.ValidateError{
			NodeName: d.Name,
			Message:  "directive must have at least one location",
		}
	}
	for _, loc := range d.Locations {
		if !IsValidLocation(loc) {
			return &errors.ValidateError{
				NodeName: d.Name,
				Message:  fmt.Sprintf("invalid location: %s", loc),
			}
		}
	}
	// for _, directive := range d.Directives {
	// 	for _, arg := range directive.Args {

	// 	}
	// }
	return nil
}

type FragmentNode struct {
	Name       string            `json:"name"`
	On         string            `json:"on"`
	Object     *ObjectNode       `json:"-"`
	Directives []*Directive      `json:"-"`
	Fields     map[string]*Field `json:"fields"`
}

func (f *FragmentNode) Validate(store *NodeStore) error {
	objectNode, ok := store.Objects[f.On]
	if !ok {
		return &errors.ValidateError{
			NodeName: f.Name,
			Message:  fmt.Sprintf("type %s not found", f.On),
		}
	}
	f.Object = objectNode
	if err := ValidateDirectives(f.Name, f.Directives, store, LocationFragmentDefinition); err != nil {
		return err
	}

	for _, field := range f.Fields {
		if err := field.Validate(store, f.Object.Fields, f.Object, LocationInlineFragment); err != nil {
			return err
		}
	}
	return nil
}

func ValidateDirectives(name string, directives []*Directive, store *NodeStore, location Location) error {
	directiveNames := make(map[string]int)
	for _, directive := range directives {
		if err := directive.Validate(store, location); err != nil {
			return err
		}
		directiveNames[directive.Name]++
	}
	for directiveName, count := range directiveNames {
		directiveDefinition := store.Directives[directiveName]
		if directiveDefinition == nil {
			return &errors.ValidateError{
				NodeName: name,
				Message:  fmt.Sprintf("directive %s not found", directiveName),
			}
		}
		if !directiveDefinition.Repeatable && count > 1 {
			return &errors.ValidateError{
				NodeName: name,
				Message:  fmt.Sprintf("directive %s is not repeatable but used %d times", directiveName, count),
			}
		}
	}
	return nil
}

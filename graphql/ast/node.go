package ast

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/utils"
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
	GetPossibleTypes() map[string]*ObjectNode
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
	Name          string                    `json:"name"`
	Description   *string                   `json:"description"`
	Fields        map[string]*Field         `json:"fields"`
	InputFields   map[string]*Field         `json:"inputFields"`
	Interfaces    map[string]*InterfaceNode `json:"interfaces"`
	EnumValues    map[string]*EnumValue     `json:"enumValues"`
	PossibleTypes map[string]*ObjectNode    `json:"possibleTypes"`

	Kind       Kind         `json:"kind"`
	Directives []*Directive `json:"-"`
	IsReserved bool         `json:"-"`
}

func (n *BaseNode) GetName() string { return n.Name }
func (n *BaseNode) GetKind() Kind   { return n.Kind }
func (n *BaseNode) GetDescription() string {
	if n.Description == nil {
		return ""
	}
	return *n.Description
}
func (n *BaseNode) GetDirectivesByName(name string) []*Directive {
	return GetDirective(name, n.Directives)
}
func (n *BaseNode) GetDirectives() []*Directive              { return n.Directives }
func (n *BaseNode) GetFields() map[string]*Field             { return nil }
func (n *BaseNode) Validate(store *NodeStore) error          { return nil }
func (n *BaseNode) GetPossibleTypes() map[string]*ObjectNode { return n.PossibleTypes }

type ObjectNode struct {
	BaseNode
	Fields         map[string]*Field `json:"fields"`
	InterfaceNames []string          `json:"-"`
	IsModel        bool              `json:"-"`
}

func (o *ObjectNode) GetFields() map[string]*Field { return o.Fields }
func (o *ObjectNode) Validate(store *NodeStore) error {
	o.Fields["__typename"] = &Field{
		Name: "__typename",
		Type: &TypeRef{
			Kind: KindNonNull,
			OfType: &TypeRef{
				Kind: KindScalar,
				Name: "String",
			},
		},
	}
	if err := o.ParseObjectDirectives(store); err != nil {
		return err
	}

	for _, field := range o.Fields {
		if err := field.Validate(store, o.Fields, o, LocationFieldDefinition, nil, nil); err != nil {
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
			if o.Interfaces == nil {
				o.Interfaces = make(map[string]*InterfaceNode)
			}
			o.Interfaces[interfaceName] = interfaceNode
			if interfaceNode.PossibleTypes == nil {
				interfaceNode.PossibleTypes = make(map[string]*ObjectNode)
			}
			interfaceNode.PossibleTypes[o.GetName()] = o
		}
	}
	return nil
}

type InterfaceNode struct {
	BaseNode
	Fields map[string]*Field `json:"fields"`
}

func (o *InterfaceNode) GetFields() map[string]*Field { return o.Fields }
func (o *InterfaceNode) Validate(store *NodeStore) error {
	for _, field := range o.Fields {
		if err := field.Validate(store, o.Fields, o, LocationFieldDefinition, nil, nil); err != nil {
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
	TypeNames map[string]string `json:"-"`
}

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
	Name              string       `json:"name"`
	Description       *string      `json:"description"`
	Directives        []*Directive `json:"-"`
	Value             int8         `json:"-"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason *string      `json:"deprecationReason"`
}

func (e *EnumValue) Validate(store *NodeStore) error {
	if err := ValidateDirectives(e.Name, e.Directives, store, LocationEnumValue); err != nil {
		return err
	}
	enum := GetDirective("enum", e.Directives)
	if len(enum) == 1 {
		directive := enum[0]
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
	deprecated := GetDirective("deprecated", e.Directives)
	if len(deprecated) == 1 {
		e.IsDeprecated = true
		reason := deprecated[0].GetArg("reason")
		if reason != nil {
			if reason.Value != nil {
				e.DeprecationReason = utils.StrPtr(reason.Value.(string))
			} else {
				e.DeprecationReason = utils.StrPtr(reason.DefaultValue.(string))
			}
		}
	}
	return nil
}

type InputObjectNode struct {
	BaseNode
	Fields map[string]*Field `json:"inputFields"`
}

func (i *InputObjectNode) GetFields() map[string]*Field { return i.Fields }

func (i *InputObjectNode) Validate(store *NodeStore) error {
	for _, field := range i.Fields {
		if err := field.Validate(store, i.Fields, i, LocationInputFieldDefinition, nil, nil); err != nil {
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

func (s *ScalarNode) Validate(store *NodeStore) error {
	if err := ValidateDirectives(s.GetName(), s.GetDirectives(), store, LocationScalar); err != nil {
		return err
	}
	return nil
}

type ScalarType interface {
	ParseValue(v string) (interface{}, error)
	Serialize(v interface{}) (string, error)
	ParseLiteral(v interface{}) (interface{}, error)
	GoType() string
}

type Field struct {
	Alias             string               `json:"-"`
	Name              string               `json:"name"`
	Description       *string              `json:"description"`
	Args              map[string]*Argument `json:"args"`
	Type              *TypeRef             `json:"type"`
	IsDeprecated      bool                 `json:"isDeprecated"`
	DeprecationReason *string              `json:"deprecationReason"`

	Children   map[string]*Field `json:"-"`
	Directives []*Directive      `json:"-"`
	IsFragment bool              `json:"-"`
	IsUnion    bool              `json:"-"`
	Fragment   *Fragment         `json:"-"`

	DefinitionDirectives []*Directive `json:"-"`
	Relation             *Relation    `json:"-"`
}

type RelationType string

const (
	RelationTypeBelongsTo RelationType = "RelationTypeBelongsTo"
	RelationTypeHasMany   RelationType = "RelationTypeHasMany"
	RelationTypeHasOne    RelationType = "RelationTypeHasOne"
)

type Relation struct {
	Name         string       `json:"relation"`
	ForeignKey   string       `json:"foreignKey"`
	Reference    string       `json:"reference"`
	RelationType RelationType `json:"relationType"`
}

func (f *Field) Validate(store *NodeStore, objectFields map[string]*Field, objectNode Node, location Location, fragments map[string]*Fragment, args map[string]*Argument) error {
	if f.IsFragment {
		if fragments == nil {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  "fragments not found",
			}
		}
		if fragments[f.Type.Name] == nil {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  fmt.Sprintf("fragment %s not found", f.Type.Name),
			}
		}
		frag := fragments[f.Type.Name]
		if frag.Object == nil {
			if err := frag.Validate(store, fragments); err != nil {
				return err
			}
		}
		f.Type.TypeNode = frag.Object
		f.Children = frag.Fields
		location = LocationFragmentSpread
	}

	if err := ValidateDirectives(f.Name, f.Directives, store, location); err != nil {
		return err
	}
	if err := f.ParseFieldDirectives(store, objectNode); err != nil {
		return err
	}

	if f.Type != nil {
		if !f.IsFragment {
			if err := f.Type.Validate(store); err != nil {
				return err
			}
		}
	} else {
		if objectNode == nil {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  "field type must be object type, but got nil",
			}
		} else {
			if objectNode.GetFields()[f.Name] != nil {
				f.Type = objectNode.GetFields()[f.Name].Type
				realType := f.Type.GetRealType()
				if realType.TypeNode == nil || (realType.TypeNode.GetKind() == KindScalar && realType.TypeNode.(*ScalarNode).ScalarType == nil) {
					return &errors.ValidateError{
						NodeName: f.Name,
						Message:  "field type must be scalar type",
					}
				}
				// merge
				f.DefinitionDirectives = append(f.DefinitionDirectives, objectNode.GetFields()[f.Name].Directives...)
			} else {
				// if the field is not found, it means the field is a fragment field
				// we need to validate the fragment field
				if f.IsFragment {
					for _, child := range f.Children {
						if err := child.Validate(store, objectFields, objectNode, location, fragments, args); err != nil {
							return err
						}
					}
				} else {
					return &errors.ValidateError{
						NodeName: f.Name,
						Message:  "field type must be object type, and the field is not a fragment field",
					}
				}
			}
		}
	}

	if f.IsUnion {
		if f.Type.TypeNode.GetKind() != KindObject {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  "field type must be object type",
			}
		}
		obj := objectNode.GetPossibleTypes()[f.Type.Name]
		if obj == nil {
			return &errors.ValidateError{
				NodeName: f.Name,
				Message:  "field type must be a possible type of the union or interface",
			}
		}
		for _, child := range f.Children {
			if err := child.Validate(store, obj.GetFields(), obj, location, fragments, nil); err != nil {
				return err
			}
		}
		return nil
	}

	if f.Children != nil {
		// the field is a fragment field or query field, so we need to validate the children
		// the children are the fields of the object type

		var obj Node
		t := f.Type
		for t.Kind == KindList || t.Kind == KindNonNull {
			t = t.OfType
		}
		obj = t.TypeNode

		for _, child := range f.Children {
			// the child is a field of the object type
			if err := child.Validate(store, obj.GetFields(), obj, location, fragments, nil); err != nil {
				return err
			}
		}
	}
	for _, arg := range f.Args {
		if err := arg.Validate(store, args, objectNode.GetFields()[f.Name]); err != nil {
			return err
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

func (t *TypeRef) GetGoName() string {
	switch t.Kind {
	case KindList:
		return t.OfType.GetGoName()
	case KindNonNull:
		return t.OfType.GetGoName()
	default:
		return utils.UcFirst(t.Name)
	}
}

func (t *TypeRef) GetRealType() *TypeRef {
	if t.Kind == KindNonNull {
		return t.OfType.GetRealType()
	}
	if t.Kind == KindList {
		return t.OfType.GetRealType()
	}
	return t
}

func (t *TypeRef) GetGoType(NonNull bool) string {
	if t == nil {
		return "interface{}"
	}

	switch t.Kind {
	case KindScalar:
		return t.TypeNode.(*ScalarNode).ScalarType.GoType()
	case KindEnum, KindObject, KindInputObject:
		if NonNull {
			return t.Name
		}
		return "*" + t.Name
	case KindList:
		if NonNull {
			return "[]" + t.OfType.GetGoType(false)
		}
		return "*[]" + t.OfType.GetGoType(false)
	case KindNonNull:
		return t.OfType.GetGoType(true)
	}
	return "any"
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
			if t.Kind == KindScalar && t.TypeNode.(*ScalarNode).ScalarType == nil {
				t.TypeNode.(*ScalarNode).ScalarType = store.ScalarTypes[t.Name]
			}
		}
	}
	return nil
}

func (t *TypeRef) ValidateValue(v interface{}, isVariable bool) error {
	switch t.Kind {
	case KindScalar:
		return t.validateScalarValue(v, isVariable)
	case KindEnum:
		return t.validateEnumValue(v)
	case KindObject:
		return t.validateObjectValue(v, isVariable)
	case KindInputObject:
		return t.validateInputObjectValue(v, isVariable)
	case KindList:
		return t.validateListValue(v, isVariable)
	case KindNonNull:
		if v == nil {
			return fmt.Errorf("value cannot be nil for non-null type %s", t.Name)
		}
		return t.OfType.ValidateValue(v, isVariable)
	default:
		return fmt.Errorf("unsupported type kind: %s", t.Kind)
	}
}

func (t *TypeRef) validateScalarValue(v interface{}, isVariable bool) error {
	var err error
	if isVariable {
		_, err = t.TypeNode.(*ScalarNode).ScalarType.ParseValue(v.(string))
	} else {
		_, err = t.TypeNode.(*ScalarNode).ScalarType.ParseLiteral(v)
	}
	if err != nil {
		return err
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

func (t *TypeRef) validateObjectValue(v interface{}, isVariable bool) error {
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
			if err := field.Type.ValidateValue(fieldValue, isVariable); err != nil {
				return fmt.Errorf("invalid value for field %s: %v", field.Name, err)
			}
		}
	}

	return nil
}

func (t *TypeRef) validateInputObjectValue(v interface{}, isVariable bool) error {
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
			if err := field.Type.ValidateValue(fieldValue, isVariable); err != nil {
				return fmt.Errorf("invalid value for input field %s: %v", field.Name, err)
			}
		}
	}

	return nil
}

func (t *TypeRef) validateListValue(v interface{}, isVariable bool) error {
	list, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("expected list, got %T", v)
	}

	for i, item := range list {
		if err := t.OfType.ValidateValue(item, isVariable); err != nil {
			return fmt.Errorf("invalid value for list item at index %d: %v", i, err)
		}
	}

	return nil
}

type Location string

const (
	LocationQuery                Location = `QUERY`
	LocationMutation             Location = `MUTATION`
	LocationSubscription         Location = `SUBSCRIPTION`
	LocationField                Location = `FIELD`
	LocationFragmentDefinition   Location = `FRAGMENT_DEFINITION`
	LocationFragmentSpread       Location = `FRAGMENT_SPREAD`
	LocationInlineFragment       Location = `INLINE_FRAGMENT`
	LocationSchema               Location = `SCHEMA` // deprecated
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
	Description *string              `json:"description"`
	Args        map[string]*Argument `json:"args"`
	Locations   []Location           `json:"locations"`
	Repeatable  bool                 `json:"repeatable"`

	Directives []*Directive `json:"-"`
}

func (d *DirectiveDefinition) GetDescription() string {
	if d.Description == nil {
		return ""
	}
	return *d.Description
}

func (d *DirectiveDefinition) Validate(store *NodeStore) error {
	for _, arg := range d.Args {
		if err := arg.Validate(store, nil, nil); err != nil {
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
	for _, directive := range d.Directives {
		for _, arg := range directive.Args {
			defArg := d.Args[arg.Name]
			if defArg == nil {
				return &errors.ValidateError{
					NodeName: d.Name,
					Message:  fmt.Sprintf("required argument %s is missing", arg.Name),
				}
			}
			kind := defArg.Type.Kind
			for {
				if kind == KindNonNull || kind == KindList {
					kind = defArg.Type.OfType.Kind
				} else {
					break
				}
			}
			if kind == KindEnum {
				arg.Value = arg.Type.Name
			}
			arg.Type = defArg.Type
		}
	}
	return nil
}

type Fragment struct {
	Name       string            `json:"name"`
	On         string            `json:"on"`
	Object     *ObjectNode       `json:"-"`
	Directives []*Directive      `json:"-"`
	Fields     map[string]*Field `json:"fields"`
}

func (f *Fragment) Validate(store *NodeStore, fragments map[string]*Fragment) error {
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
		if err := field.Validate(store, f.Object.Fields, f.Object, LocationInlineFragment, fragments, nil); err != nil {
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

type Argument struct {
	Name         string       `json:"name"`
	Description  *string      `json:"description"`
	Directives   []*Directive `json:"-"`
	Type         *TypeRef     `json:"type"`
	DefaultValue any          `json:"default_value"`
	Value        any          `json:"-"`
	IsVariable   bool         `json:"-"`
	IsReference  bool         `json:"-"`
}

func (a *Argument) GetDefaultValue() *string {
	if a.DefaultValue == nil {
		return nil
	}
	str := fmt.Sprintf("%v", a.DefaultValue)
	return &str
}

func (a *Argument) Validate(store *NodeStore, args map[string]*Argument, field *Field) error {
	location := LocationArgumentDefinition
	if a.IsVariable {
		location = LocationVariableDefinition
	}
	if a.IsReference {
		if args == nil {
			return &errors.ValidateError{
				NodeName: a.Name,
				Message:  "variable arguments not found",
			}
		}

		name, ok := a.Value.(string)

		if !ok {
			return &errors.ValidateError{
				NodeName: a.Name,
				Message:  "variable argument must be string",
			}
		}
		if args[name] == nil {
			return &errors.ValidateError{
				NodeName: a.Name,
				Message:  "variable argument not found",
			}
		}
		a.Value = args[name].Value
		a.Type = args[name].Type
	}
	if a.Type == nil {
		if field == nil {
			return &errors.ValidateError{
				NodeName: a.Name,
				Message:  "argument type not found, the field is not a query field",
			}
		}
		a.Type = field.Args[a.Name].Type
	}
	if err := a.Type.Validate(store); err != nil {
		return err
	}
	if a.Value != nil && !a.IsReference {
		if err := a.Type.ValidateValue(a.Value, false); err != nil {
			return err
		}
	}
	if a.DefaultValue != nil && !a.IsReference {
		if err := a.Type.ValidateValue(a.DefaultValue, false); err != nil {
			return err
		}
	}
	if a.Value != nil && a.IsReference {
		if err := a.Type.ValidateValue(a.Value, true); err != nil {
			return err
		}
	}

	err := ValidateDirectives(a.Name, a.Directives, store, location)
	if err != nil {
		return err
	}
	return nil
}

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
	Validate(store *NodeStore) errors.GraphqlErrorInterface

	// GetDirectives returns the directives of the node.
	GetDirectives() []*Directive
	GetFields() map[string]*Field
	GetPossibleTypes() map[string]*ObjectNode
}

type Locationable interface {
	GetLocation() errors.GraphqlLocation
}

type BaseLocation struct {
	Line   int `json:"-"`
	Column int `json:"-"`
}

func (l *BaseLocation) GetLocation() *errors.GraphqlLocation {
	return &errors.GraphqlLocation{
		Line:   l.Line,
		Column: l.Column,
	}
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

	// IsMain is true when the node is the main service
	IsMain bool `json:"-"`
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
func (n *BaseNode) GetDirectives() []*Directive                            { return n.Directives }
func (n *BaseNode) GetFields() map[string]*Field                           { return nil }
func (n *BaseNode) Validate(store *NodeStore) errors.GraphqlErrorInterface { return nil }
func (n *BaseNode) GetPossibleTypes() map[string]*ObjectNode               { return n.PossibleTypes }

type ObjectNode struct {
	BaseNode
	BaseLocation
	Fields         map[string]*Field `json:"fields"`
	InterfaceNames []string          `json:"-"`
	IsModel        bool              `json:"-"`
	Scopes         []string          `json:"-"`
	Table          string            `json:"-"`
}

func (o *ObjectNode) GetFields() map[string]*Field { return o.Fields }
func (o *ObjectNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
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
				return &errors.GraphQLError{
					Message:   fmt.Sprintf("interface %s not found", interfaceName),
					Locations: []*errors.GraphqlLocation{o.GetLocation()},
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
	BaseLocation
	Fields map[string]*Field `json:"fields"`
}

func (o *InterfaceNode) GetFields() map[string]*Field { return o.Fields }
func (o *InterfaceNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
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
	BaseLocation
	TypeNames map[string]string `json:"-"`
}

func (u *UnionNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
	if err := ValidateDirectives(u.GetName(), u.GetDirectives(), store, LocationUnion); err != nil {
		return err
	}
	for _, typeName := range u.TypeNames {
		typeNode, ok := store.Objects[typeName]
		if !ok {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("type %s not found", typeName),
				Locations: []*errors.GraphqlLocation{u.GetLocation()},
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
	BaseLocation
	EnumValues map[string]*EnumValue `json:"enumValues"`
}

func (e *EnumNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
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
	BaseLocation
	Name              string       `json:"name"`
	Description       *string      `json:"description"`
	Directives        []*Directive `json:"-"`
	Value             int8         `json:"-"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason *string      `json:"deprecationReason"`
}

func (e *EnumValue) Validate(store *NodeStore) errors.GraphqlErrorInterface {
	if err := ValidateDirectives(e.Name, e.Directives, store, LocationEnumValue); err != nil {
		return err
	}
	enum := GetDirective("enum", e.Directives)
	if len(enum) == 1 {
		directive := enum[0]
		value := directive.GetArg("value")
		if value != nil {
			if value.Value == nil {
				return &errors.GraphQLError{
					Message:   "enum value must be a number",
					Locations: []*errors.GraphqlLocation{directive.GetLocation()},
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
	BaseLocation
	Fields map[string]*Field `json:"inputFields"`
}

func (i *InputObjectNode) GetFields() map[string]*Field { return i.Fields }

func (i *InputObjectNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
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
	BaseLocation
	ScalarType ScalarType `json:"-"`
}

func (s *ScalarNode) Validate(store *NodeStore) errors.GraphqlErrorInterface {
	if err := ValidateDirectives(s.GetName(), s.GetDirectives(), store, LocationScalar); err != nil {
		return err
	}
	return nil
}

type ScalarType interface {
	ParseValue(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface)
	Serialize(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface)
	ParseLiteral(v interface{}, location *errors.GraphqlLocation) (interface{}, errors.GraphqlErrorInterface)
	GoType() string
}

type Field struct {
	BaseLocation
	Alias             string               `json:"-"`
	Name              string               `json:"name"`
	Description       *string              `json:"description"`
	Args              map[string]*Argument `json:"args"`
	Type              *TypeRef             `json:"type"`
	IsDeprecated      bool                 `json:"isDeprecated"`
	DeprecationReason *string              `json:"deprecationReason"`

	Children             map[string]*Field    `json:"-"`
	Directives           []*Directive         `json:"-"`
	IsFragment           bool                 `json:"-"`
	IsUnion              bool                 `json:"-"`
	Fragment             *Fragment            `json:"-"`
	IsAttr               bool                 `json:"-"`
	IsSearchable         bool                 `json:"-"`
	DefinitionDirectives []*Directive         `json:"-"`
	DefinitionArgs       map[string]*Argument `json:"-"`
	Relation             *Relation            `json:"-"`

	ServiceInfo *ServiceInfo `json:"-"`
}

type ServiceInfo struct {
	ServiceName string `json:"serviceName"`
}

type RelationType string

const (
	RelationTypeBelongsTo     RelationType = "RelationTypeBelongsTo"
	RelationTypeHasMany       RelationType = "RelationTypeHasMany"
	RelationTypeHasOne        RelationType = "RelationTypeHasOne"
	RelationTypeMorphTo       RelationType = "RelationTypeMorphTo"
	RelationTypeMorphMany     RelationType = "RelationTypeMorphMany"
	RelationTypeBelongsToMany RelationType = "RelationTypeBelongsToMany"
)

type Relation struct {
	RelationType       RelationType `json:"relationType"`
	Name               string       `json:"name"`
	CurrentType        string       `json:"currentType"`
	ForeignKey         string       `json:"foreignKey"`
	Reference          string       `json:"reference"`
	MorphType          string       `json:"morphType"`
	MorphKey           string       `json:"morphKey"`
	PivotForeignKey    string       `json:"pivotForeignKey"`
	PivotReference     string       `json:"pivotReference"`
	RelationForeignKey string       `json:"relationForeignKey"`
	Pivot              string       `json:"pivot"`
}

func (f *Field) Validate(store *NodeStore, objectFields map[string]*Field, objectNode Node, location Location, fragments map[string]*Fragment, args map[string]*Argument) errors.GraphqlErrorInterface {
	if f.IsFragment {
		if fragments == nil {
			return &errors.GraphQLError{
				Message:   "fragments not found",
				Locations: []*errors.GraphqlLocation{f.GetLocation()},
			}
		}
		if fragments[f.Type.Name] == nil {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("fragment %s not found", f.Type.Name),
				Locations: []*errors.GraphqlLocation{f.GetLocation()},
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
			return &errors.GraphQLError{
				Message:   "object node not found",
				Locations: []*errors.GraphqlLocation{f.GetLocation()},
			}
		} else {
			if f.Name == "__typename" {
				f.Type = &TypeRef{
					Kind: KindScalar,
					Name: "String",
				}
				return nil
			}
			if objectNode.GetFields()[f.Name] != nil {
				f.Type = objectNode.GetFields()[f.Name].Type
				realType := f.Type.GetRealType()
				if realType.TypeNode == nil || (realType.TypeNode.GetKind() == KindScalar && realType.TypeNode.(*ScalarNode).ScalarType == nil) {
					return &errors.GraphQLError{
						Message:   fmt.Sprintf("field %s type not found", f.Name),
						Locations: []*errors.GraphqlLocation{f.GetLocation()},
					}
				}
				field := objectNode.GetFields()[f.Name]
				// merge
				f.DefinitionDirectives = append(f.DefinitionDirectives, field.Directives...)
				f.Relation = field.Relation
				f.DefinitionArgs = field.Args
				f.IsAttr = field.IsAttr
				for _, defArg := range f.DefinitionArgs {
					if defArg.DefaultValue != nil && f.Args[defArg.Name] == nil {
						if f.Args == nil {
							f.Args = make(map[string]*Argument)
						}
						f.Args[defArg.Name] = &Argument{
							Name:        defArg.Name,
							Value:       defArg.DefaultValue,
							IsReference: defArg.IsReference,
						}
					}
				}
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
					return &errors.GraphQLError{
						Message:   fmt.Sprintf("field %s not found in function %s", f.Name, "Validate"),
						Locations: []*errors.GraphqlLocation{f.GetLocation()},
					}
				}
			}
		}
	}

	if f.IsUnion {
		if f.Type.TypeNode.GetKind() != KindObject {
			return &errors.GraphQLError{
				Message:   "field type must be a object type",
				Locations: []*errors.GraphqlLocation{f.GetLocation()},
			}
		}
		obj := objectNode.GetPossibleTypes()[f.Type.Name]
		if obj == nil {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("union type %s not found", f.Type.Name),
				Locations: []*errors.GraphqlLocation{f.GetLocation()},
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
	BaseLocation
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

func (t *TypeRef) IsList() bool {
	if t.Kind == KindList {
		return true
	}
	if t.Kind == KindNonNull {
		return t.OfType.IsList()
	}
	return false
}

func (t *TypeRef) IsScalar() bool {
	if t.Kind == KindNonNull {
		return t.OfType.IsScalar()
	}
	if t.Kind == KindList {
		return t.OfType.IsScalar()
	}
	return t.Kind == KindScalar
}

func (t *TypeRef) IsObject() bool {
	if t.Kind == KindNonNull {
		return t.OfType.IsObject()
	}
	if t.Kind == KindList {
		return t.OfType.IsObject()
	}
	return t.Kind == KindObject
}

func (t *TypeRef) IsEnum() bool {
	return t.GetRealType().Kind == KindEnum
}

func (t *TypeRef) IsNullable() bool {
	if t.Kind == KindNonNull {
		return t.OfType.IsNullable()
	}
	return false
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
		if NonNull {
			return t.TypeNode.(*ScalarNode).ScalarType.GoType()
		}
		return "*" + t.TypeNode.(*ScalarNode).ScalarType.GoType()
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
	return "interface{}"
}

func (t *TypeRef) GetGoRealType() string {
	if t == nil {
		return "interface{}"
	}

	switch t.Kind {
	case KindScalar:
		return t.TypeNode.(*ScalarNode).ScalarType.GoType()
	case KindEnum, KindObject, KindInputObject:
		return t.Name
	case KindList:
		return "[]" + t.OfType.GetGoRealType()
	case KindNonNull:
		return t.OfType.GetGoRealType()
	}
	return "interface{}"
}

func (t *TypeRef) Validate(store *NodeStore) errors.GraphqlErrorInterface {
	if t.Kind == KindNonNull {
		if t.OfType == nil {
			return &errors.GraphQLError{
				Message:   "non-null type cannot be null, function: Validate",
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
		}
		if err := t.OfType.Validate(store); err != nil {
			return err
		}
	} else if t.Kind == KindList {
		if t.OfType == nil {
			return &errors.GraphQLError{
				Message:   "list type cannot be nil",
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
		}
		if err := t.OfType.Validate(store); err != nil {
			return err
		}
	} else {
		if t.Name == "" {
			return &errors.GraphQLError{
				Message:   "type name cannot be empty",
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
		}
		if node, ok := store.Nodes[t.Name]; !ok {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("type %s not found", t.Name),
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
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

func (t *TypeRef) ValidateValue(v interface{}, isVariable bool) errors.GraphqlErrorInterface {

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
			return &errors.GraphQLError{
				Message:   "non-null type cannot be null, function: ValidateValue " + t.GetGoName(),
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
		}
		return t.OfType.ValidateValue(v, isVariable)
	case KindUnion:
		return t.validateUnionValue(v)
	default:
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid type: %s", t.Name),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}
}

func (t *TypeRef) validateUnionValue(v interface{}) errors.GraphqlErrorInterface {
	objValue, ok := v.(map[string]interface{})
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("expected object value, got %T", v),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}
	typename, ok := objValue["__typename"].(string)
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("union type must have __typename field, object value %+v， %+v", objValue, t),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	unionNode, ok := t.TypeNode.(*UnionNode)
	if !ok {
		return &errors.GraphQLError{
			Message:   "invalid union type node",
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}
	memberName := utils.UcFirst(utils.CamelCase(typename))
	// Check if typename is a valid union member
	var memberType *ObjectNode
	for _, member := range unionNode.PossibleTypes {
		if member.Name == memberName {
			memberType = member
			break
		}
	}

	if memberType == nil {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("invalid union member type: %s", typename),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	// Validate the value against the member type
	return nil
}

func (t *TypeRef) validateScalarValue(v interface{}, isVariable bool) errors.GraphqlErrorInterface {
	var err errors.GraphqlErrorInterface
	if isVariable {
		_, err = t.TypeNode.(*ScalarNode).ScalarType.ParseValue(v, t.GetLocation())
	} else {
		_, err = t.TypeNode.(*ScalarNode).ScalarType.ParseLiteral(v, t.GetLocation())
	}
	if err != nil {
		return err
	}
	return nil
}

func (t *TypeRef) validateEnumValue(v interface{}) errors.GraphqlErrorInterface {
	strValue, ok := v.(string)
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("expected string value, got %T", v),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	enumNode, ok := t.TypeNode.(*EnumNode)
	if !ok {
		return &errors.GraphQLError{
			Message:   "invalid enum type node",
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	for _, enumValue := range enumNode.EnumValues {
		if enumValue.Name == strValue {
			return nil
		}
	}

	return &errors.GraphQLError{
		Message:   fmt.Sprintf("invalid enum value: %s", strValue),
		Locations: []*errors.GraphqlLocation{t.GetLocation()},
	}
}

func (t *TypeRef) validateObjectValue(v interface{}, isVariable bool) errors.GraphqlErrorInterface {
	objValue, ok := v.(map[string]interface{})
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("expected object value, got %T", v),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	objNode, ok := t.TypeNode.(*ObjectNode)
	if !ok {
		return &errors.GraphQLError{
			Message:   "invalid object type node",
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	for _, field := range objNode.Fields {
		fieldValue, exists := objValue[field.Name]
		if !exists && field.Type.Kind == KindNonNull && !isVariable {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("required field %s is missing", field.Name),
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
		}
		if exists {
			if err := field.Type.ValidateValue(fieldValue, isVariable); err != nil {
				return &errors.GraphQLError{
					Message:   err.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
		}
	}

	return nil
}

func (t *TypeRef) validateInputObjectValue(v interface{}, isVariable bool) errors.GraphqlErrorInterface {
	inputObjValue, ok := v.(map[string]interface{})
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("expected object value, got %T", v),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	inputObjNode, ok := t.TypeNode.(*InputObjectNode)
	if !ok {
		return &errors.GraphQLError{
			Message:   "invalid input object type node",
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	for _, field := range inputObjNode.Fields {
		fieldValue, exists := inputObjValue[field.Name]
		if !exists && field.Type.Kind == KindNonNull {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("required field %s is missing", field.Name),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}
		if exists {
			if err := field.Type.ValidateValue(fieldValue, isVariable); err != nil {
				return &errors.GraphQLError{
					Message:   err.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
		}
	}

	return nil
}

func (t *TypeRef) validateListValue(v interface{}, isVariable bool) errors.GraphqlErrorInterface {
	var list []interface{}
	switch v := v.(type) {
	case []interface{}:
		list = v
	case map[string]interface{}:
		list = []interface{}{v}
	default:
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("expected list value, got %T", v),
			Locations: []*errors.GraphqlLocation{t.GetLocation()},
		}
	}

	for _, item := range list {
		if err := t.OfType.ValidateValue(item, isVariable); err != nil {
			return &errors.GraphQLError{
				Message:   err.Error(),
				Locations: []*errors.GraphqlLocation{t.GetLocation()},
			}
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
	BaseLocation
	Name       string               `json:"name"`
	Args       map[string]*Argument `json:"args"`
	Definition *DirectiveDefinition `json:"-"`
}

func (d *Directive) GetArg(name string) *Argument {
	return d.Args[name]
}

func (d *Directive) Validate(store *NodeStore, location Location) errors.GraphqlErrorInterface {
	d.Definition = store.Directives[d.Name]
	if d.Definition == nil {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("directive %s not found", d.Name),
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
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
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("directive %s is not valid for location %s", d.Name, location),
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
		}
	}

	d.Definition.Directives = append(d.Definition.Directives, d)
	return nil
}

type DirectiveDefinition struct {
	BaseLocation
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

func (d *DirectiveDefinition) Validate(store *NodeStore) errors.GraphqlErrorInterface {
	for _, arg := range d.Args {
		if err := arg.Validate(store, nil, nil); err != nil {
			return err
		}
	}
	if len(d.Locations) == 0 {
		return &errors.GraphQLError{
			Message:   "directive locations cannot be empty",
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
		}
	}
	for _, loc := range d.Locations {
		if !IsValidLocation(loc) {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("invalid location: %s", loc),
				Locations: []*errors.GraphqlLocation{d.GetLocation()},
			}
		}
	}
	for _, directive := range d.Directives {
		for _, arg := range directive.Args {
			defArg := d.Args[arg.Name]
			if defArg == nil {
				return &errors.GraphQLError{
					Message:   fmt.Sprintf("argument %s not found", arg.Name),
					Locations: []*errors.GraphqlLocation{directive.GetLocation()},
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
	BaseLocation
	Name       string            `json:"name"`
	On         string            `json:"on"`
	Object     *ObjectNode       `json:"-"`
	Directives []*Directive      `json:"-"`
	Fields     map[string]*Field `json:"fields"`
}

func (f *Fragment) Validate(store *NodeStore, fragments map[string]*Fragment) errors.GraphqlErrorInterface {
	objectNode, ok := store.Objects[f.On]
	if !ok {
		return &errors.GraphQLError{
			Message:   fmt.Sprintf("fragment %s not found", f.On),
			Locations: []*errors.GraphqlLocation{f.GetLocation()},
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

func ValidateDirectives(name string, directives []*Directive, store *NodeStore, location Location) errors.GraphqlErrorInterface {
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
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("directive %s not found", directiveName),
				Locations: []*errors.GraphqlLocation{{Column: 1, Line: 1}},
			}
		}
		if !directiveDefinition.Repeatable && count > 1 {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("directive %s is not repeatable", directiveName),
				Locations: []*errors.GraphqlLocation{directiveDefinition.GetLocation()},
			}
		}
	}
	return nil
}

type Argument struct {
	BaseLocation
	Name         string       `json:"name"`
	Description  *string      `json:"description"`
	Directives   []*Directive `json:"-"`
	Type         *TypeRef     `json:"type"`
	DefaultValue any          `json:"default_value"`
	Value        any          `json:"-"`
	IsVariable   bool         `json:"-"`
	IsReference  bool         `json:"-"`
}

func (a *Argument) GetValue() (interface{}, errors.GraphqlErrorInterface) {
	var err errors.GraphqlErrorInterface
	if a.Value == nil {
		return nil, nil
	}

	// If type is non-null or list, get the inner type
	t := a.Type
	for t.Kind == KindNonNull || t.Kind == KindList {
		t = t.OfType
	}

	// For scalar and enum types, parse the value directly
	if t.Kind == KindScalar {
		val, err := t.TypeNode.(*ScalarNode).ScalarType.ParseValue(a.Value, a.GetLocation())
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	if t.Kind == KindEnum {
		return a.Value.(string), nil
	}

	// For object and input object types, recursively get values
	if t.Kind == KindObject || t.Kind == KindInputObject {
		if objValue, ok := a.Value.(map[string]interface{}); ok {
			result := make(map[string]interface{})

			var fields map[string]*Field
			if t.Kind == KindObject {
				fields = t.TypeNode.(*ObjectNode).Fields
			} else {
				fields = t.TypeNode.(*InputObjectNode).Fields
			}

			for fieldName, fieldValue := range objValue {
				if field, exists := fields[fieldName]; exists {
					arg := &Argument{
						Type:  field.Type,
						Value: fieldValue,
					}
					result[fieldName], err = arg.GetValue()
					if err != nil {
						return nil, err
					}
				}
			}
			return result, nil
		}
	}

	return a.Value, nil
}

func (a *Argument) GetDefaultValue() *string {
	if a.DefaultValue == nil {
		return nil
	}
	str := fmt.Sprintf("%v", a.DefaultValue)
	return &str
}

func (a *Argument) Validate(store *NodeStore, args map[string]*Argument, field *Field) errors.GraphqlErrorInterface {

	location := LocationArgumentDefinition
	if a.IsVariable {
		location = LocationVariableDefinition
	}
	if a.IsReference {
		if args == nil {
			return &errors.GraphQLError{
				Message:   "variable argument not found",
				Locations: []*errors.GraphqlLocation{a.GetLocation()},
			}
		}

		name, ok := a.Value.(string)

		if !ok {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("expected string value, got %T", a.Value),
				Locations: []*errors.GraphqlLocation{a.GetLocation()},
			}
		}
		if args[name] == nil {
			return &errors.GraphQLError{
				Message:   "variable argument not found",
				Locations: []*errors.GraphqlLocation{a.GetLocation()},
			}
		}
		a.Type = args[name].Type
		a.Value = args[name].Value
	}
	if a.Type == nil {
		if field == nil {
			return &errors.GraphQLError{
				Message:   "argument type not found",
				Locations: []*errors.GraphqlLocation{a.GetLocation()},
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

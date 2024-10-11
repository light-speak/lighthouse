package ast

type Node interface {
	GetName() string
	GetNodeType() NodeType
	GetDescription() string
	IsDeprecated() (bool, string)
	GetField(name string) *FieldNode
	GetDirective(name string) *DirectiveNode
	GetArg(name string) *ArgumentNode
	GetParent() Node

	GetDirectives() []*DirectiveNode
	GetDirectivesByName(name string) []*DirectiveNode
	GetArgs() []*ArgumentNode
	GetFields() []*FieldNode

	GoType() string
}

type ScalarType interface {
	ParseValue(v string) (interface{}, error)
	Serialize(v interface{}) (string, error)
	ParseLiteral(v Value) (interface{}, error)
	GoType() string
}

type NodeType string

const (
	NodeTypeOperation           NodeType = "Operation"
	NodeTypeSubOperation        NodeType = "SubOperation"
	NodeTypeType                NodeType = "Type"
	NodeTypeField               NodeType = "Field"
	NodeTypeArgument            NodeType = "Argument"
	NodeTypeDirective           NodeType = "Directive"
	NodeTypeDirectiveDefinition NodeType = "DirectiveDefinition"
	NodeTypeScalar              NodeType = "Scalar"
	NodeTypeUnion               NodeType = "Union"
	NodeTypeEnum                NodeType = "Enum"
	NodeTypeInterface           NodeType = "Interface"
	NodeTypeInput               NodeType = "Input"
	NodeTypeEnumValue           NodeType = "EnumValue"
	NodeTypeFragment            NodeType = "Fragment"
)

type TypeCategory string

const (
	TypeCategoryScalar TypeCategory = "Scalar"
	TypeCategoryEnum   TypeCategory = "Enum"
	TypeCategoryInput  TypeCategory = "Input"
	TypeCategoryUnion  TypeCategory = "Union"
	TypeCategoryType   TypeCategory = "Type"
	TypeFragmentType   TypeCategory = "Fragment"
)

// FieldType represents the type of a field in the GraphQL schema
type FieldType struct {
	Name         string       // Name of the field type
	Type         Node         // The underlying type node
	TypeCategory TypeCategory // Category of the type (e.g., Scalar, Enum, etc.)
	IsList       bool         // Indicates if the field is a list type
	ElemType     *FieldType   // Element type if IsList is true
	IsNonNull    bool         // Indicates if the field is non-nullable
	Level        int          // Level of nesting for list types
}

func (f *FieldType) GoType() string {
	if f.IsList {
		return f.listGoType()
	}
	return f.singleGoType()
}

func (f *FieldType) listGoType() string {
	elemType := f.ElemType.GoType()
	if !f.IsNonNull {
		return "*[]" + elemType
	}
	return "[]" + elemType
}

func (f *FieldType) singleGoType() string {
	var baseType string
	if scalar, ok := f.Type.(*ScalarNode); ok {
		baseType = scalar.Scalar.GoType()
	} else {
		baseType = f.Type.GoType()
	}

	if !f.IsNonNull {
		return "*" + baseType
	}
	return baseType
}

// ArgumentValue represents a value for an argument
// it can also contain children values
type ArgumentValue struct {
	Value    Value
	Type     *FieldType
	Children []*ArgumentValue
}

type BaseNode struct {
	Name        string
	Description string
	Directives  []*DirectiveNode
}

// GetName returns the name of the node
func (b *BaseNode) GetName() string {
	return b.Name
}

// GetDescription returns the description of the node
func (b *BaseNode) GetDescription() string {
	return b.Description
}

// GetDirective returns a directive by name
func (b *BaseNode) GetDirective(name string) *DirectiveNode {
	for _, directive := range b.Directives {
		if directive.Name == name {
			return directive
		}
	}
	return nil
}

// GetDirectivesByName returns all directives by name
func (b *BaseNode) GetDirectivesByName(name string) []*DirectiveNode {
	var directives []*DirectiveNode
	for _, directive := range b.Directives {
		if directive.Name == name {
			directives = append(directives, directive)
		}
	}
	return directives
}

// GetDirectives returns all directives
func (b *BaseNode) GetDirectives() []*DirectiveNode {
	return b.Directives
}

// IsDeprecated checks if the node is deprecated
func (b *BaseNode) IsDeprecated() (bool, string) {
	directive := b.GetDirective("deprecated")
	if directive == nil {
		return false, ""
	}
	reason, err := ExtractValue(directive.GetArg("reason").Value.Value)
	if err != nil {
		return false, ""
	}
	return true, reason.(string)
}

// GoType returns the Go type of the node
func (b *BaseNode) GoType() string { return b.GetName() }

// Common methods for all node types
func (b *BaseNode) GetArg(name string) *ArgumentNode { return nil }
func (b *BaseNode) GetParent() Node                  { return nil }
func (b *BaseNode) GetArgs() []*ArgumentNode         { return nil }
func (b *BaseNode) GetFields() []*FieldNode          { return nil }
func (b *BaseNode) GetField(name string) *FieldNode  { return nil }

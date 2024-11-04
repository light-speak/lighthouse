package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/utils"
)

// parseObject parse a object node
//
//	type User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseObject() {
	extend := false
	if p.PreviousToken().Type == lexer.Extend {
		extend = true
	}

	object := &ast.ObjectNode{
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Type),
		},
		InterfaceNames: p.parseImplements(),
	}

	object.Directives = p.parseDirectives()

	if p.currToken.Type != lexer.LeftBrace {
		object.Fields = make(map[string]*ast.Field)
		p.AddObject(object, extend)
		return
	}
	p.expect(lexer.LeftBrace)

	fields := make(map[string]*ast.Field)
	for {
		field := p.parseField(false, "")
		if field == nil {
			continue
		}
		if _, ok := fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		fields[field.Name] = field
		if p.currToken.Type == lexer.RightBrace {
			break
		}
	}
	object.Fields = fields
	p.AddObject(object, extend)
}

// parseDescription parses a description if present
func (p *Parser) parseDescription() *string {
	if p.PreviousToken().Type == lexer.Message {
		description := strings.Split(p.PreviousToken().Value, "\"")[1]
		return &description
	}
	return nil
}

// parseImplements parses implements if present
func (p *Parser) parseImplements() []string {
	var implements []string
	if p.currToken.Type != lexer.Implements {
		return implements
	}
	p.expect(lexer.Implements)
	for {
		implementName := p.currToken.Value
		if utils.Contains(implements, implementName) {
			panic("duplicate implement: " + implementName)
		}
		implements = append(implements, implementName)
		p.nextToken()
		if p.currToken.Type != lexer.And {
			break
		}
		p.expect(lexer.And)
	}
	return implements
}

// parseDirectives parses directives if present
func (p *Parser) parseDirectives() []*ast.Directive {
	var directives []*ast.Directive
	if p.currToken.Type != lexer.At {
		return directives
	}
	for {
		directives = append(directives, p.parseDirective())
		if p.currToken.Type != lexer.At {
			break
		}
	}

	return directives
}

// parseDirective parses a directive if present
func (p *Parser) parseDirective() *ast.Directive {
	directive := &ast.Directive{
		Name: p.expectAndGetValue(lexer.At),
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	var args map[string]*ast.Argument
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments()
	}
	directive.Args = args
	return directive
}

// parseField parses a field node
// Examples:
// "It is ID"
// id: ID!
// name: String!
// age: Int
// email: String
// createdAt: DateTime
func (p *Parser) parseField(isOperation bool, alias string) *ast.Field {
	// Handle comments
	if p.currToken.Type == lexer.Comment || p.currToken.Type == lexer.Message {
		p.nextToken()
		return nil
	}

	// Handle fragment spreads
	if p.currToken.Type == lexer.TripleDot {
		p.expect(lexer.TripleDot)

		// Handle union types
		if p.currToken.Type == lexer.On {
			p.expect(lexer.On)
			field := &ast.Field{
				Name:    p.currToken.Value,
				Alias:   alias,
				IsUnion: true,
				Type: &ast.TypeRef{
					Name: p.currToken.Value,
					BaseLocation: ast.BaseLocation{
						Line:   p.currToken.Line,
						Column: p.currToken.LinePosition,
					},
				},
				BaseLocation: ast.BaseLocation{
					Line:   p.currToken.Line,
					Column: p.currToken.LinePosition,
				},
			}
			p.nextToken()
			p.expect(lexer.LeftBrace)
			field.Children = make(map[string]*ast.Field)
			for p.currToken.Type != lexer.RightBrace {
				cField := p.parseField(isOperation, "")
				if cField != nil {
					field.Children[cField.Name] = cField
				}
			}
			p.expect(lexer.RightBrace)
			return field
		}

		// Handle fragments
		fragmentField := &ast.Field{
			Name:       p.currToken.Value,
			IsFragment: true,
			Type: &ast.TypeRef{
				Name: p.currToken.Value,
				BaseLocation: ast.BaseLocation{
					Line:   p.currToken.Line,
					Column: p.currToken.LinePosition,
				},
			},
			BaseLocation: ast.BaseLocation{
				Line:   p.currToken.Line,
				Column: p.currToken.LinePosition,
			},
		}
		p.nextToken()
		return fragmentField
	}

	field := &ast.Field{
		Name:        p.currToken.Value,
		Alias:       alias,
		Description: p.parseDescription(),
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}
	p.nextToken()

	if isOperation && p.currToken.Type == lexer.Colon {
		alias = field.Name
		p.expect(lexer.Colon)
		return p.parseField(isOperation, alias)
	}

	if p.currToken.Type == lexer.LeftParent {
		field.Args = p.parseArguments()
	}

	switch p.currToken.Type {
	case lexer.Colon:
		p.expect(lexer.Colon)
		field.Type, _ = p.parseTypeReferenceAndValue()
		field.Directives = p.parseDirectives()
	case lexer.LeftBrace:
		field.Children = make(map[string]*ast.Field)
		p.expect(lexer.LeftBrace)
		for p.currToken.Type != lexer.RightBrace {
			childField := p.parseField(isOperation, "")
			if childField != nil {
				field.Children[childField.Name] = childField
			}
		}
		p.expect(lexer.RightBrace)
	}

	// Skip any trailing comments
	for p.currToken.Type == lexer.Comment || p.currToken.Type == lexer.Message {
		p.nextToken()
	}

	return field
}

// parseArguments parse arguments
// (id: ID!, name: String!)
func (p *Parser) parseArguments() map[string]*ast.Argument {
	args := make(map[string]*ast.Argument)
	if p.currToken.Type != lexer.LeftParent {
		return args
	}
	p.expect(lexer.LeftParent)
	for p.currToken.Type != lexer.RightParent {
		arg := p.parseArgument()
		args[arg.Name] = arg
		if p.currToken.Type != lexer.RightParent {
			p.expect(lexer.Comma)
		}
	}

	p.expect(lexer.RightParent)
	return args
}

func (p *Parser) parseDefaultValue() any {
	if p.currToken.Type == lexer.Equal {
		p.expect(lexer.Equal) // skip =
		return p.parseValue()
	}
	return nil
}

// parseTypeReference parses a type reference
// Examples: ID, String, [Int], [[Int]], [User], [[User]], [User!], [1,2,3], [[1,2,3],[4,6]]
func (p *Parser) parseTypeReferenceAndValue() (*ast.TypeRef, any) {
	var fieldType *ast.TypeRef
	var value any

	if p.currToken.Type == lexer.LeftBracket {
		p.expect(lexer.LeftBracket)
		innerType, innerValue := p.parseTypeReferenceAndValue()

		// 处理多重数组
		if p.currToken.Type == lexer.Comma {
			value = []any{innerValue}
			for p.currToken.Type == lexer.Comma {
				p.expect(lexer.Comma)
				_, nextValue := p.parseTypeReferenceAndValue()
				value = append(value.([]any), nextValue)
			}
		} else if innerValue != nil {
			// 单层数组值
			value = []any{innerValue}
		} else {
			// 处理数组类型
			fieldType = &ast.TypeRef{
				Kind:   ast.KindList,
				OfType: innerType,
				BaseLocation: ast.BaseLocation{
					Line:   p.currToken.Line,
					Column: p.currToken.LinePosition,
				},
			}
		}

		p.expect(lexer.RightBracket)
	} else {
		switch p.currToken.Type {
		case lexer.IntNumber, lexer.FloatNumber, lexer.Message, lexer.Boolean, lexer.Null:
			value = p.parseValue()
		case lexer.Letter:
			fieldType = &ast.TypeRef{
				Kind: "",
				Name: p.currToken.Value,
				BaseLocation: ast.BaseLocation{
					Line:   p.currToken.Line,
					Column: p.currToken.LinePosition,
				},
			}
			p.expect(lexer.Letter)
		default:
			panic(fmt.Sprintf("Unexpected token type in type reference parsing: %v", p.currToken.Type))
		}
	}

	if fieldType != nil && p.currToken.Type == lexer.Exclamation {
		p.expect(lexer.Exclamation)
		fieldType = &ast.TypeRef{
			Kind:   ast.KindNonNull,
			OfType: fieldType,
			BaseLocation: ast.BaseLocation{
				Line:   p.currToken.Line,
				Column: p.currToken.LinePosition,
			},
		}
	}

	return fieldType, value
}

// parseArgument parse an argument node
func (p *Parser) parseArgument() *ast.Argument {
	description := p.parseDescription()
	isVariable := false
	isReference := false
	if p.currToken.Type == lexer.Variable {
		isVariable = true
	}
	name := p.currToken.Value

	p.nextToken()         // skip name
	p.expect(lexer.Colon) // skip :

	var typeRef *ast.TypeRef
	var defaultValue, value any

	switch p.currToken.Type {
	case lexer.Letter, lexer.LeftBracket:
		// Case 1: Normal parameter with type (id: ID!, name: String!)
		// Case 2: Normal parameter with type and default value (id: ID = 123, name: String = "123")
		typeRef, value = p.parseTypeReferenceAndValue()
		defaultValue = p.parseDefaultValue()
	case lexer.IntNumber, lexer.FloatNumber, lexer.Message, lexer.Boolean, lexer.Null, lexer.LeftBrace:
		// Case 3: Normal parameter with value (id: 123, name: "123")
		value = p.parseValue()
	case lexer.Variable:
		// Case 4: Normal parameter with variable (id: $id)
		value = p.currToken.Value
		isReference = true
		p.expect(lexer.Variable)
	default:
		panic("Unexpected token type in argument parsing: " + p.currToken.Value)
	}

	directives := p.parseDirectives()

	return &ast.Argument{
		Name:         name,
		Description:  description,
		Directives:   directives,
		Type:         typeRef,
		Value:        value,
		DefaultValue: defaultValue,
		IsVariable:   !isReference && isVariable,
		IsReference:  isReference,
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}
}

func (p *Parser) parseValue() any {
	if p.currToken.Type == lexer.Null {
		p.expect(lexer.Null)
		return nil
	}
	if p.currToken.Type == lexer.Message {
		value := strings.Trim(p.currToken.Value, "\"")
		p.expect(lexer.Message)
		return value
	}
	if p.currToken.Type == lexer.IntNumber {
		v, err := strconv.ParseInt(p.currToken.Value, 10, 64)
		p.expect(lexer.IntNumber)
		if err != nil {
			panic("invalid integer value: " + err.Error())
		}
		return v
	}
	if p.currToken.Type == lexer.FloatNumber {
		v, err := strconv.ParseFloat(p.currToken.Value, 64)
		if err != nil {
			panic("invalid float value: " + err.Error())
		}
		p.expect(lexer.FloatNumber)
		return v
	}
	if p.currToken.Type == lexer.Boolean {
		value := p.currToken.Value == "true"
		p.expect(lexer.Boolean)
		return value
	}
	if p.currToken.Type == lexer.Letter {
		value := p.currToken.Value
		p.expect(lexer.Letter)
		return value
	}
	if p.currToken.Type == lexer.LeftBracket {
		p.expect(lexer.LeftBracket)
		var values []any
		for p.currToken.Type != lexer.RightBracket {
			values = append(values, p.parseValue())
			if p.currToken.Type == lexer.Comma {
				p.expect(lexer.Comma)
			}
		}
		p.expect(lexer.RightBracket)
		return values
	}
	if p.currToken.Type == lexer.LeftBrace {
		p.expect(lexer.LeftBrace)
		values := make(map[string]any)
		for p.currToken.Type != lexer.RightBrace {
			key := p.currToken.Value
			p.expect(lexer.Letter)
			p.expect(lexer.Colon)
			values[key] = p.parseValue()
			if p.currToken.Type == lexer.Comma {
				p.expect(lexer.Comma)
			}
		}
		p.expect(lexer.RightBrace)
		return values
	}
	return nil
}

// parseDirectiveDefinition parses a directive definition node
// Example: directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
func (p *Parser) parseDirectiveDefinition() {
	description := p.parseDescription()
	p.expect(lexer.Directive)
	node := &ast.DirectiveDefinition{
		Name:        p.expectAndGetValue(lexer.At),
		Description: description,
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	if p.currToken.Type == lexer.LeftParent {
		node.Args = p.parseArguments()
	}
	if p.currToken.Type == lexer.Repeatable {
		node.Repeatable = true
		p.expect(lexer.Repeatable)
	}
	p.expect(lexer.On)
	node.Locations = p.parseLocations()

	p.AddDirectiveDefinition(node)
}

// parseLocations parse locations
// ON FIELD_DEFINITION | ARGUMENT_DEFINITION | INTERFACE | UNION | ENUM | INPUT_OBJECT | SCALAR | OBJECT
func (p *Parser) parseLocations() []ast.Location {
	locations := make([]ast.Location, 0)

	for {
		locations = append(locations, ast.Location(p.currToken.Value))
		p.nextToken()
		if p.currToken.Type != lexer.Pipe {
			break
		}
		p.expect(lexer.Pipe)
	}

	return locations
}

// parseEnum parse an enum node
//
//	enum Role {
//	  ADMIN
//	  USER
//	  GUEST
//	}
func (p *Parser) parseEnum() {
	extend := false
	if p.PreviousToken().Type == lexer.Extend {
		extend = true
	}

	node := &ast.EnumNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Enum),
			Directives:  p.parseDirectives(),
		},
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	if p.currToken.Type != lexer.LeftBrace {
		p.AddEnum(node, extend)
		return
	}

	p.expect(lexer.LeftBrace)
	node.EnumValues = make(map[string]*ast.EnumValue)
	for p.currToken.Type != lexer.RightBrace {
		enumValue := p.parseEnumValue()
		if _, ok := node.EnumValues[enumValue.Name]; ok {
			panic("duplicate enum value: " + enumValue.Name)
		}
		node.EnumValues[enumValue.Name] = enumValue
	}

	p.AddEnum(node, extend)
}

func (p *Parser) parseEnumValue() *ast.EnumValue {
	description := p.parseDescription()
	name := p.currToken.Value
	p.nextToken()
	directives := p.parseDirectives()

	return &ast.EnumValue{
		Name:        name,
		Description: description,
		Directives:  directives,
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}
}

// parseExtend parse an extend node
//
//	extend type User {
//	  role: Role!
//	}
func (p *Parser) parseExtend() {
	p.parseDescription()   // Skip extend description
	p.expect(lexer.Extend) // Skip 'extend'
}

// parseInput 解析输入节点
//
//	input User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseInput() {
	extend := false
	if p.PreviousToken().Type == lexer.Extend {
		extend = true
	}

	node := &ast.InputObjectNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Input),
			Directives:  p.parseDirectives(),
		},
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	if p.currToken.Type != lexer.LeftBrace {
		p.AddInput(node, extend)
		return
	}

	p.expect(lexer.LeftBrace)
	node.Fields = make(map[string]*ast.Field)

	for {
		field := p.parseField(false, "")
		if _, ok := node.Fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		node.Fields[field.Name] = field
		if p.currToken.Type == lexer.RightBrace {
			break
		}
	}

	p.AddInput(node, extend)
}

// parseInterface parse an interface node
//
//	interface Node {
//	  id: ID!
//	}
func (p *Parser) parseInterface() {
	extend := false
	if p.PreviousToken().Type == lexer.Extend {
		extend = true
	}

	node := &ast.InterfaceNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Interface),
		},
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	node.Directives = p.parseDirectives()
	if p.currToken.Type != lexer.LeftBrace {
		p.AddInterface(node, extend)
		return
	}
	p.expect(lexer.LeftBrace)
	fields := make(map[string]*ast.Field)
	for p.currToken.Type == lexer.Letter || p.currToken.Type == lexer.TripleDot {
		field := p.parseField(false, "")
		if _, ok := fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		fields[field.Name] = field
	}

	node.Fields = fields

	p.AddInterface(node, extend)
}

// parseScalar parses a scalar node
// Example: scalar DateTime
func (p *Parser) parseScalar() {
	node := &ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Scalar),
			Directives:  p.parseDirectives(),
		},
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	p.AddScalar(node)
}

// parseUnion parses a union node
// Example: union User = Product | Order
func (p *Parser) parseUnion() {
	extend := false
	if p.PreviousToken().Type == lexer.Extend {
		extend = true
	}

	node := &ast.UnionNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Union),
			Directives:  p.parseDirectives(),
		},
		BaseLocation: ast.BaseLocation{
			Line:   p.currToken.Line,
			Column: p.currToken.LinePosition,
		},
	}

	if p.currToken.Type != lexer.Equal {
		p.AddUnion(node, extend)
		return
	}

	p.expect(lexer.Equal)
	node.TypeNames = p.parseUnionTypes()
	p.AddUnion(node, extend)
}

// parseUnionTypes parses the types in a union definition
func (p *Parser) parseUnionTypes() map[string]string {
	types := make(map[string]string)
	for {
		types[p.currToken.Value] = p.currToken.Value
		p.nextToken()
		if p.currToken.Type != lexer.Pipe {
			break
		}
		p.expect(lexer.Pipe)
	}
	return types
}

// expectAndGetValue expects a token type and returns its value
func (p *Parser) expectAndGetValue(tokenType lexer.TokenType) string {
	p.expect(tokenType)
	value := p.currToken.Value
	p.nextToken()
	return value
}

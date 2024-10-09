package parser

import (
	"strconv"
	"strings"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
)

// parseType parse a type node
//
//	type User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseType(extends ...bool) {
	node := &ast.TypeNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Type),
		},
		Implements: p.parseImplements(),
	}
	node.Directives = p.parseDirectives()
	p.expect(lexer.LeftBrace)
	var fields []*ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		field := p.parseField(node)
		fields = append(fields, field)
	}
	p.expect(lexer.RightBrace)
	node.Fields = fields

	p.AddType(node.GetName(), node, len(extends) > 0)
}

// parseDescription parses a description if present
func (p *Parser) parseDescription() string {
	if p.PreviousToken().Type == lexer.Message {
		description := strings.Split(p.PreviousToken().Value, "\"")[1]
		return description
	}
	return ""
}

// parseImplements parses implements if present
func (p *Parser) parseImplements() []string {
	var implements []string
	if p.currToken.Type != lexer.Implements {
		return implements
	}
	p.expect(lexer.Implements)
	for {
		implements = append(implements, p.currToken.Value)
		p.nextToken()
		if p.currToken.Type == lexer.LeftBrace {
			break
		}
		p.expect(lexer.And)
	}
	return implements
}

// parseDirectives parses directives if present
func (p *Parser) parseDirectives() []*ast.DirectiveNode {
	var directives []*ast.DirectiveNode
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
func (p *Parser) parseDirective() *ast.DirectiveNode {
	directive := &ast.DirectiveNode{
		BaseNode: ast.BaseNode{
			Name: p.expectAndGetValue(lexer.At),
		},
	}

	var args []*ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments(directive)
	}
	directive.Args = args
	return directive
}

// parseField parse a field node
// "It is ID"
// id: ID!
// name: String!
// age: Int
// email: String
// createdAt: DateTime
func (p *Parser) parseField(parent ast.Node) *ast.FieldNode {
	field := &ast.FieldNode{
		BaseNode: ast.BaseNode{
			Name:        p.currToken.Value,
			Description: p.parseDescription(),
		},
		Parent: parent,
	}
	p.nextToken()

	if p.currToken.Type == lexer.LeftParent {
		field.Args = p.parseArguments(field)
	}

	switch p.currToken.Type {
	case lexer.Colon:
		p.expect(lexer.Colon)
		field.Type = p.parseTypeReference()
		field.Directives = p.parseDirectives()
	case lexer.LeftBrace:
		field.Children = p.parseChildren(field)
	}

	return field
}

// parseChildren parse children
func (p *Parser) parseChildren(parent ast.Node) []*ast.FieldNode {
	var children []*ast.FieldNode
	p.expect(lexer.LeftBrace)
	for p.currToken.Type != lexer.RightBrace {
		children = append(children, p.parseField(parent))
	}
	p.expect(lexer.RightBrace)
	return children
}

// parseArguments parse arguments
// (id: ID!, name: String!)
func (p *Parser) parseArguments(parent ast.Node) []*ast.ArgumentNode {
	var args []*ast.ArgumentNode
	if p.currToken.Type != lexer.LeftParent {
		return args
	}
	p.expect(lexer.LeftParent)
	for p.currToken.Type != lexer.RightParent {
		args = append(args, p.parseArgument(parent))
		if p.currToken.Type != lexer.RightParent {
			p.expect(lexer.Comma)
		}
	}
	p.expect(lexer.RightParent)
	return args
}

// parseDefaultValue parse a default value
// = 123, = "123", = true, = false, = null
// = [123, 456] , = ["123", "456"]
func (p *Parser) parseDefaultValue() *ast.ArgumentValue {
	if p.currToken.Type == lexer.Equal {
		p.expect(lexer.Equal) // skip =
		return p.parseArgumentValue()
	}
	return nil
}

// parseArgument parse an argument node
func (p *Parser) parseArgument(parent ast.Node) *ast.ArgumentNode {
	description := p.parseDescription()
	name := p.currToken.Value

	p.nextToken()         // skip name
	p.expect(lexer.Colon) // skip :

	var fieldType *ast.FieldType
	var defaultValue, value *ast.ArgumentValue

	if parent.GetNodeType() == ast.NodeTypeDirective {
		// Assigned when using @directive
		value = p.parseArgumentValue()
	} else {
		fieldType = p.parseTypeReference()   // parse type reference
		defaultValue = p.parseDefaultValue() // parse default value
	}

	directives := p.parseDirectives()

	return &ast.ArgumentNode{
		BaseNode: ast.BaseNode{
			Name:        name,
			Description: description,
			Directives:  directives,
		},
		Type:         fieldType,
		Value:        value,
		DefaultValue: defaultValue,
		Parent:       parent,
	}
}

// parseArgumentValue parse a directive argument value
// @directive(arg: "123")
// @directive(arg: 123)
// @directive(arg: [123, 456])
// @directive(arg: ["123", "456"])
// @directive(arg: true, arg2: false)
// @directive(arg: Boolean, arg2: String, arg3: Int, arg4: [[User]!]!, arg5: ID)
// The colon has been parsed in the previous step, so only the value needs to be parsed here
func (p *Parser) parseArgumentValue() *ast.ArgumentValue {
	switch p.currToken.Type {
	case lexer.LeftBracket:
		return p.parseListArgumentValue()
	default:
		return p.parseSingleArgumentValue()
	}
}

func (p *Parser) parseListArgumentValue() *ast.ArgumentValue {
	p.expect(lexer.LeftBracket)
	values := []*ast.ArgumentValue{}

	for p.currToken.Type != lexer.RightBracket {
		values = append(values, p.parseArgumentValue())
		if p.currToken.Type == lexer.Comma {
			p.expect(lexer.Comma)
		}
	}
	p.expect(lexer.RightBracket)

	argValue := &ast.ArgumentValue{
		Children: values,
		Type: &ast.FieldType{
			Name:   "List",
			IsList: true,
		},
	}

	if p.currToken.Type == lexer.Exclamation {
		argValue.Type.IsNonNull = true
		p.expect(lexer.Exclamation)
	}

	return argValue
}

func (p *Parser) parseSingleArgumentValue() *ast.ArgumentValue {
	var v ast.Value
	var typeName string

	switch p.currToken.Type {
	case lexer.Letter:
		v = &ast.StringValue{Value: p.currToken.Value}
		typeName = "String"
	case lexer.IntNumber:
		intValue, err := strconv.ParseInt(p.currToken.Value, 10, 64)
		if err != nil {
			panic("invalid integer value: " + err.Error())
		}
		v = &ast.IntValue{Value: intValue}
		typeName = "Int"
	case lexer.Boolean:
		boolValue := p.currToken.Value == "true"
		v = &ast.BooleanValue{Value: boolValue}
		typeName = "Boolean"
	default:
		panic("unsupported token type: " + p.currToken.Type)
	}

	argValue := &ast.ArgumentValue{
		Value: v,
		Type: &ast.FieldType{
			Name: typeName,
		},
	}

	p.nextToken()

	if p.currToken.Type == lexer.Exclamation {
		argValue.Type.IsNonNull = true
		p.expect(lexer.Exclamation)
	}

	return argValue
}

// parseDirectiveDefinition parses a directive definition node
// Example: directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
func (p *Parser) parseDirectiveDefinition() {
	description := p.parseDescription()
	p.expect(lexer.Directive)
	node := &ast.DirectiveDefinitionNode{
		BaseNode: ast.BaseNode{
			Name:        p.expectAndGetValue(lexer.At),
			Description: description,
		},
	}

	if p.currToken.Type == lexer.LeftParent {
		node.Args = p.parseArguments(node)
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
	node := &ast.EnumNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Enum),
			Directives:  p.parseDirectives(),
		},
	}

	p.expect(lexer.LeftBrace)

	for p.currToken.Type != lexer.RightBrace {
		node.Values = append(node.Values, p.parseEnumValue(node))
	}

	p.expect(lexer.RightBrace)

	p.AddEnum(node)
}

func (p *Parser) parseEnumValue(parent ast.Node) *ast.EnumValueNode {
	description := p.parseDescription()
	name := p.currToken.Value
	p.nextToken()
	directives := p.parseDirectives()

	return &ast.EnumValueNode{
		BaseNode: ast.BaseNode{
			Name:        name,
			Description: description,
			Directives:  directives,
		},
		Parent: parent,
	}
}

// parseExtend parse an extend node
//
//	extend type User {
//	  role: Role!
//	}
func (p *Parser) parseExtend() {
	p.parseDescription() // Skip extend description

	p.nextToken()               // Skip 'extend'
	p.expect(lexer.Type, false) // Ensure the next token is 'type', but not move to next token, continue parsing

	// Parse the extended type using parseType
	p.parseType(true)
}

func (p *Parser) parseFragment() {
	node := &ast.FragmentNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Fragment),
			Directives:  p.parseDirectives(),
		},
		On: p.expectAndGetValue(lexer.On),
	}

	p.expect(lexer.LeftBrace)
	for p.currToken.Type != lexer.RightBrace {
		node.Fields = append(node.Fields, p.parseField(node))
	}

	p.AddFragment(node)
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
	node := &ast.InputNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Input),
			Directives:  p.parseDirectives(),
		},
	}

	p.expect(lexer.LeftBrace)

	for p.currToken.Type != lexer.RightBrace {
		node.Fields = append(node.Fields, p.parseField(node))
	}

	p.AddInput(node)
}

// parseInterface parse an interface node
//
//	interface Node {
//	  id: ID!
//	}
func (p *Parser) parseInterface() {
	node := &ast.InterfaceNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Interface),
		},
	}

	p.expect(lexer.LeftBrace)

	var fields []*ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		fields = append(fields, p.parseField(node))
	}

	node.Fields = fields

	p.AddInterface(node)
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
	}

	p.AddScalar(node)
}

// parseTypeReference parses a type reference
// Examples: ID, String, [Int], [[Int]], [User], [[User]], [User!]
func (p *Parser) parseTypeReference() *ast.FieldType {
	var fieldType *ast.FieldType

	if p.currToken.Type == lexer.LeftBracket {
		p.expect(lexer.LeftBracket)
		elemType := p.parseTypeReference()
		p.expect(lexer.RightBracket)
		fieldType = &ast.FieldType{
			Name:     "List",
			IsList:   true,
			ElemType: elemType,
		}
	} else {
		fieldType = &ast.FieldType{
			Name: p.currToken.Value,
		}
		p.expect(lexer.Letter)
	}

	if p.currToken.Type == lexer.Exclamation {
		fieldType.IsNonNull = true
		p.expect(lexer.Exclamation)
	}

	return fieldType
}

// parseUnion parses a union node
// Example: union User = Product | Order
func (p *Parser) parseUnion() {
	node := &ast.UnionNode{
		BaseNode: ast.BaseNode{
			Description: p.parseDescription(),
			Name:        p.expectAndGetValue(lexer.Union),
			Directives:  p.parseDirectives(),
		},
	}

	p.expect(lexer.Equal)
	node.Types = p.parseUnionTypes()
	p.AddUnion(node)
}

// parseUnionTypes parses the types in a union definition
func (p *Parser) parseUnionTypes() []string {
	var types []string
	for {
		types = append(types, p.currToken.Value)
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
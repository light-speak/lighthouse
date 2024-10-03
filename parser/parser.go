package parser

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

// Parser is responsible for parsing the GraphQL schema.
// It contains a lexer for tokenizing the input and the current token being processed.
// The various maps are used to store different types of AST nodes for quick lookup and management during parsing.
type Parser struct {
	// lexer is the Lexer instance used for lexical analysis, converting the input GraphQL text into a stream of tokens.
	lexer *lexer.Lexer

	// currToken is the current token being processed, which helps the parser determine its state.
	currToken *lexer.Token

	// typeMap, enumMap, scalarMap, unionMap, inputMap, interfaceMap, and directiveMap are all maps
	// that store parsed AST nodes. The keys are the names of the respective types, enums, scalars, unions,
	// input types, interfaces, and directives, while the values are pointers to their corresponding AST node structures.
	typeMap      map[string]*ast.TypeNode
	enumMap      map[string]*ast.EnumNode
	scalarMap    map[string]*ast.ScalarNode
	unionMap     map[string]*ast.UnionNode
	inputMap     map[string]*ast.InputNode
	interfaceMap map[string]*ast.InterfaceNode
	directiveMap map[string]*ast.DirectiveDefinitionNode
}

// ReadGraphQLFile read graphql file and return a lexer
func ReadGraphQLFile(path string) (*lexer.Lexer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return lexer.NewLexer(string(content)), nil
}

// NewParser create a new parser
func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken() // Initialize currToken
	return p
}

// nextToken move to next token
func (p *Parser) nextToken() {
	p.currToken = p.lexer.NextToken()
}

// ParseSchema parse schema
// return a list of ast nodes
// the nodes is a list of type, enum, interface, input, scalar, union, directive, extend
func (p *Parser) ParseSchema() map[string]ast.Node {
	nodes := make(map[string]ast.Node)

	// 定义类型映射
	tokenTypeToParseFunc := map[lexer.TokenType]func() ast.Node{
		lexer.Type:      func() ast.Node { return p.parseType() },
		lexer.Extend:    func() ast.Node { return p.parseExtend() },
		lexer.Enum:      func() ast.Node { return p.parseEnum() },
		lexer.Interface: func() ast.Node { return p.parseInterface() },
		lexer.Input:     func() ast.Node { return p.parseInput() },
		lexer.Scalar:    func() ast.Node { return p.parseScalar() },
		lexer.Union:     func() ast.Node { return p.parseUnion() },
		lexer.Directive: func() ast.Node { return p.parseDirectiveDefinition() },
	}

	for p.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.currToken.Type]; ok {
			node := parseFunc()
			if node != nil {
				nodes[node.GetName()] = node
			}
		}
		p.nextToken()
	}

	return nodes
}

// parseDescription parses a description if present
func (p *Parser) parseDescription() string {
	if p.currToken.Type == lexer.Message {
		description := p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
		return description
	}
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description := p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
		return description
	}
	return ""
}

// parseExtend parse an extend node
//
//	extend type User {
//	  role: Role!
//	}
func (p *Parser) parseExtend() *ast.TypeNode {
	p.parseDescription() // Skip extend description

	p.nextToken()               // Skip 'extend'
	p.expect(lexer.Type, false) // Ensure the next token is 'type', but not move to next token, continue parsing

	// Parse the extended type using parseType
	return p.parseType()
}

func (p *Parser) parseImplements() []string {
	var implements []string
	if p.currToken.Type != lexer.Implements {
		return implements
	}
	p.nextToken()
	for p.currToken.Type != lexer.LeftBrace {
		implements = append(implements, p.currToken.Value)
		p.nextToken()
	}
	return implements
}

func (p *Parser) parseDirectives() []ast.DirectiveNode {
	var directives []ast.DirectiveNode
	if p.currToken.Type != lexer.At {
		return directives
	}
	for p.currToken.Type != lexer.LeftBrace {
		directives = append(directives, p.parseDirective())
	}
	return directives
}

func (p *Parser) parseDirective() ast.DirectiveNode {
	p.nextToken() // skip @
	name := p.currToken.Value
	p.nextToken() // skip name

	var args []ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments()
	}

	return ast.DirectiveNode{
		Name: name,
		Args: args,
	}
}

// parseType parse a type node
//
//	type User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseType() *ast.TypeNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()

	implements := p.parseImplements()
	directives := p.parseDirectives()
	p.expect(lexer.LeftBrace)

	var fields []ast.FieldNode

	for p.currToken.Type != lexer.RightBrace {
		field := p.parseField()
		fields = append(fields, field)
	}
	p.expect(lexer.RightBrace)
	operationType := ast.OperationTypeEntity

	switch name {
	case "Query":
		operationType = ast.OperationTypeQuery
	case "Mutation":
		operationType = ast.OperationTypeMutation
	case "Subscription":
		operationType = ast.OperationTypeSubscription
	}
	node := &ast.TypeNode{
		Name:          name,
		Fields:        fields,
		Description:   description,
		Implements:    implements,
		OperationType: operationType,
		Directives:    directives,
	}
	if p.typeMap == nil {
		p.typeMap = make(map[string]*ast.TypeNode)
	}
	if existingTypeNode, ok := p.typeMap[name]; ok {

		existingTypeNode.Fields = append(existingTypeNode.Fields, node.Fields...)
	} else {
		p.typeMap[name] = node
	}
	return p.typeMap[name]
}

// parseField parse a field node
// "It is ID"
// id: ID!
// name: String!
// age: Int
// email: String
// createdAt: DateTime
func (p *Parser) parseField() ast.FieldNode {
	description := p.parseDescription()
	p.nextToken()
	name := p.currToken.Value

	// parse arguments
	var args []ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments()
	}

	p.expect(lexer.Colon)
	fieldType := p.parseTypeReference()

	return ast.FieldNode{
		Name:        name,
		Type:        fieldType,
		Args:        args,
		Description: description,
	}
}

// parseTypeReference parse a type reference
// ID, String, [Int], [[Int]], [User], [[User]], [User!]
func (p *Parser) parseTypeReference() *ast.FieldType {
	var fieldType *ast.FieldType

	if p.currToken.Type == lexer.LeftBracket {
		p.nextToken() // skip [
		elemType := p.parseTypeReference()
		p.expect(lexer.RightBracket) // skip ]
		fieldType = &ast.FieldType{
			IsList:   true,
			ElemType: elemType,
		}
	} else {

		fieldType = &ast.FieldType{
			Name: p.currToken.Value,
		}
	}
	p.nextToken()
	if p.currToken.Type == lexer.Exclamation {
		fieldType.IsNonNull = true
		p.nextToken()

	}
	return fieldType
}

// parseArguments parse arguments
// (id: ID!, name: String!)
func (p *Parser) parseArguments() []ast.ArgumentNode {
	var args []ast.ArgumentNode
	if p.currToken.Type != lexer.LeftParent {
		return args
	}
	p.expect(lexer.LeftParent)
	for p.currToken.Type != lexer.RightParent {
		args = append(args, p.parseArgument())
		if p.currToken.Type != lexer.RightParent {
			p.expect(lexer.Comma)
		}
	}
	p.expect(lexer.RightParent)
	return args
}

// parseDefaultValue parse a default value
// = 123, = "123", = true, = false, = null
func (p *Parser) parseDefaultValue() string {
	if p.currToken.Type == lexer.Equal {
		p.nextToken() // skip =
		value := p.currToken.Value
		p.nextToken() // skip value
		return value
	}
	return ""
}

// parseArgument parse an argument node
// id: ID!
// name: String!
func (p *Parser) parseArgument() ast.ArgumentNode {
	description := p.parseDescription()
	name := p.currToken.Value

	p.nextToken()                         // skip name
	p.expect(lexer.Colon)                 // skip :
	fieldType := p.parseTypeReference()   // parse type reference
	defaultValue := p.parseDefaultValue() // parse default value

	// parse directives
	directives := p.parseDirectives()

	return ast.ArgumentNode{
		Name:         name,
		Type:         fieldType,
		Description:  description,
		Value:        "",
		Directives:   directives,
		DefaultValue: defaultValue,
	}
}

// parseEnum parse an enum node
//
//	enum Role {
//	  ADMIN
//	  USER
//	  GUEST
//	}
func (p *Parser) parseEnum() *ast.EnumNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()

	p.expect(lexer.LeftBrace)

	var values []string
	for p.currToken.Type != lexer.RightBrace {
		values = append(values, p.currToken.Value)
		p.nextToken()
	}
	p.expect(lexer.RightBrace)
	node := &ast.EnumNode{Name: name, Values: values, Description: description}
	if p.enumMap == nil {
		p.enumMap = make(map[string]*ast.EnumNode)
	}
	p.enumMap[name] = node
	return node
}

// parseInterface parse an interface node
//
//	interface Node {
//	  id: ID!
//	}
func (p *Parser) parseInterface() *ast.InterfaceNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.LeftBrace)

	var fields []ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		fields = append(fields, p.parseField())
	}

	p.expect(lexer.RightBrace)
	node := &ast.InterfaceNode{Name: name, Fields: fields, Description: description}
	if p.interfaceMap == nil {
		p.interfaceMap = make(map[string]*ast.InterfaceNode)
	}
	p.interfaceMap[name] = node
	return node
}

// parseInput parse an input node
//
//	input User {
//	  id: ID!
//	  name: String!
//	  age: Int
//	  email: String
//	  createdAt: DateTime
//	}
func (p *Parser) parseInput() *ast.InputNode {
	// Input types are similar to regular types
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()

	p.expect(lexer.LeftBrace)

	var fields []ast.FieldNode

	for p.currToken.Type != lexer.RightBrace {
		field := p.parseField()
		fields = append(fields, field)
	}
	p.expect(lexer.RightBrace)
	node := &ast.InputNode{
		Name:        name,
		Fields:      fields,
		Description: description,
	}

	if p.inputMap == nil {
		p.inputMap = make(map[string]*ast.InputNode)
	}
	p.inputMap[node.Name] = node
	return node
}

// parseScalar parse a scalar node
// scalar DateTime
func (p *Parser) parseScalar() *ast.ScalarNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	node := &ast.ScalarNode{Name: name, Description: description}
	if p.scalarMap == nil {
		p.scalarMap = make(map[string]*ast.ScalarNode)
	}
	p.scalarMap[name] = node
	return node
}

// parseUnion parse a union node
// union User = Product | Order
func (p *Parser) parseUnion() *ast.UnionNode {
	description := p.parseDescription()

	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.Equal)

	var types []string
	for {
		p.nextToken()
		types = append(types, p.currToken.Value)
		if p.peekToken().Type != lexer.Pipe {
			break
		}
		p.nextToken() // consume the '|'
	}
	node := &ast.UnionNode{Name: name, Types: types, Description: description}
	if p.unionMap == nil {
		p.unionMap = make(map[string]*ast.UnionNode)
	}
	p.unionMap[name] = node
	return node
}

// parseLocations parse locations
// ON FIELD_DEFINITION | ARGUMENT_DEFINITION | INTERFACE | UNION | ENUM | INPUT_OBJECT | SCALAR | OBJECT
func (p *Parser) parseLocations() []ast.DirectiveDefinitionNodeLocation {
	var locations []ast.DirectiveDefinitionNodeLocation

	for {
		locations = append(locations, ast.DirectiveDefinitionNodeLocation(p.currToken.Value))
		p.nextToken()
		if p.currToken.Type != lexer.Pipe {
			break
		}
		p.nextToken()
	}

	return locations
}

// parseDirective parse a directive node
// @skip(if: true)
func (p *Parser) parseDirectiveDefinition() *ast.DirectiveDefinitionNode {
	description := p.parseDescription()

	p.expect(lexer.Directive)
	p.expect(lexer.At)
	name := p.currToken.Value
	p.nextToken() // skip name

	var args []ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments()
	}
	p.expect(lexer.On) // skip ON
	locations := p.parseLocations()

	node := &ast.DirectiveDefinitionNode{Name: name, Args: args, Description: description, Locations: locations}
	if p.directiveMap == nil {
		p.directiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	p.directiveMap[name] = node
	return node
}

// expect check if the current token is the expected token
// if not, panic
func (p *Parser) expect(t lexer.TokenType, options ...bool) {
	if p.currToken.Type != t {
		panic(fmt.Sprintf("expect token: %s, but got: %s at line %d position %d", t, p.currToken.Value, p.currToken.Line, p.currToken.LinePosition))
	}

	if len(options) == 0 || options[0] {
		p.nextToken()
	}
}

// peekToken return the next token
func (p *Parser) peekToken() *lexer.Token {
	return p.lexer.PeekToken()
}

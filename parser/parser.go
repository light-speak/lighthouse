package parser

import (
	"os"

	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
)

type Parser struct {
	lexer     *lexer.Lexer
	currToken *lexer.Token
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
	return &Parser{lexer: lexer, currToken: lexer.NextToken()}
}

// nextToken move to next token
func (p *Parser) nextToken() {
	p.currToken = p.lexer.NextToken()
}

// ParseSchema parse schema
// return a list of ast nodes
// the nodes is a list of type, enum, interface, input, scalar, union, directive
func (p *Parser) ParseSchema() []ast.ASTNode {
	var nodes []ast.ASTNode
	for p.currToken.Type != lexer.EOF {
		switch p.currToken.Type {
		case lexer.Type:
			nodes = append(nodes, p.parseType())
		case lexer.Enum:
			nodes = append(nodes, p.parseEnum())
		case lexer.Interface:
			nodes = append(nodes, p.parseInterface())
		case lexer.Input:
			nodes = append(nodes, p.parseInput())
		case lexer.Scalar:
			nodes = append(nodes, p.parseScalar())
		case lexer.Union:
			nodes = append(nodes, p.parseUnion())
		case lexer.Directive:
			nodes = append(nodes, p.parseDirective())
		default:
			p.nextToken()
		}
	}
	return nodes
}

// parseType parse a type node
// type User {
//   id: ID!
//   name: String!
//   age: Int
//   email: String
//   createdAt: DateTime
// }
func (p *Parser) parseType() *ast.TypeNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.LeftBrace)

	var fields []ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		fields = append(fields, p.parseField())
		p.nextToken()
	}
	p.expect(lexer.RightBrace)
	return &ast.TypeNode{Name: name, Fields: fields, Description: description}
}


// parseField parse a field node
// id: ID!
// name: String!
// age: Int
// email: String
// createdAt: DateTime
func (p *Parser) parseField() ast.FieldNode {
	var description string
	if p.currToken.Type == lexer.Message {
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.Colon)
	p.nextToken()
	fieldType := p.currToken.Value

	var args []ast.ArgumentNode
	if p.peekToken().Type == lexer.LeftParent {
		p.nextToken()
		args = p.parseArguments()
	}

	return ast.FieldNode{Name: name, Type: ast.FieldType(fieldType), Args: args, Description: description}
}

// parseArguments parse arguments
// (id: ID!, name: String!)
func (p *Parser) parseArguments() []ast.ArgumentNode {
	var args []ast.ArgumentNode
	p.expect(lexer.LeftParent)
	for p.currToken.Type != lexer.RightParent {
		args = append(args, p.parseArgument())
		if p.peekToken().Type != lexer.RightParent {
			p.expect(lexer.Comma)
		}
	}
	p.expect(lexer.RightParent)
	return args
}

// parseArgument parse an argument node
// id: ID!
// name: String!
func (p *Parser) parseArgument() ast.ArgumentNode {
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.Colon)
	p.nextToken()
	argType := p.currToken.Value
	p.nextToken()
	return ast.ArgumentNode{Name: name, Type: ast.FieldType(argType)}
}

// parseEnum parse an enum node
// enum Role {
//   ADMIN
//   USER
//   GUEST
// }
func (p *Parser) parseEnum() *ast.EnumNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
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
	return &ast.EnumNode{Name: name, Values: values, Description: description}
}

// parseInterface parse an interface node
// interface Node {
//   id: ID!
// }
func (p *Parser) parseInterface() *ast.InterfaceNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	p.expect(lexer.LeftBrace)

	var fields []ast.FieldNode
	for p.currToken.Type != lexer.RightBrace {
		fields = append(fields, p.parseField())
		p.nextToken()
	}
	p.expect(lexer.RightBrace)
	return &ast.InterfaceNode{Name: name, Fields: fields, Description: description}
}

// parseInput parse an input node
// input User {
//   id: ID!
//   name: String!
//   age: Int
//   email: String
//   createdAt: DateTime
// }
func (p *Parser) parseInput() *ast.TypeNode {
	// Input types are similar to regular types
	return p.parseType()
}

// parseScalar parse a scalar node
// scalar DateTime
func (p *Parser) parseScalar() *ast.ScalarNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
	p.nextToken()
	name := p.currToken.Value
	p.nextToken()
	return &ast.ScalarNode{Name: name, Description: description}
}

// parseUnion parse a union node
// union User = Product | Order
func (p *Parser) parseUnion() *ast.UnionNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
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
	return &ast.UnionNode{Name: name, Types: types, Description: description}
}

// parseDirective parse a directive node
// @skip(if: true)
func (p *Parser) parseDirective() *ast.DirectiveNode {
	var description string
	if p.peekToken().Type == lexer.Message {
		p.nextToken()
		description = p.currToken.Value[1 : len(p.currToken.Value)-1] // Remove quotes
		p.nextToken()
	}
	
	p.nextToken()
	name := p.currToken.Value
	p.nextToken()

	var args []ast.ArgumentNode
	if p.currToken.Type == lexer.LeftParent {
		args = p.parseArguments()
	}

	return &ast.DirectiveNode{Name: name, Args: args, Description: description}
}

// expect check if the current token is the expected token
// if not, panic
func (p *Parser) expect(t lexer.TokenType) {
	if p.currToken.Type != t {
		panic("unexpected token: " + p.currToken.Value)
	}
	p.nextToken()
}

// peekToken return the next token
func (p *Parser) peekToken() *lexer.Token {
	return p.lexer.PeekToken()
}
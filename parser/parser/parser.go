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

	// TypeMap, enumMap, scalarMap, unionMap, inputMap, interfaceMap, and directiveMap are all maps
	// that store parsed AST nodes. The keys are the names of the respective types, enums, scalars, unions,
	// input types, interfaces, and directives, while the values are pointers to their corresponding AST node structures.
	TypeMap      map[string]*ast.TypeNode
	EnumMap      map[string]*ast.EnumNode
	ScalarMap    map[string]*ast.ScalarNode
	UnionMap     map[string]*ast.UnionNode
	InputMap     map[string]*ast.InputNode
	InterfaceMap map[string]*ast.InterfaceNode
	DirectiveMap map[string]*ast.DirectiveDefinitionNode
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
	p.AddReserved()
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
		lexer.Fragment:  func() ast.Node { return p.parseFragment() },
	}

	for p.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.currToken.Type]; ok {
			node := parseFunc()
			if node != nil {
				nodes[node.GetName()] = node
			}
		}
		if p.currToken.Type != lexer.Directive && p.currToken.Type != lexer.Union {
			p.nextToken()
		}
	}

	return nodes
}

// PreviousToken return Previous Token
func (p *Parser) PreviousToken() *lexer.Token {
	return p.lexer.PreviousToken()
}

// expect check if the current token is the expected token
// if not, panic
func (p *Parser) expect(t lexer.TokenType, options ...bool) {
	if p.currToken.Type != t {
		panic(fmt.Sprintf("expect: %s, but got: %s at line %d position %d", t, p.currToken.Value, p.currToken.Line, p.currToken.LinePosition))
	}

	if len(options) == 0 || options[0] {
		p.nextToken()
	}
}

func (p *Parser) AddScalar(node *ast.ScalarNode) {
	if p.ScalarMap == nil {
		p.ScalarMap = make(map[string]*ast.ScalarNode)
	}
	p.ScalarMap[node.Name] = node
}

func (p *Parser) AddDirective(node *ast.DirectiveDefinitionNode) {
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	p.DirectiveMap[node.Name] = node
}

func (p *Parser) AddInput(node *ast.InputNode) {
	if p.InputMap == nil {
		p.InputMap = make(map[string]*ast.InputNode)
	}
	p.InputMap[node.Name] = node
}

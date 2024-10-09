package parser

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
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

	ScalarTypeMap map[string]ast.ScalarType
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

	p.MergeScalarType()

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
func (p *Parser) AddScalar(node *ast.ScalarNode) ast.Node {
	if p.ScalarMap == nil {
		p.ScalarMap = make(map[string]*ast.ScalarNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Scalar '%s' already defined", node.Name))
	}
	p.ScalarMap[node.Name] = node
	return p.ScalarMap[node.Name]
}

func (p *Parser) AddInput(node *ast.InputNode) ast.Node {
	if p.InputMap == nil {
		p.InputMap = make(map[string]*ast.InputNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Input '%s' already defined", node.Name))
	}
	p.InputMap[node.Name] = node
	return p.InputMap[node.Name]
}

func (p *Parser) AddInterface(node *ast.InterfaceNode) ast.Node {
	if p.InterfaceMap == nil {
		p.InterfaceMap = make(map[string]*ast.InterfaceNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Interface '%s' already defined", node.Name))
	}
	p.InterfaceMap[node.Name] = node
	return p.InterfaceMap[node.Name]
}

func (p *Parser) AddDirectiveDefinition(node *ast.DirectiveDefinitionNode) ast.Node {
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	if _, exists := p.DirectiveMap[node.Name]; exists {
		panic(fmt.Sprintf("Duplicate directive definition: '%s'", node.Name))
	}
	p.DirectiveMap[node.Name] = node
	return p.DirectiveMap[node.Name]
}

func (p *Parser) AddEnum(node *ast.EnumNode) ast.Node {
	if p.EnumMap == nil {
		p.EnumMap = make(map[string]*ast.EnumNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Enum '%s' already defined", node.Name))
	}
	p.EnumMap[node.Name] = node
	return p.EnumMap[node.Name]
}

func (p *Parser) AddUnion(node *ast.UnionNode) ast.Node {
	if p.UnionMap == nil {
		p.UnionMap = make(map[string]*ast.UnionNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Union '%s' already defined", node.Name))
	}
	p.UnionMap[node.Name] = node
	return p.UnionMap[node.Name]
}

func (p *Parser) AddType(name string, node *ast.TypeNode, extends bool) ast.Node {
	if p.TypeMap == nil {
		p.TypeMap = make(map[string]*ast.TypeNode)
	}
	if existingTypeNode, ok := p.TypeMap[name]; ok {
		if extends {
			for _, field := range node.Fields {
				if existingField := existingTypeNode.GetField(field.Name); existingField != nil {
					panic(fmt.Sprintf("Duplicate field: '%s' in type: '%s'", field.Name, name))
				}
			}
			existingTypeNode.Fields = append(existingTypeNode.Fields, node.Fields...)
		} else {
			panic(fmt.Sprintf("Duplicate type definition: '%s'", name))
		}
	} else {
		if p.isNameConflict(name) {
			panic(fmt.Sprintf("Name conflict: Type '%s' already defined", name))
		}
		p.TypeMap[name] = node
	}
	return p.TypeMap[name]
}

func (p *Parser) AddScalarType(name string, scalarType ast.ScalarType) {
	if p.ScalarTypeMap == nil {
		p.ScalarTypeMap = make(map[string]ast.ScalarType)
	}
	if _, exists := p.ScalarTypeMap[name]; exists {
		panic(fmt.Sprintf("Duplicate ScalarType definition: '%s'", name))
	}
	p.ScalarTypeMap[name] = scalarType
}

func (p *Parser) AddDirective(node *ast.DirectiveDefinitionNode) ast.Node {
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	if _, exists := p.DirectiveMap[node.Name]; exists {
		panic(fmt.Sprintf("Duplicate directive definition: '%s'", node.Name))
	}
	p.DirectiveMap[node.Name] = node
	return p.DirectiveMap[node.Name]
}

func (p *Parser) isNameConflict(name string) bool {
	return (p.ScalarMap != nil && p.ScalarMap[name] != nil) ||
		(p.TypeMap != nil && p.TypeMap[name] != nil) ||
		(p.InputMap != nil && p.InputMap[name] != nil) ||
		(p.InterfaceMap != nil && p.InterfaceMap[name] != nil) ||
		(p.EnumMap != nil && p.EnumMap[name] != nil) ||
		(p.UnionMap != nil && p.UnionMap[name] != nil)
}

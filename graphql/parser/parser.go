package parser

import (
	"fmt"
	"os"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/log"
)

// Parser is responsible for parsing the GraphQL schema.
// It contains a lexer for tokenizing the input and the current token being processed.
// The various maps are used to store different types of AST nodes for quick lookup and management during parsing.
type Parser struct {
	// lexer is the Lexer instance used for lexical analysis, converting the input GraphQL text into a stream of tokens.
	lexer *lexer.Lexer

	// currToken is the current token being processed, which helps the parser determine its state.
	currToken *lexer.Token

	// Nodes is a map of all parsed nodes
	Nodes map[string]ast.Node

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
	FragmentMap  map[string]*ast.FragmentNode

	ScalarTypeMap map[string]ast.ScalarType
}

// ReadGraphQLFile read graphql file and return a lexer
func ReadGraphQLFile(path string) (*lexer.Lexer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return lexer.NewLexer([]*lexer.Content{{Path: &path, Content: string(content)}}), nil
}

func ReadGraphQLFiles(paths []string) (*lexer.Lexer, error) {
	contents := make([]*lexer.Content, 0)
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		contents = append(contents, &lexer.Content{Path: &path, Content: string(content)})
	}
	return lexer.NewLexer(contents), nil
}

// NewParser create a new parser
func NewParser(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken() // Initialize currToken
	return p
}

// nextToken move to next token
func (p *Parser) nextToken() error {
	var err error
	p.currToken, err = p.lexer.NextToken()
	if err != nil {
		return err
	}
	return nil
}

// ParseSchema parse schema
// return a list of ast nodes
// the nodes is a list of type, enum, interface, input, scalar, union, directive, extend
func (p *Parser) ParseSchema() map[string]ast.Node {
	p.Nodes = make(map[string]ast.Node)
	log.Debug().Msgf("currToken: %+v", p.Nodes)
	tokenTypeToParseFunc := map[lexer.TokenType]func(){
		lexer.Type:      func() { p.parseType() },
		lexer.Extend:    func() { p.parseExtend() },
		lexer.Enum:      func() { p.parseEnum() },
		lexer.Interface: func() { p.parseInterface() },
		lexer.Input:     func() { p.parseInput() },
		lexer.Scalar:    func() { p.parseScalar() },
		lexer.Union:     func() { p.parseUnion() },
		lexer.Directive: func() { p.parseDirectiveDefinition() },
		lexer.Fragment:  func() { p.parseFragment() },
	}

	for p.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.currToken.Type]; ok {
			parseFunc()
		}
		if p.currToken.Type != lexer.Directive && p.currToken.Type != lexer.Union {
			p.nextToken()
		}
	}
	p.AddReserved()
	p.MergeScalarType()
	return p.Nodes
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
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Scalar '%s' already defined", node.Name))
	}
	p.ScalarMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) AddInput(node *ast.InputNode) {
	if p.InputMap == nil {
		p.InputMap = make(map[string]*ast.InputNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Input '%s' already defined", node.Name))
	}
	p.InputMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) AddInterface(node *ast.InterfaceNode) {
	if p.InterfaceMap == nil {
		p.InterfaceMap = make(map[string]*ast.InterfaceNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Interface '%s' already defined", node.Name))
	}
	p.InterfaceMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) AddDirectiveDefinition(node *ast.DirectiveDefinitionNode) {
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Directive '%s' already defined", node.Name))
	}
	p.DirectiveMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) AddEnum(node *ast.EnumNode) {
	if p.EnumMap == nil {
		p.EnumMap = make(map[string]*ast.EnumNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Enum '%s' already defined", node.Name))
	}
	p.EnumMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) AddUnion(node *ast.UnionNode) {
	if p.UnionMap == nil {
		p.UnionMap = make(map[string]*ast.UnionNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Union '%s' already defined", node.Name))
	}
	p.UnionMap[node.Name] = node
	p.Nodes[node.Name] = node
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
	p.Nodes[name] = p.TypeMap[name]
	return p.TypeMap[name]
}

func (p *Parser) AddScalarType(name string, scalarType ast.ScalarType) {
	if p.ScalarTypeMap == nil {
		p.ScalarTypeMap = make(map[string]ast.ScalarType)
	}
	if p.isNameConflict(name) {
		panic(fmt.Sprintf("Name conflict: Scalar type '%s' already defined", name))
	}
	p.ScalarTypeMap[name] = scalarType
}

func (p *Parser) AddDirective(node *ast.DirectiveDefinitionNode) {
	if p.DirectiveMap == nil {
		p.DirectiveMap = make(map[string]*ast.DirectiveDefinitionNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Directive '%s' already defined", node.Name))
	}
	p.DirectiveMap[node.Name] = node
	p.Nodes[node.Name] = node
}

func (p *Parser) isNameConflict(name string) bool {
	return (p.ScalarMap != nil && p.ScalarMap[name] != nil) ||
		(p.TypeMap != nil && p.TypeMap[name] != nil) ||
		(p.InputMap != nil && p.InputMap[name] != nil) ||
		(p.InterfaceMap != nil && p.InterfaceMap[name] != nil) ||
		(p.EnumMap != nil && p.EnumMap[name] != nil) ||
		(p.UnionMap != nil && p.UnionMap[name] != nil) ||
		(p.DirectiveMap != nil && p.DirectiveMap[name] != nil) ||
		(p.FragmentMap != nil && p.FragmentMap[name] != nil)
}

func (p *Parser) AddFragment(node *ast.FragmentNode) {
	if p.FragmentMap == nil {
		p.FragmentMap = make(map[string]*ast.FragmentNode)
	}
	if p.isNameConflict(node.Name) {
		panic(fmt.Sprintf("Name conflict: Fragment '%s' already defined", node.Name))
	}
	p.FragmentMap[node.Name] = node
	p.Nodes[node.Name] = node
}

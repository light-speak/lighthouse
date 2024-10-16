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
	QueryParser *QueryParser

	// lexer is the Lexer instance used for lexical analysis, converting the input GraphQL text into a stream of tokens.
	lexer *lexer.Lexer

	// currToken is the current token being processed, which helps the parser determine its state.
	currToken *lexer.Token

	// NodeStore is a store of all parsed nodes
	NodeStore *ast.NodeStore
}

type QueryParser struct {
	Parser  *Parser
	QueryId string
	// OperationNode *ast.OperationNode
	Fragments map[string]*ast.FragmentNode
}

func (p *QueryParser) AddFragment(node *ast.FragmentNode) {
	if p.Fragments == nil {
		p.Fragments = make(map[string]*ast.FragmentNode)
	}
	p.Fragments[node.Name] = node
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

// NewQueryParser create a new query parser, which is used to parse query
// it will parse the operation and fragments
// and store them in the QueryParser
func (p *Parser) NewQueryParser(queryLexer *lexer.Lexer) *QueryParser {
	p.lexer = queryLexer
	p.QueryParser = &QueryParser{Parser: p}
	p.nextToken()
	return p.QueryParser
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
	p.NodeStore = &ast.NodeStore{}
	p.NodeStore.InitStore()
	tokenTypeToParseFunc := map[lexer.TokenType]func(){
		// lexer.LowerQuery:        func() { p.parseOperation() },
		// lexer.LowerMutation:     func() { p.parseOperation() },
		// lexer.LowerSubscription: func() { p.parseOperation() },
		lexer.Type:      func() { p.parseObject() },
		lexer.Extend:    func() { p.parseExtend() },
		lexer.Enum:      func() { p.parseEnum() },
		lexer.Interface: func() { p.parseInterface() },
		lexer.Input:     func() { p.parseInput() },
		lexer.Scalar:    func() { p.parseScalar() },
		lexer.Union:     func() { p.parseUnion() },
		lexer.Directive: func() { p.parseDirectiveDefinition() },
	}

	for p.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.currToken.Type]; ok {
			parseFunc()
		}
		if p.currToken.Type != lexer.Directive && p.currToken.Type != lexer.Union {
			p.nextToken()
		}
	}

	if p.QueryParser == nil {
		p.AddReserved()
		p.MergeScalarType()
	}
	return p.NodeStore.Nodes
}

// PreviousToken return Previous Token
func (p *Parser) PreviousToken() *lexer.Token {
	return p.lexer.PreviousToken()
}

// expect check if the current token is the expected token
// if not, panic
func (p *Parser) expect(t lexer.TokenType, options ...bool) {
	if p.currToken.Type != t {
		panic(fmt.Sprintf("expect: %s but got: %s at line %d position %d", t, p.currToken.Value, p.currToken.Line, p.currToken.LinePosition))
	}

	if len(options) == 0 || options[0] {
		p.nextToken()
	}
}

func (p *Parser) AddScalar(node *ast.ScalarNode) {
	p.isNameConflict(node.Name)
	p.NodeStore.Scalars[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddInput(node *ast.InputObjectNode) {
	p.isNameConflict(node.Name)
	p.NodeStore.Inputs[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddInterface(node *ast.InterfaceNode) {
	p.isNameConflict(node.Name)
	p.NodeStore.Interfaces[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddDirectiveDefinition(node *ast.DirectiveDefinition) {
	p.isNameConflict(node.Name)
	p.NodeStore.Directives[node.Name] = node
	p.NodeStore.Names[node.Name] = node
}

func (p *Parser) AddEnum(node *ast.EnumNode) {
	p.isNameConflict(node.Name)
	p.NodeStore.Enums[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddUnion(node *ast.UnionNode) {
	p.isNameConflict(node.Name)
	p.NodeStore.Unions[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddObject(node *ast.ObjectNode, extend bool) ast.Node {
	if extend {
		if _, ok := p.NodeStore.Objects[node.Name]; !ok {
			p.NodeStore.Objects[node.Name] = node
		} else {
			for _, field := range node.Fields {
				if _, ok := p.NodeStore.Objects[node.Name].Fields[field.Name]; ok {
					panic(fmt.Sprintf("Name conflict: Field '%s' already defined", field.Name))
				}
				p.NodeStore.Objects[node.Name].Fields[field.Name] = field
			}
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Objects[node.Name] = node
	}
	p.NodeStore.Nodes[node.Name] = p.NodeStore.Objects[node.Name]
	p.NodeStore.Names[node.Name] = node
	return node
}

func (p *Parser) isNameConflict(name string) {
	if p.NodeStore.Names[name] != nil {
		panic(fmt.Sprintf("Name conflict: '%s' already defined", name))
	}
}

// AddScalarType add a scalar type
func (p *Parser) AddScalarType(name string, scalarType ast.ScalarType) {
	if p.NodeStore.ScalarTypes == nil {
		p.NodeStore.ScalarTypes = make(map[string]ast.ScalarType)
	}
	p.NodeStore.ScalarTypes[name] = scalarType
}

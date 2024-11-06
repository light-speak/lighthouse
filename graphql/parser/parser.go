package parser

import (
	_ "embed"
	"fmt"
	"os"
	"sync"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
)

//go:embed base.graphql
var baseSchema string

var parserPool = sync.Pool{
	New: func() interface{} {
		return &Parser{
			NodeStore: &ast.NodeStore{},
		}
	},
}

type ParserInterface interface {
	GetLexer() *lexer.Lexer
	NextToken() error
}

type BaseParser struct {
	lexer     *lexer.Lexer
	currToken *lexer.Token
}

func (p *BaseParser) GetLexer() *lexer.Lexer {
	return p.lexer
}

func (p *BaseParser) NextToken() error {
	var err error
	p.currToken, err = p.lexer.NextToken()
	if err != nil {
		return err
	}
	return nil
}

type Parser struct {
	BaseParser
	NodeStore *ast.NodeStore
}

func ReadGraphQLFile(path string) (*lexer.Lexer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return lexer.NewLexer([]*lexer.Content{{Path: &path, Content: string(content) + "\n"}}), nil
}

func ReadGraphQLFiles(paths []string) (*lexer.Lexer, errors.GraphqlErrorInterface) {
	contents := make([]*lexer.Content, 0, len(paths)+1)
	contents = append(contents, &lexer.Content{Path: nil, Content: baseSchema + "\n"})

	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, &errors.ParserError{
				Message:   err.Error(),
				Locations: &errors.GraphqlLocation{Line: 1, Column: 1},
			}
		}
		contents = append(contents, &lexer.Content{Path: &path, Content: string(content) + "\n"})
	}
	return lexer.NewLexer(contents), nil
}

func NewParser(lexer *lexer.Lexer) *Parser {
	p := parserPool.Get().(*Parser)
	p.lexer = lexer
	p.nextToken()
	return p
}

func (p *Parser) NewQueryParser(queryLexer *lexer.Lexer) *QueryParser {
	qp := &QueryParser{BaseParser: BaseParser{lexer: queryLexer}, Parser: p}
	qp.NextToken()
	return qp
}

func (p *Parser) nextToken() error {
	var err error
	p.currToken, err = p.lexer.NextToken()
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) ParseSchema() map[string]ast.Node {
	p.NodeStore = &ast.NodeStore{}
	p.NodeStore.InitStore()

	tokenTypeToParseFunc := map[lexer.TokenType]func(){
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
		if p.currToken.Type != lexer.Directive &&
			p.currToken.Type != lexer.Union &&
			p.currToken.Type != lexer.Type &&
			p.currToken.Type != lexer.Extend &&
			p.currToken.Type != lexer.Input &&
			p.currToken.Type != lexer.Enum &&
			p.currToken.Type != lexer.Interface {
			p.nextToken()
		}
	}

	p.AddReserved()

	defer parserPool.Put(p)

	return p.NodeStore.Nodes
}

func (p *BaseParser) PreviousToken() *lexer.Token {
	return p.lexer.PreviousToken()
}

func (p *BaseParser) expect(t lexer.TokenType, options ...bool) {
	if p.currToken.Type != t {
		panic(fmt.Sprintf("expect: %s but got: %s at line %d position %d", t, p.currToken.Value, p.currToken.Line, p.currToken.LinePosition))
	}

	if len(options) == 0 || options[0] {
		p.NextToken()
	}
}

func (p *Parser) AddInput(node *ast.InputObjectNode, extend bool) {
	if extend {
		if existingNode, ok := p.NodeStore.Inputs[node.Name]; ok {
			if existingNode.Fields == nil {
				existingNode.Fields = make(map[string]*ast.Field)
			}
			for name, field := range node.Fields {
				if _, ok := existingNode.Fields[name]; ok {
					panic(fmt.Sprintf("Name conflict: Field '%s' already defined in Input '%s'", name, node.Name))
				}
				existingNode.Fields[name] = field
			}
			existingNode.Directives = append(existingNode.Directives, node.Directives...)
			node = existingNode
		} else {
			p.isNameConflict(node.Name)
			p.NodeStore.Inputs[node.Name] = node
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Inputs[node.Name] = node
	}
	node.Kind = ast.KindInputObject
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddInterface(node *ast.InterfaceNode, extend bool) {
	if extend {
		if existingNode, ok := p.NodeStore.Interfaces[node.Name]; ok {
			if existingNode.Fields == nil {
				existingNode.Fields = make(map[string]*ast.Field)
			}
			for name, field := range node.Fields {
				if _, ok := existingNode.Fields[name]; ok {
					panic(fmt.Sprintf("Name conflict: Field '%s' already defined in Interface '%s'", name, node.Name))
				}
				existingNode.Fields[name] = field
			}
			existingNode.Directives = append(existingNode.Directives, node.Directives...)
			node = existingNode
		} else {
			p.isNameConflict(node.Name)
			p.NodeStore.Interfaces[node.Name] = node
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Interfaces[node.Name] = node
	}
	node.Kind = ast.KindInterface
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddEnum(node *ast.EnumNode, extend bool) {
	if extend {
		if existingNode, ok := p.NodeStore.Enums[node.Name]; ok {
			if existingNode.EnumValues == nil {
				existingNode.EnumValues = make(map[string]*ast.EnumValue)
			}
			for name, value := range node.EnumValues {
				if _, ok := existingNode.EnumValues[name]; ok {
					panic(fmt.Sprintf("Name conflict: EnumValue '%s' already defined in Enum '%s'", name, node.Name))
				}
				existingNode.EnumValues[name] = value
			}
			existingNode.Directives = append(existingNode.Directives, node.Directives...)
			node = existingNode
		} else {
			p.isNameConflict(node.Name)
			p.NodeStore.Enums[node.Name] = node
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Enums[node.Name] = node
	}
	node.Kind = ast.KindEnum
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddUnion(node *ast.UnionNode, extend bool) {
	if extend {
		if existingNode, ok := p.NodeStore.Unions[node.Name]; ok {
			if existingNode.TypeNames == nil {
				existingNode.TypeNames = make(map[string]string)
			}
			for name, typeName := range node.TypeNames {
				if _, ok := existingNode.TypeNames[name]; ok {
					panic(fmt.Sprintf("Name conflict: TypeName '%s' already defined in Union '%s'", name, node.Name))
				}
				existingNode.TypeNames[name] = typeName
			}
			existingNode.Directives = append(existingNode.Directives, node.Directives...)
			node = existingNode
		} else {
			p.isNameConflict(node.Name)
			p.NodeStore.Unions[node.Name] = node
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Unions[node.Name] = node
	}
	node.Kind = ast.KindUnion
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) AddObject(node *ast.ObjectNode, extend bool) ast.Node {
	if extend {
		if existingNode, ok := p.NodeStore.Objects[node.Name]; ok {
			if existingNode.Fields == nil {
				existingNode.Fields = make(map[string]*ast.Field)
			}
			for name, field := range node.Fields {
				if _, ok := existingNode.Fields[name]; ok {
					panic(fmt.Sprintf("Name conflict: Field '%s' already defined in Object '%s'", name, node.Name))
				}
				existingNode.Fields[name] = field
			}
			existingNode.Directives = append(existingNode.Directives, node.Directives...)
			existingNode.InterfaceNames = append(existingNode.InterfaceNames, node.InterfaceNames...)
			node = existingNode
		} else {
			p.isNameConflict(node.Name)
			p.NodeStore.Objects[node.Name] = node
		}
	} else {
		p.isNameConflict(node.Name)
		p.NodeStore.Objects[node.Name] = node
	}
	node.Kind = ast.KindObject
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
	return node
}

func (p *Parser) AddDirectiveDefinition(node *ast.DirectiveDefinition) {
	p.isNameConflict(node.Name)
	p.NodeStore.Directives[node.Name] = node
	p.NodeStore.Names[node.Name] = node
}

func (p *Parser) AddScalar(node *ast.ScalarNode) {
	p.isNameConflict(node.Name)
	node.Kind = ast.KindScalar
	p.NodeStore.Scalars[node.Name] = node
	p.NodeStore.Names[node.Name] = node
	p.NodeStore.Nodes[node.Name] = node
}

func (p *Parser) isNameConflict(name string) {
	if p.NodeStore.Names[name] != nil {
		panic(fmt.Sprintf("Name conflict: '%s' already defined", name))
	}
}

func (p *Parser) AddScalarType(name string, scalarType ast.ScalarType) {
	if p.NodeStore.ScalarTypes == nil {
		p.NodeStore.ScalarTypes = make(map[string]ast.ScalarType)
	}
	p.NodeStore.ScalarTypes[name] = scalarType
}

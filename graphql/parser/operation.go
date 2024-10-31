package parser

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/utils"
)

type QueryParser struct {
	Parser  *Parser
	QueryId string

	Name          string
	Fragments     map[string]*ast.Fragment
	Variables     map[string]any
	Args          map[string]*ast.Argument
	Fields        map[string]*ast.Field
	Directives    []*ast.Directive
	Location      ast.Location
	OperationType string
}

func (p *QueryParser) Validate(store *ast.NodeStore) errors.GraphqlErrorInterface {
	for _, arg := range p.Args {
		if p.Variables[arg.Name] == nil {
			return &errors.GraphQLError{
				Message:   fmt.Sprintf("variable %s not found", arg.Name),
				Locations: []*errors.GraphqlLocation{arg.GetLocation()},
			}
		}
		arg.Value = p.Variables[arg.Name]
		arg.Validate(store, p.Args, nil)
	}
	for _, field := range p.Fields {
		obj := p.Parser.NodeStore.Objects[p.OperationType]
		if obj == nil {
			return &errors.ParserError{
				Message:   "Operation " + p.OperationType + " not found",
				Locations: field.GetLocation(),
			}
		}
		if err := field.Validate(store, obj.Fields, obj, ast.LocationField, p.Fragments, p.Args); err != nil {
			return err
		}
	}
	return nil
}

func (p *QueryParser) ParseSchema() *QueryParser {
	defer func() {
		e := recover()
		if e != nil {
			panic(&errors.ParserError{
				Message:   e.(string),
				Locations: &errors.GraphqlLocation{Line: 1, Column: 1},
			})
		}
	}()
	tokenTypeToParseFunc := map[lexer.TokenType]func(){
		lexer.LowerQuery:        func() { p.parseOperation() },
		lexer.LowerMutation:     func() { p.parseOperation() },
		lexer.LowerSubscription: func() { p.parseOperation() },
		lexer.Fragment:          func() { p.parseFragment() },
	}

	for p.Parser.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.Parser.currToken.Type]; ok {
			parseFunc()
		}
		p.Parser.nextToken()
	}

	return p
}

func (p *QueryParser) AddFragment(node *ast.Fragment) {
	if p.Fragments == nil {
		p.Fragments = make(map[string]*ast.Fragment)
	}
	p.Fragments[node.Name] = node
}

func (p *QueryParser) parseFragment() {
	node := &ast.Fragment{
		Name: p.Parser.expectAndGetValue(lexer.Fragment),
		On:   p.Parser.expectAndGetValue(lexer.On),
	}

	node.Directives = p.Parser.parseDirectives()

	p.Parser.expect(lexer.LeftBrace)
	node.Fields = make(map[string]*ast.Field)
	for {
		field := p.Parser.parseField(false, "")
		if _, ok := node.Fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		node.Fields[field.Name] = field
		if p.Parser.currToken.Type == lexer.RightBrace {
			break
		}
	}

	p.AddFragment(node)
}

func (p *QueryParser) parseOperation() error {
	operationType := utils.UcFirst(string(p.Parser.currToken.Type))
	location := ast.LocationQuery
	if operationType == "Query" {
		location = ast.LocationQuery
	} else if operationType == "Mutation" {
		location = ast.LocationMutation
	} else if operationType == "Subscription" {
		location = ast.LocationSubscription
	}
	p.Location = location
	p.OperationType = operationType

	p.Parser.nextToken()
	if p.Parser.currToken.Type == lexer.Letter {
		p.Name = p.Parser.currToken.Value
		p.Parser.nextToken()
	}

	if p.Parser.currToken.Type == lexer.LeftParent {
		p.Args = p.Parser.parseArguments()
	}
	p.Directives = p.Parser.parseDirectives()

	p.Parser.expect(lexer.LeftBrace)
	fields := make(map[string]*ast.Field)
	for p.Parser.currToken.Type != lexer.RightBrace {
		field := p.Parser.parseField(true, "")
		if field == nil {
			continue
		}
		fields[field.Name] = field
	}
	p.Fields = fields

	return nil
}

package parser

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
)

// import (
// 	"github.com/light-speak/lighthouse/graphql/ast"
// 	"github.com/light-speak/lighthouse/graphql/parser/lexer"
// )

// func (p *Parser) parseOperation() {
// 	node := &ast.OperationNode{}
// 	operationTypes := map[lexer.TokenType]ast.OperationType{
// 		lexer.LowerQuery:        ast.QueryOperation,
// 		lexer.LowerMutation:     ast.MutationOperation,
// 		lexer.LowerSubscription: ast.SubscriptionOperation,
// 	}

// 	if opType, ok := operationTypes[p.currToken.Type]; ok {
// 		node.Type = opType
// 		node.Name = p.expectAndGetValue(p.currToken.Type)
// 	} else {
// 		panic("invalid operation type: " + p.currToken.Value)
// 	}

// 	if p.currToken.Type == lexer.LeftParent {
// 		node.Args = p.parseArguments(node)
// 	}

// 	p.expect(lexer.LeftBrace)
// 	var fields []*ast.FieldNode
// 	for p.currToken.Type != lexer.RightBrace {
// 		field := p.parseField(node)
// 		fields = append(fields, field)
// 	}
// 	node.Fields = fields
// 	p.QueryParser.OperationNode = node

// }

type QueryParser struct {
	Parser  *Parser
	QueryId string
	// OperationNode *ast.OperationNode
	Fragments map[string]*ast.FragmentNode
}

func (p *QueryParser) ParseSchema() *QueryParser {
	tokenTypeToParseFunc := map[lexer.TokenType]func(){
		// lexer.LowerQuery:        func() { p.parseOperation() },
		// lexer.LowerMutation:     func() { p.parseOperation() },
		// lexer.LowerSubscription: func() { p.parseOperation() },
		lexer.Fragment: func() { p.parseFragment() },
	}

	for p.Parser.currToken.Type != lexer.EOF {
		if parseFunc, ok := tokenTypeToParseFunc[p.Parser.currToken.Type]; ok {
			parseFunc()
		}
		p.Parser.nextToken()
	}

	return p
}

func (p *QueryParser) AddFragment(node *ast.FragmentNode) {
	if p.Fragments == nil {
		p.Fragments = make(map[string]*ast.FragmentNode)
	}
	p.Fragments[node.Name] = node
}

func (p *QueryParser) parseFragment() {
	node := &ast.FragmentNode{
		Name: p.Parser.expectAndGetValue(lexer.Fragment),
		On:   p.Parser.expectAndGetValue(lexer.On),
	}

	node.Directives = p.Parser.parseDirectives()

	p.Parser.expect(lexer.LeftBrace)
	node.Fields = make(map[string]*ast.Field)
	for p.Parser.currToken.Type == lexer.Letter || p.Parser.currToken.Type == lexer.TripleDot {
		field := p.Parser.parseField()
		if _, ok := node.Fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		node.Fields[field.Name] = field
	}

	p.AddFragment(node)
}

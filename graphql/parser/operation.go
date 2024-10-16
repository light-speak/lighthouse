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


func (p *QueryParser) parseFragment() {
	node := &ast.FragmentNode{
		Name: p.Parser.expectAndGetValue(lexer.Fragment),
		On:   p.Parser.expectAndGetValue(lexer.On),
	}

	p.Parser.expect(lexer.LeftBrace)
	node.Fields = make(map[string]*ast.Field)
	for p.Parser.currToken.Type != lexer.RightBrace {
		field := p.Parser.parseField()
		if _, ok := node.Fields[field.Name]; ok {
			panic("duplicate field: " + field.Name)
		}
		node.Fields[field.Name] = field
	}

	p.AddFragment(node)
}
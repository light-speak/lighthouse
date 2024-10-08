package parser

import (
	"strconv"

	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/lexer"
	"github.com/light-speak/lighthouse/parser/value"
)

// parseArguments parse arguments
// (id: ID!, name: String!)
func (p *Parser) parseArguments(parent ast.Node) []*ast.ArgumentNode {
	var args []*ast.ArgumentNode
	if p.currToken.Type != lexer.LeftParent {
		return args
	}
	p.expect(lexer.LeftParent)
	for p.currToken.Type != lexer.RightParent {
		args = append(args, p.parseArgument(parent))
		if p.currToken.Type != lexer.RightParent {
			p.expect(lexer.Comma)
		}
	}
	p.expect(lexer.RightParent)
	return args
}

// parseDefaultValue parse a default value
// = 123, = "123", = true, = false, = null
// = [123, 456] , = ["123", "456"]
func (p *Parser) parseDefaultValue() *ast.ArgumentValue {
	if p.currToken.Type == lexer.Equal {
		p.expect(lexer.Equal) // skip =
		return p.parseArgumentValue()
	}
	return nil
}

// parseArgument parse an argument node
func (p *Parser) parseArgument(parent ast.Node) *ast.ArgumentNode {
	description := p.parseDescription()
	name := p.currToken.Value

	p.nextToken()         // skip name
	p.expect(lexer.Colon) // skip :
	var fieldType *ast.FieldType
	var defaultValue *ast.ArgumentValue
	var value *ast.ArgumentValue

	if parent.GetType() == ast.NodeTypeDirective {
		// Assigned when using @directive
		fieldType = nil
		defaultValue = nil
		value = p.parseArgumentValue()
	} else {
		fieldType = p.parseTypeReference()   // parse type reference
		defaultValue = p.parseDefaultValue() // parse default value
		value = nil
	}

	// parse directives
	directives := p.parseDirectives()

	return &ast.ArgumentNode{
		Name:         name,
		Type:         fieldType,
		Description:  description,
		Value:        value,
		Directives:   directives,
		DefaultValue: defaultValue,
		Parent:       parent,
	}
}

// parseArgumentValue parse a directive argument value
// @directive(arg: "123")
// @directive(arg: 123)
// @directive(arg: [123, 456])
// @directive(arg: ["123", "456"])
// @directive(arg: true, arg2: false)
// @directive(arg: Boolean, arg2: String, arg3: Int, arg4: [[User]!]!, arg5: ID)
// The colon has been parsed in the previous step, so only the value needs to be parsed here
func (p *Parser) parseArgumentValue() *ast.ArgumentValue {
	var argValue *ast.ArgumentValue

	if p.currToken.Type == lexer.LeftBracket {
		p.expect(lexer.LeftBracket) // skip [
		values := []*ast.ArgumentValue{}

		for p.currToken.Type != lexer.RightBracket {
			values = append(values, p.parseArgumentValue())
			if p.currToken.Type == lexer.Comma {
				p.expect(lexer.Comma) // skip ,
			}
		}
		p.expect(lexer.RightBracket) // skip ]

		argValue = &ast.ArgumentValue{
			Children: values,
			Type: &ast.FieldType{
				Name:   "List",
				IsList: true,
			},
		}
	} else {
		var values []*ast.ArgumentValue
		for {
			var v value.Value
			switch p.currToken.Type {
			case lexer.Letter:
				v = &value.StringValue{Value: p.currToken.Value}
			case lexer.IntNumber:
				intValue, err := strconv.ParseInt(p.currToken.Value, 10, 64)
				if err != nil {
					panic("invalid integer value: " + err.Error())
				}
				v = &value.IntValue{Value: intValue}
			case lexer.Boolean:
				boolValue := p.currToken.Value == "true"
				v = &value.BooleanValue{Value: boolValue}
			default:
				panic("unsupported token type: " + p.currToken.Type)
			}
			vi := &ast.ArgumentValue{
				Value: v,
				Type: &ast.FieldType{
					Name: string(p.currToken.Type),
				},
			}
			values = append(values, vi)
			p.nextToken()

			if p.currToken.Type != lexer.Comma {
				break
			}
			p.expect(lexer.Comma)
		}

		if len(values) == 1 {
			argValue = values[0]
		} else {
			argValue = &ast.ArgumentValue{
				Children: values,
				Type: &ast.FieldType{
					IsList: true,
					Name:   "List",
				},
			}
		}
	}

	if p.currToken.Type == lexer.Exclamation {
		argValue.Type.IsNonNull = true
		p.expect(lexer.Exclamation)
	}
	return argValue
}

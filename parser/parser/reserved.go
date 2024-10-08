package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/scalar"
)

func (p *Parser) AddReserved() {
	p.addReservedScalarType()
	p.addReservedScalar()
}

func (p *Parser) MergeScalarType() {
	for name, scalar := range p.ScalarMap {
		if _, ok := p.ScalarTypeMap[name]; ok {
			scalar.Scalar = p.ScalarTypeMap[name]
		}
	}
}

func (p *Parser) addReservedScalar() {
	p.AddScalar(&ast.ScalarNode{
		Name:        "Boolean",
		Description: "The Boolean scalar type represents a boolean value. It can be either true or false.",
	})

	p.AddScalar(&ast.ScalarNode{
		Name:        "Int",
		Description: "The Int scalar type represents a signed 32-bit integer.",
	})

	p.AddScalar(&ast.ScalarNode{
		Name:        "Float",
		Description: "The Float scalar type represents a signed double-precision floating-point number.",
	})

	p.AddScalar(&ast.ScalarNode{
		Name:        "String",
		Description: "The String scalar type represents a string value.",
	})

	p.AddScalar(&ast.ScalarNode{
		Name:        "ID",
		Description: "The ID scalar type represents a unique identifier.",
	})
}

func (p *Parser) addReservedScalarType() {
	p.AddScalarType("Boolean", &scalar.BooleanScalar{})
	p.AddScalarType("Int", &scalar.IntScalar{})
	p.AddScalarType("Float", &scalar.FloatScalar{})
	p.AddScalarType("String", &scalar.StringScalar{})
	p.AddScalarType("ID", &scalar.IDScalar{})
}

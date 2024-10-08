package parser

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/scalar"
)

func (p *Parser) AddReserved() {
	p.addReservedScalarType()
	p.addReservedScalar()
	p.addReservedDirective()
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
		BaseNode: ast.BaseNode{
			Name:        "Boolean",
			Description: "The Boolean scalar type represents a boolean value. It can be either true or false.",
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "Int",
			Description: "The Int scalar type represents a signed 32-bit integer.",
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "Float",
			Description: "The Float scalar type represents a signed double-precision floating-point number.",
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "String",
			Description: "The String scalar type represents a string value.",
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "ID",
			Description: "The ID scalar type represents a unique identifier.",
		},
	})
}

func (p *Parser) addReservedDirective() {
	p.AddDirective(&ast.DirectiveDefinitionNode{
		BaseNode: ast.BaseNode{
			Name:        "skip",
			Description: "Skips the current field or fragment when the argument is true.",
		},
		Locations: []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
		Args: []*ast.ArgumentNode{
			{
				BaseNode: ast.BaseNode{
					Name: "if",
				},
				Type: &ast.FieldType{
					Name:      "Boolean",
					Type:      p.ScalarMap["Boolean"],
					IsNonNull: false,
				},
			},
		},
	})

	p.AddDirective(&ast.DirectiveDefinitionNode{
		BaseNode: ast.BaseNode{
			Name:        "include",
			Description: "Includes the current field or fragment when the argument is true.",
		},
		Locations: []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
		Args: []*ast.ArgumentNode{
			{
				BaseNode: ast.BaseNode{
					Name: "if",
				},
				Type: &ast.FieldType{
					Name:      "Boolean",
					Type:      p.ScalarMap["Boolean"],
					IsNonNull: false,
				},
			},
		},
	})

	p.AddDirective(&ast.DirectiveDefinitionNode{
		BaseNode: ast.BaseNode{
			Name:        "enum",
			Description: "enum",
		},
		Locations: []ast.Location{ast.LocationEnumValue},
		Args: []*ast.ArgumentNode{
			{
				BaseNode: ast.BaseNode{
					Name: "value",
				},
				Type: &ast.FieldType{
					Name:      "Int",
					Type:      p.ScalarMap["Int"],
					IsNonNull: true,
				},
			},
		},
	})
}

func (p *Parser) addReservedScalarType() {
	p.AddScalarType("Boolean", &scalar.BooleanScalar{})
	p.AddScalarType("Int", &scalar.IntScalar{})
	p.AddScalarType("Float", &scalar.FloatScalar{})
	p.AddScalarType("String", &scalar.StringScalar{})
	p.AddScalarType("ID", &scalar.IDScalar{})
}

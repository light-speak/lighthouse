package parser

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/scalar"
	"github.com/light-speak/lighthouse/utils"
)

func (p *Parser) AddReserved() {
	p.addReservedScalarType()
	p.addReservedScalar()
	p.addReservedDirective()
	p.addReservedEnum()
	p.addReservedObject()
}

func (p *Parser) MergeScalarType() {
	for name, scalar := range p.NodeStore.Scalars {
		if _, ok := p.NodeStore.ScalarTypes[name]; ok {
			scalar.ScalarType = p.NodeStore.ScalarTypes[name]
		}
	}
}

func (p *Parser) addReservedScalar() {
	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "Boolean",
			Description: utils.StrPtr("The Boolean scalar type represents a boolean value. It can be either true or false."),
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "Int",
			Description: utils.StrPtr("The Int scalar type represents a signed 32-bit integer."),
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "Float",
			Description: utils.StrPtr("The Float scalar type represents a signed double-precision floating-point number."),
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "String",
			Description: utils.StrPtr("The String scalar type represents a string value."),
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "ID",
			Description: utils.StrPtr("The ID scalar type represents a unique identifier."),
		},
	})

	p.AddScalar(&ast.ScalarNode{
		BaseNode: ast.BaseNode{
			Name:        "DateTime",
			Description: utils.StrPtr("The DateTime scalar type represents a date and time."),
		},
	})
}

func (p *Parser) addReservedScalarType() {
	p.AddScalarType("Boolean", &scalar.BooleanScalar{})
	p.AddScalarType("Int", &scalar.IntScalar{})
	p.AddScalarType("Float", &scalar.FloatScalar{})
	p.AddScalarType("String", &scalar.StringScalar{})
	p.AddScalarType("ID", &scalar.IDScalar{})
	p.AddScalarType("DateTime", &scalar.DateTimeScalar{})
}

func (p *Parser) addReservedObject() {
	p.AddObject(&ast.ObjectNode{
		BaseNode: ast.BaseNode{
			IsReserved:  true,
			Name:        "PaginateInfo",
			Description: utils.StrPtr("The PaginateInfo type represents information about a paginated list."),
		},
		Fields: map[string]*ast.Field{
			"currentPage": {
				Name: "currentPage",
				Type: &ast.TypeRef{
					Kind: ast.KindNonNull,
					OfType: &ast.TypeRef{
						Kind:     ast.KindScalar,
						Name:     "Int",
						TypeNode: p.NodeStore.Scalars["Int"],
					},
				},
			},
			"totalPage": {
				Name: "totalPage",
				Type: &ast.TypeRef{
					Kind: ast.KindNonNull,
					OfType: &ast.TypeRef{
						Kind:     ast.KindScalar,
						Name:     "Int",
						TypeNode: p.NodeStore.Scalars["Int"],
					},
				},
			},
			"hasNextPage": {
				Name: "hasNextPage",
				Type: &ast.TypeRef{
					Kind: ast.KindNonNull,
					OfType: &ast.TypeRef{
						Kind:     ast.KindScalar,
						Name:     "Boolean",
						TypeNode: p.NodeStore.Scalars["Boolean"],
					},
				},
			},
			"totalCount": {
				Name: "totalCount",
				Type: &ast.TypeRef{
					Kind: ast.KindNonNull,
					OfType: &ast.TypeRef{
						Kind:     ast.KindScalar,
						Name:     "Int",
						TypeNode: p.NodeStore.Scalars["Int"],
					},
				},
			},
		},
	}, false)
}

func (p *Parser) addReservedEnum() {
	p.AddEnum(&ast.EnumNode{
		BaseNode: ast.BaseNode{
			IsReserved:  true,
			Name:        "SortOrder",
			Description: utils.StrPtr("The SortOrder enum type represents the order of a list."),
		},
		EnumValues: map[string]*ast.EnumValue{
			"ASC": {
				Name:        "ASC",
				Description: utils.StrPtr("The ASC enum value represents ascending order."),
				Value:       int8(1),
			},
			"DESC": {
				Name:        "DESC",
				Description: utils.StrPtr("The DESC enum value represents descending order."),
				Value:       int8(-1),
			},
		},
	}, false)
}

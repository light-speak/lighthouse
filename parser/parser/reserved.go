package parser

import (
	"github.com/light-speak/lighthouse/parser/ast"
	"github.com/light-speak/lighthouse/parser/scalar"
)

func addReservedScalar(p *Parser) {
	// Define reserved scalars
	reservedScalars := []struct {
		name        string
		description string
		scalar      ast.ScalarType
	}{
		{"ID", "The ID scalar type represents a unique identifier for a resource.", &scalar.IDScalar{}},
		{"String", "The String scalar type represents a sequence of characters.", &scalar.StringScalar{}},
		{"Int", "The Int scalar type represents a signed 64-bit integer.", &scalar.IntScalar{}},
		{"Float", "The Float scalar type represents a signed double-precision floating-point number.", &scalar.FloatScalar{}},
		{"Boolean", "The Boolean scalar type represents a boolean value.", &scalar.BooleanScalar{}},
	}

	// Populate ScalarMap
	for _, scalar := range reservedScalars {
		p.AddScalar(&ast.ScalarNode{
			Name:        scalar.name,
			Description: scalar.description,
			Scalar:      scalar.scalar,
		})
	}
}

func addReservedDirective(p *Parser) {
	reservedDirectives := []struct {
		name        string
		description string
		locations   []ast.Location
		args        []*ast.ArgumentNode
	}{
		{
			"scalar",
			"The scalar directive is used to define a custom scalar type.",
			[]ast.Location{ast.LocationScalar},
			[]*ast.ArgumentNode{
				{
					Name:        "name",
					Description: "The name of the scalar",
					Type: &ast.FieldType{
						Name: "String",
					},
				},
			},
		},
		{
			"include",
			"The include directive is used to conditionally include fields in the response.",
			[]ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
			[]*ast.ArgumentNode{
				{
					Name:        "if",
					Description: "The condition to include the field",
					Type: &ast.FieldType{
						Name: "Boolean",
					},
				},
			},
		},
		{
			"skip",
			"The skip directive is used to conditionally skip fields in the response.",
			[]ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
			[]*ast.ArgumentNode{
				{
					Name:        "if",
					Description: "The condition to skip the field",
					Type: &ast.FieldType{
						Name: "Boolean",
					},
				},
			},
		},
		{
			"paginate",
			"The paginate directive is used to paginate the results of a query.",
			[]ast.Location{ast.LocationField},
			[]*ast.ArgumentNode{
				{
					Name:        "scopes",
					Description: "The scope of the pagination",
					Type: &ast.FieldType{
						Name:   "List",
						IsList: true,
						ElemType: &ast.FieldType{
							Name:         "String",
							Type:         p.ScalarMap["String"],
							TypeCategory: ast.TypeCategoryScalar,
						},
					},
				},
			},
		},
	}

	for _, directive := range reservedDirectives {
		p.AddDirective(&ast.DirectiveDefinitionNode{
			Name:        directive.name,
			Description: directive.description,
			Locations:   directive.locations,
			Args:        directive.args,
		})
	}
}

func (p *Parser) AddReserved() {
	addReservedScalar(p)
	addReservedDirective(p)
}

package parser

import "github.com/light-speak/lighthouse/graphql/ast"

func (p *Parser) addReservedDirective() {
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "skip", Description: "skip current field or fragment, when the parameter is true.",
		Locations: []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
		Args: map[string]*ast.Argument{
			"if": {
				Name: "if",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "Boolean",
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "include", Description: "include current field or fragment, when the parameter is true.",
		Locations: []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
		Args: map[string]*ast.Argument{
			"if": {
				Name: "if",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "Boolean",
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "enum", Description: "Change the value of the enum.",
		Locations: []ast.Location{ast.LocationEnumValue},
		Args: map[string]*ast.Argument{
			"value": {
				Name: "value",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "Int",
					OfType: &ast.TypeRef{
						Kind: ast.KindNonNull,
						OfType: &ast.TypeRef{
							Kind: ast.KindScalar,
							Name: "Int",
						},
					},
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "paginate", Description: "The response will return paginate information and a list. The field must be in the form of a list.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"scopes": {
				Name: "scopes",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
					OfType: &ast.TypeRef{
						Kind: ast.KindList,
						OfType: &ast.TypeRef{
							Kind: ast.KindScalar,
							Name: "String",
						},
					},
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "external", Description: "The field is defined in another schema.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "requires", Description: "The field is defined in another schema.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "provides", Description: "The field is defined in another schema.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "key", Description: "The field is defined in another schema.",
		Locations: []ast.Location{ast.LocationObject, ast.LocationInterface},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "extends", Description: "The field is defined in another schema.",
		Locations: []ast.Location{ast.LocationObject},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "softDeleteModel", Description: "The model is soft delete.",
		Locations: []ast.Location{ast.LocationObject},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "model", Description: "The model name.",
		Locations: []ast.Location{ast.LocationObject},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "tag", Description: "The tag of the field.",
		Locations:  []ast.Location{ast.LocationFieldDefinition},
		Repeatable: true,
		Args: map[string]*ast.Argument{
			"name": {
				Name: "name",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
			"value": {
				Name: "value",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "index", Description: "The field is indexed.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"name": {
				Name: "name",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "unique", Description: "The field is unique.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "default", Description: "The default value of the field.",
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"value": {
				Name: "value",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
}

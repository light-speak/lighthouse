package parser

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func (p *Parser) addReservedDirective() {
	// deprecated
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "deprecated", Description: utils.StrPtr("The field is deprecated."),
		Locations: []ast.Location{
			ast.LocationFieldDefinition,
			ast.LocationField,
			ast.LocationEnumValue,
			ast.LocationInputFieldDefinition,
		},
		Args: map[string]*ast.Argument{
			"reason": {
				Name: "reason",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// skip
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "skip", Description: utils.StrPtr("skip current field or fragment, when the parameter is true."),
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
	// include
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "include", Description: utils.StrPtr("include current field or fragment, when the parameter is true."),
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
	// enum
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "enum", Description: utils.StrPtr("Change the value of the enum."),
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
	// paginate
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "paginate", Description: utils.StrPtr("The response will return paginate information and a list. The field must be in the form of a list."),
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
	// external
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name:      "external",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	// requires
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name:      "requires",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	// provides
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name:      "provides",
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	// key
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name:      "key",
		Locations: []ast.Location{ast.LocationObject, ast.LocationInterface},
		Args: map[string]*ast.Argument{
			"fields": {
				Name: "fields",
				Type: &ast.TypeRef{
					Kind: ast.KindList,
					OfType: &ast.TypeRef{
						Kind: ast.KindScalar,
						Name: "String",
					},
				},
			},
		},
	})
	// extends
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name:      "extends",
		Locations: []ast.Location{ast.LocationObject},
	})
	// softDeleteModel
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "softDeleteModel", Description: utils.StrPtr("The model is soft delete."),
		Locations: []ast.Location{ast.LocationObject},
	})
	// model
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "model", Description: utils.StrPtr("The model."),
		Locations: []ast.Location{ast.LocationObject},
	})
	// tag
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "tag", Description: utils.StrPtr("The tag of the field."),
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
	// index
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "index", Description: utils.StrPtr("The field is indexed."),
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
	// unique
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "unique", Description: utils.StrPtr("The field is unique."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})
	// default
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "default", Description: utils.StrPtr("The default value of the field."),
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
	// first
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "first", Description: utils.StrPtr("The response will return only one item."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
	})



}

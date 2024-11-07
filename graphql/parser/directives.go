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

	p.addFilterDirective()
	p.addReturnDirective()
	p.addRelationDirective()
	p.addObjectDirective()

	p.addRuntimeFieldDirective()
}

func (p *Parser) addReturnDirective() {
	// first
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "first", Description: utils.StrPtr("The response will return only one item."),
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
	// find
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "find", Description: utils.StrPtr("The response will return only one item."),
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
}

func (p *Parser) addFilterDirective() {
	// eq
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "eq", Description: utils.StrPtr("The field is equal to the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// neq
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "neq", Description: utils.StrPtr("The field is not equal to the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// gt
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "gt", Description: utils.StrPtr("The field is greater than the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// gte
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "gte", Description: utils.StrPtr("The field is greater than or equal to the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// lt
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "lt", Description: utils.StrPtr("The field is less than the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// lte
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "lte", Description: utils.StrPtr("The field is less than or equal to the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// in
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "in", Description: utils.StrPtr("The field is in the value list."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// notIn
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "notIn", Description: utils.StrPtr("The field is not in the value list."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
	// like
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "like", Description: utils.StrPtr("The field is like the value."),
		Locations: []ast.Location{ast.LocationArgumentDefinition},
		Args: map[string]*ast.Argument{
			"field": {
				Name: "field",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "String",
				},
			},
		},
	})
}

func (p *Parser) addRelationDirective() {
	// belongsTo
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "belongsTo", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"relation": {
				Name: "relation",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"foreignKey": {
				Name: "foreignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"reference": {
				Name: "reference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// hasMany
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "hasMany", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"relation": {
				Name: "relation",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"foreignKey": {
				Name: "foreignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"reference": {
				Name: "reference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// hasOne
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "hasOne", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"relation": {
				Name: "relation",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"foreignKey": {
				Name: "foreignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"reference": {
				Name: "reference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// morphTo
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "morphTo", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"morphType": {
				Name: "morphType",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"morphKey": {
				Name: "morphKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"reference": {
				Name: "reference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// morphToMany
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "morphToMany", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"relation": {
				Name: "relation",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"currentType": {
				Name: "currentType",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"morphType": {
				Name: "morphType",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"morphKey": {
				Name: "morphKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"reference": {
				Name: "reference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// manyToMany
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "manyToMany", Description: utils.StrPtr("The field is a relationship with another model."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"relation": {
				Name: "relation",
				Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"}},
			},
			"pivot": {
				Name: "pivot",
				Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"}},
			},
			"currentType": {
				Name: "currentType",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"pivotForeignKey": {
				Name: "pivotForeignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"foreignKey": {
				Name: "foreignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"pivotReference": {
				Name: "pivotReference",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
			"relationForeignKey": {
				Name: "relationForeignKey",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
}

func (p *Parser) addObjectDirective() {
	// softDeleteModel
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "softDeleteModel", Description: utils.StrPtr("The model is soft delete."),
		Locations: []ast.Location{ast.LocationObject},
		Args: map[string]*ast.Argument{
			"table": {
				Name: "table",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
	})
	// model
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "model", Description: utils.StrPtr("The model."),
		Locations: []ast.Location{ast.LocationObject},
		Args: map[string]*ast.Argument{
			"table": {
				Name: "table",
				Type: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
			},
		},
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
			"unique": {
				Name: "unique",
				Type: &ast.TypeRef{
					Kind: ast.KindScalar,
					Name: "Boolean",
				},
				DefaultValue: false,
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
}

func (p *Parser) addRuntimeFieldDirective() {
	// auth
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "auth", Description: utils.StrPtr("The field is runtime auth."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"msg": {
				Name:         "msg",
				Description:  utils.StrPtr("The message of the auth."),
				Type:         &ast.TypeRef{Kind: ast.KindScalar, Name: "String"},
				DefaultValue: "Unauthorized! please login",
			},
		},
	})
	// cache
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "cache", Description: utils.StrPtr("The field is cached."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"auth": {
				Name:        "auth",
				Description: utils.StrPtr("Cache with auth info. if true, the cache will be valid for the current user."),
				Type:        &ast.TypeRef{Kind: ast.KindScalar, Name: "Boolean"},
			},
			"ttl": {
				Name:        "ttl",
				Description: utils.StrPtr("Cache ttl."),
				Type:        &ast.TypeRef{Kind: ast.KindScalar, Name: "Int"},
			},
			"tags": {
				Name:        "tags",
				Description: utils.StrPtr("Cache tags."),
				Type:        &ast.TypeRef{Kind: ast.KindList, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"}},
			},
		},
	})
	// cacheClear
	p.AddDirectiveDefinition(&ast.DirectiveDefinition{
		Name: "cacheClear", Description: utils.StrPtr("Clear the cache."),
		Locations: []ast.Location{ast.LocationFieldDefinition},
		Args: map[string]*ast.Argument{
			"tags": {
				Name:        "tags",
				Description: utils.StrPtr("Cache tags."),
				Type:        &ast.TypeRef{Kind: ast.KindList, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "String"}},
			},
			"auth": {
				Name:        "auth",
				Description: utils.StrPtr("Clear the cache with auth info. if true, the cache will be cleared only for the current user."),
				Type:        &ast.TypeRef{Kind: ast.KindScalar, Name: "Boolean"},
			},
		},
	})
}

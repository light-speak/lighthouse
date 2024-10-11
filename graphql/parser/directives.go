package parser

import "github.com/light-speak/lighthouse/graphql/ast"

func (p *Parser) addReservedDirective() {
	directives := []struct {
		name        string
		description string
		locations   []ast.Location
		args        []*ast.ArgumentNode
		repeatable  bool
	}{
		{
			name:        "skip",
			description: "skip current field or fragment, when the parameter is true.",
			locations:   []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "if"},
				Type:     &ast.FieldType{Name: "Boolean", Type: p.ScalarMap["Boolean"]},
			}},
		},
		{
			name:        "include",
			description: "include current field or fragment, when the parameter is true.",
			locations:   []ast.Location{ast.LocationField, ast.LocationFragmentSpread, ast.LocationInlineFragment},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "if"},
				Type:     &ast.FieldType{Name: "Boolean", Type: p.ScalarMap["Boolean"]},
			}},
		},
		{
			name:        "enum",
			description: "Change the value of the enum.",
			locations:   []ast.Location{ast.LocationEnumValue},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "value"},
				Type:     &ast.FieldType{Name: "Int", Type: p.ScalarMap["Int"], IsNonNull: true},
			}},
		},
		{
			name:        "paginate",
			description: "The response will return paginate information and a list. The field must be in the form of a list.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "scopes"},
				Type: &ast.FieldType{
					Name:   "Int",
					IsList: true,
					ElemType: &ast.FieldType{
						Name:      "String",
						Type:      p.ScalarMap["String"],
						IsNonNull: true,
					},
				},
			}},
		},
		{
			name:        "external",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
		},
		{
			name:        "requires",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "fields"},
				Type: &ast.FieldType{
					Name:   "String",
					IsList: true,
					ElemType: &ast.FieldType{
						Name:      "String",
						Type:      p.ScalarMap["String"],
						IsNonNull: true,
					},
				},
			}},
		},
		{
			name:        "provides",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "fields"},
				Type: &ast.FieldType{
					Name:   "String",
					IsList: true,
					ElemType: &ast.FieldType{
						Name:      "String",
						Type:      p.ScalarMap["String"],
						IsNonNull: true,
					},
				},
			}},
		},
		{
			name:        "key",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationObject, ast.LocationInterface},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "fields"},
				Type: &ast.FieldType{
					Name:   "String",
					IsList: true,
					ElemType: &ast.FieldType{
						Name:      "String",
						Type:      p.ScalarMap["String"],
						IsNonNull: true,
					},
				},
			}},
		},
		{
			name:        "extends",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationObject},
		},
		{
			name:        "softDeleteModel",
			description: "The model is soft delete.",
			locations:   []ast.Location{ast.LocationObject},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "name"},
				Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: false},
			}},
		},
		{
			name:        "model",
			description: "The model name.",
			locations:   []ast.Location{ast.LocationObject},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "name"},
				Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: false},
			}},
		},
		{
			name:        "tag",
			description: "The tag of the field.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			repeatable:  true,
			args: []*ast.ArgumentNode{
				{
					BaseNode: ast.BaseNode{Name: "name"},
					Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: true},
				},
				{
					BaseNode: ast.BaseNode{Name: "value"},
					Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: true},
				},
			},
		},
		{
			name:        "index",
			description: "The field is indexed.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "name", Description: "The name of the index."},
				Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: false},
			}},
		},
		{
			name:        "unique",
			description: "The field is unique.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
		},
		{
			name:        "defaultString",
			description: "The default value of the field.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "value"},
				Type:     &ast.FieldType{Name: "String", Type: p.ScalarMap["String"], IsNonNull: true},
			}},
		},
		{
			name:        "defaultInt",
			description: "The default value of the field.",
			locations:   []ast.Location{ast.LocationFieldDefinition},
			args: []*ast.ArgumentNode{{
				BaseNode: ast.BaseNode{Name: "value"},
				Type:     &ast.FieldType{Name: "Int", Type: p.ScalarMap["Int"], IsNonNull: true},
			}},
		},
	}

	for _, d := range directives {
		p.AddDirective(&ast.DirectiveDefinitionNode{
			BaseNode:   ast.BaseNode{Name: d.name, Description: d.description},
			Locations:  d.locations,
			Args:       d.args,
			Repeatable: d.repeatable,
		})
	}
}

package parser

import "github.com/light-speak/lighthouse/graphql/ast"

func (p *Parser) addReservedDirective() {
	directives := []struct {
		name        string
		description string
		locations   []ast.Location
		args        []*ast.ArgumentNode
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
			name:        "softDelete",
			description: "The field is defined in another schema.",
			locations:   []ast.Location{ast.LocationObject},
		},
	}

	for _, d := range directives {
		p.AddDirective(&ast.DirectiveDefinitionNode{
			BaseNode:  ast.BaseNode{Name: d.name, Description: d.description},
			Locations: d.locations,
			Args:      d.args,
		})
	}
}

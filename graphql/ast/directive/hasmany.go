package directive

import (
		"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func handlerHasMany(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	relation := &ast.Relation{
		Name:         d.GetArg("relation").Value.(string),
		RelationType: ast.RelationTypeHasMany,
	}
	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Name = relationName.Value.(string)
	} else {
		return &errors.GraphQLError{
			Message:   "relation name is required for hasMany directive",
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
		}
	}
	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = foreignKey.Value.(string)
	} else {
		return &errors.GraphQLError{
			Message:   "foreign key is required for hasMany directive",
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
		}
	}
	if reference := d.GetArg("reference"); reference != nil {
		relation.Reference = reference.Value.(string)
	} else {
		relation.Reference = "id"
	}
	f.Relation = relation
	return nil
}

func init() {
	ast.AddFieldDirective("hasMany", handlerHasMany)
}

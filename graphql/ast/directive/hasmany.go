package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func handlerHasMany(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	relation := &ast.Relation{
		Name:         d.GetArg("relation").Value.(string),
		RelationType: ast.RelationTypeHasMany,
	}
	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Name = utils.SnakeCase(relationName.Value.(string))
	} else {
		relation.Name = utils.SnakeCase(f.Type.GetRealType().Name)
	}
	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = utils.SnakeCase(foreignKey.Value.(string))
	} else {
		relation.ForeignKey = utils.SnakeCase(relation.Name) + "_id"
	}
	if reference := d.GetArg("reference"); reference != nil {
		relation.Reference = utils.SnakeCase(reference.Value.(string))
	} else {
		relation.Reference = "id"
	}
	f.Relation = relation
	return nil
}

func init() {
	ast.AddFieldDirective("hasMany", handlerHasMany)
}

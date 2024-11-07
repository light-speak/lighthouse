package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func handlerManyToMany(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	relation := &ast.Relation{
		RelationType: ast.RelationTypeBelongsToMany,
	}

	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Name = utils.SnakeCase(relationName.Value.(string))
	} else {
		return &errors.GraphQLError{
			Message:   "relation argument is required for @manyToMany directive",
			Locations: []*errors.GraphqlLocation{d.GetLocation()},
		}
	}

	if pivot := d.GetArg("pivot"); pivot != nil {
		relation.Pivot = utils.SnakeCase(pivot.Value.(string))
	} else {
		relation.Pivot = utils.SnakeCase(parent.GetName() + "_" + f.Type.GetRealType().Name)
	}

	if currentType := d.GetArg("currentType"); currentType != nil {
		relation.CurrentType = utils.SnakeCase(currentType.Value.(string))
	} else {
		relation.CurrentType = utils.SnakeCase(parent.GetName())
	}

	if pivotForeignKey := d.GetArg("pivotForeignKey"); pivotForeignKey != nil {
		relation.PivotForeignKey = utils.SnakeCase(pivotForeignKey.Value.(string))
	} else {
		relation.PivotForeignKey = utils.SnakeCase(parent.GetName()) + "_id"
	}

	if pivotReference := d.GetArg("pivotReference"); pivotReference != nil {
		relation.PivotReference = utils.SnakeCase(pivotReference.Value.(string))
	} else {
		relation.PivotReference = utils.SnakeCase(f.Type.GetRealType().Name) + "_id"
	}

	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = utils.SnakeCase(foreignKey.Value.(string))
	} else {
		relation.ForeignKey = "id"
	}

	if relationForeignKey := d.GetArg("relationForeignKey"); relationForeignKey != nil {
		relation.RelationForeignKey = utils.SnakeCase(relationForeignKey.Value.(string))
	} else {
		relation.RelationForeignKey = "id"
	}

	f.Relation = relation
	return nil
}

func init() {
	ast.AddFieldDirective("manyToMany", handlerManyToMany)
}

package directive

import (
	"fmt"

	"github.com/light-speak/lighthouse/graphql/ast"
)

func handlerHasMany(f *ast.Field, d *ast.Directive, store *ast.NodeStore) error {
	relation := &ast.Relation{
		Relation:     d.GetArg("relation").Value.(string),
		RelationType: ast.RelationTypeHasMany,
	}
	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Relation = relationName.Value.(string)
	} else {
		return fmt.Errorf("relation name is required for hasMany directive")
	}
	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = foreignKey.Value.(string)
	} else {
		return fmt.Errorf("foreign key is required for hasMany directive")
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

package directive

import (
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func handlerBelongsTo(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) error {
	relation := &ast.Relation{
		RelationType: ast.RelationTypeBelongsTo,
	}
	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Name = relationName.Value.(string)
	} else {
		relation.Name = utils.LcFirst(f.Name)
	}
	if reference := d.GetArg("reference"); reference != nil {
		relation.Reference = reference.Value.(string)
	} else {
		relation.Reference = "id"
	}
	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = foreignKey.Value.(string)
	} else {
		relation.ForeignKey = utils.LcFirst(f.Name) + "_id"
	}
	f.Relation = relation
	return nil
}

func init() {
	ast.AddFieldDirective("belongsTo", handlerBelongsTo)
}

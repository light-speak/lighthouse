package directive

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func handlerHasOne(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	relation := &ast.Relation{
		RelationType: ast.RelationTypeHasOne,
	}

	// 处理关系名称
	if relationName := d.GetArg("relation"); relationName != nil {
		relation.Name = utils.SnakeCase(relationName.Value.(string))
	} else {
		relation.Name = utils.LcFirst(f.Name)
	}

	// 处理外键
	if foreignKey := d.GetArg("foreignKey"); foreignKey != nil {
		relation.ForeignKey = utils.SnakeCase(foreignKey.Value.(string))
	} else {
		// 默认外键是父对象名称加_id
		relation.ForeignKey = utils.SnakeCase(parent.GetName()) + "_id"
	}

	// 处理引用键
	if reference := d.GetArg("reference"); reference != nil {
		relation.Reference = utils.SnakeCase(reference.Value.(string))
	} else {
		relation.Reference = "id"
	}

	f.Relation = relation
	return nil
}

func init() {
	ast.AddFieldDirective("hasOne", handlerHasOne)
}

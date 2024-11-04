package directive

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/utils"
)

func handlerMorphTo(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) errors.GraphqlErrorInterface {
	relation := &ast.Relation{
		RelationType: ast.RelationTypeMorphTo,
	}
	if morphType := d.GetArg("morphType"); morphType != nil {
		relation.MorphType = morphType.Value.(string)
	} else {
		relation.MorphType = fmt.Sprintf("%s_type", utils.LcFirst(f.Name))
	}
	if morphKey := d.GetArg("morphKey"); morphKey != nil {
		relation.MorphKey = morphKey.Value.(string)
	} else {
		relation.MorphKey = fmt.Sprintf("%s_id", utils.LcFirst(f.Name))
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
	ast.AddFieldDirective("morphTo", handlerMorphTo)
}

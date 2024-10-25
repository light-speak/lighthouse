package directive

import "github.com/light-speak/lighthouse/graphql/ast"

func handlerSoftDeleteModel(o *ast.ObjectNode, d *ast.Directive, store *ast.NodeStore) error {
	o.IsModel = true
	o.Fields["id"] = &ast.Field{
		Name: "id",
		Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "ID", TypeNode: store.Scalars["ID"]}},
	}
	o.Fields["created_at"] = &ast.Field{
		Name: "created_at",
		Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "DateTime", TypeNode: store.Scalars["DateTime"]}},
	}
	o.Fields["updated_at"] = &ast.Field{
		Name: "updated_at",
		Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "DateTime", TypeNode: store.Scalars["DateTime"]}},
	}
	o.Fields["deleted_at"] = &ast.Field{
		Name: "deleted_at",
		Type: &ast.TypeRef{Kind: ast.KindNonNull, OfType: &ast.TypeRef{Kind: ast.KindScalar, Name: "DateTime", TypeNode: store.Scalars["DateTime"]}},
	}
	return nil
}

func init() {
	ast.AddObjectDirective("softDeleteModel", handlerSoftDeleteModel)
}

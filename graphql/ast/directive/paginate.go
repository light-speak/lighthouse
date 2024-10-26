package directive

import (
	"fmt"

	"github.com/light-speak/lighthouse/graphql/ast"
)

func handlerPaginate(f *ast.Field, d *ast.Directive, store *ast.NodeStore, parent ast.Node) error {
	if parent.GetName() != "Query" {
		return fmt.Errorf("paginate directive can only be used on Query type")
	}
	addPaginationResponseType(f, store)
	addPaginationArguments(f, store)
	return nil
}

func init() {
	ast.AddFieldDirective("paginate", handlerPaginate)
}

// addPaginationResponseType adds a pagination response type to the store
func addPaginationResponseType(f *ast.Field, store *ast.NodeStore) {
	curType := f.Type
	for curType.Kind == ast.KindList || curType.Kind == ast.KindNonNull {
		curType = curType.OfType
	}
	typeName := curType.Name
	responseName := fmt.Sprintf("%sPaginateResponse", typeName)
	if _, ok := store.Objects[responseName]; ok {
		return
	}
	description := fmt.Sprintf("The %sPaginateResponse type represents a paginated list of %s.", typeName, typeName)
	store.AddObject(responseName, &ast.ObjectNode{
		BaseNode: ast.BaseNode{
			Name:        responseName,
			Kind:        ast.KindObject,
			Description: &description,
		},
		Fields: map[string]*ast.Field{
			"data": {
				Name: "data",
				Type: f.Type,
			},
			"paginateInfo": {
				Name: "paginateInfo",
				Type: &ast.TypeRef{
					Kind: ast.KindNonNull,
					OfType: &ast.TypeRef{
						Kind:     ast.KindObject,
						Name:     "PaginateInfo",
						TypeNode: store.Objects["PaginateInfo"],
					},
				},
			},
		},
	})
	f.Type = &ast.TypeRef{
		Kind:     ast.KindNonNull,
		OfType: &ast.TypeRef{
			Kind:     ast.KindObject,
			Name:     responseName,
			TypeNode: store.Objects[responseName],
		},
	}
}

// addPaginationArguments adds a pagination arguments to the field
func addPaginationArguments(f *ast.Field, store *ast.NodeStore) {
	if f.Args == nil {
		f.Args = make(map[string]*ast.Argument)
	}
	f.Args["page"] = &ast.Argument{
		Name: "page",
		Type: &ast.TypeRef{
			Kind:     ast.KindScalar,
			Name:     "Int",
			TypeNode: store.Scalars["Int"],
		},
		DefaultValue: int64(1),
	}
	f.Args["size"] = &ast.Argument{
		Name: "size",
		Type: &ast.TypeRef{
			Kind:     ast.KindScalar,
			Name:     "Int",
			TypeNode: store.Scalars["Int"],
		},
		DefaultValue: int64(10),
	}
	f.Args["sort"] = &ast.Argument{
		Name: "sort",
		Type: &ast.TypeRef{
			Kind:     ast.KindScalar,
			Name:     "SortOrder",
			TypeNode: store.Enums["SortOrder"],
		},
		DefaultValue: "ASC",
	}
}

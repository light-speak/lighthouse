package ast

import "fmt"

func (f *Field) ParseDirectives(store *NodeStore) error {
	directives := GetDirective("deprecated", f.Directives)
	if len(directives) > 1 {
		directive := directives[0]
		deprecationReason := directive.GetArg("reason")
		if deprecationReason != nil {
			f.IsDeprecated = true
			if deprecationReason.Value != nil {
				f.DeprecationReason = deprecationReason.Value.(string)
			} else if deprecationReason.DefaultValue != nil {
				f.DeprecationReason = deprecationReason.DefaultValue.(string)
			} else {
				f.DeprecationReason = "field is deprecated"
			}
		}
	}
	directives = GetDirective("paginate", f.Directives)
	if len(directives) > 0 {
		f.addPaginationResponseType(store)
		f.addPaginationArguments(store)
	}
	return nil
}

// addPaginationResponseType adds a pagination response type to the store
func (f *Field) addPaginationResponseType(store *NodeStore) {
	curType := f.Type
	for curType.Kind == KindList || curType.Kind == KindNonNull {
		curType = curType.OfType
	}
	typeName := curType.Name
	responseName := fmt.Sprintf("%sPaginateResponse", typeName)
	if _, ok := store.Objects[responseName]; ok {
		return
	}
	store.AddObject(responseName, &ObjectNode{
		BaseNode: BaseNode{
			Name:        responseName,
			Description: fmt.Sprintf("The %sPaginateResponse type represents a paginated list of %s.", typeName, typeName),
		},
		Fields: map[string]*Field{
			"data": {
				Name: "data",
				Type: &TypeRef{
					Kind:   KindNonNull,
					OfType: f.Type,
				},
			},
			"paginateInfo": {
				Name: "paginateInfo",
				Type: &TypeRef{
					Kind: KindNonNull,
					OfType: &TypeRef{
						Kind:     KindObject,
						Name:     "PaginateInfo",
						TypeNode: store.Nodes["PaginateInfo"],
					},
				},
			},
		},
	})
}

// addPaginationArguments adds a pagination arguments to the field
func (f *Field) addPaginationArguments(store *NodeStore) {
	if f.Args == nil {
		f.Args = make(map[string]*Argument)
	}
	f.Args["page"] = &Argument{
		Name: "page",
		Type: &TypeRef{
			Kind:     KindScalar,
			Name:     "Int",
			TypeNode: store.Scalars["Int"],
		},
		DefaultValue: int64(1),
	}
	f.Args["size"] = &Argument{
		Name: "size",
		Type: &TypeRef{
			Kind:     KindScalar,
			Name:     "Int",
			TypeNode: store.Scalars["Int"],
		},
		DefaultValue: int64(10),
	}
	f.Args["sort"] = &Argument{
		Name: "sort",
		Type: &TypeRef{
			Kind:     KindScalar,
			Name:     "SortOrder",
			TypeNode: store.Enums["SortOrder"],
		},
		DefaultValue: "ASC",
	}
}

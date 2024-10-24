package ast

import (
	"fmt"
)

func (f *Field) ParseFieldDirectives(store *NodeStore) error {
	deprecated := GetDirective("deprecated", f.Directives)
	if len(deprecated) > 0 {
		directive := deprecated[0]
		deprecationReason := directive.GetArg("reason")
		if deprecationReason != nil {
			f.IsDeprecated = true
			reason := "field is deprecated"
			if deprecationReason.Value != nil {
				reason = deprecationReason.Value.(string)
			} else if deprecationReason.DefaultValue != nil {
				reason = deprecationReason.DefaultValue.(string)
			}
			f.DeprecationReason = &reason
		}
	}
	deprecated = GetDirective("paginate", f.Directives)
	if len(deprecated) > 0 {
		f.addPaginationResponseType(store)
		f.addPaginationArguments(store)
	}
	return nil
}

func (o *ObjectNode) ParseObjectDirectives(store *NodeStore) error {
	model := GetDirective("model", o.Directives)
	modelSoftDelete := GetDirective("softDelete", o.Directives)
	if len(model) > 0 || len(modelSoftDelete) > 0 {
		o.IsModel = true
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
	description := fmt.Sprintf("The %sPaginateResponse type represents a paginated list of %s.", typeName, typeName)
	store.AddObject(responseName, &ObjectNode{
		BaseNode: BaseNode{
			Name:        responseName,
			Kind:        KindObject,
			Description: &description,
		},
		Fields: map[string]*Field{
			"data": {
				Name: "data",
				Type: f.Type,
			},
			"paginateInfo": {
				Name: "paginateInfo",
				Type: &TypeRef{
					Kind: KindNonNull,
					OfType: &TypeRef{
						Kind:     KindObject,
						Name:     "PaginateInfo",
						TypeNode: store.Objects["PaginateInfo"],
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

package server

import (
	"user/graph/generate"

	"github.com/light-speak/lighthouse/graphql/resolver"
)

func SetDirective(d *generate.DirectiveRoot) {
	d.All = resolver.All
	d.First = resolver.First
	d.Eq = resolver.Eq
	d.Create = resolver.CreateOrUpdate
	d.Update = resolver.CreateOrUpdate
	d.Scope = resolver.Scope
	d.Page = resolver.Page
	d.Size = resolver.Size
	d.Count = resolver.Count
	d.Sum = resolver.Sum
	d.Resolve = resolver.Resolve
}

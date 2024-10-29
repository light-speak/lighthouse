package ast

import "github.com/light-speak/lighthouse/errors"

var fieldDirectiveMap = make(map[string]func(f *Field, d *Directive, store *NodeStore, parent Node) errors.GraphqlErrorInterface)
var objectDirectiveMap = make(map[string]func(o *ObjectNode, d *Directive, store *NodeStore) errors.GraphqlErrorInterface)

func AddFieldDirective(name string, fn func(f *Field, d *Directive, store *NodeStore, parent Node) errors.GraphqlErrorInterface) {
	fieldDirectiveMap[name] = fn
}

func AddObjectDirective(name string, fn func(o *ObjectNode, d *Directive, store *NodeStore) errors.GraphqlErrorInterface) {
	objectDirectiveMap[name] = fn
}

func (f *Field) ParseFieldDirectives(store *NodeStore, parent Node) errors.GraphqlErrorInterface {
	for _, directive := range f.Directives {
		if fn, ok := fieldDirectiveMap[directive.Name]; ok {
			if err := fn(f, directive, store, parent); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *ObjectNode) ParseObjectDirectives(store *NodeStore) errors.GraphqlErrorInterface {
	for _, directive := range o.Directives {
		if fn, ok := objectDirectiveMap[directive.Name]; ok {
			if err := fn(o, directive, store); err != nil {
				return err
			}
		}
	}
	return nil
}

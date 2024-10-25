package ast

var fieldDirectiveMap = make(map[string]func(f *Field, d *Directive, store *NodeStore) error)
var objectDirectiveMap = make(map[string]func(o *ObjectNode, d *Directive, store *NodeStore) error)

func AddFieldDirective(name string, fn func(f *Field, d *Directive, store *NodeStore) error) {
	fieldDirectiveMap[name] = fn
}

func AddObjectDirective(name string, fn func(o *ObjectNode, d *Directive, store *NodeStore) error) {
	objectDirectiveMap[name] = fn
}

func (f *Field) ParseFieldDirectives(store *NodeStore) error {
	for _, directive := range f.Directives {
		if fn, ok := fieldDirectiveMap[directive.Name]; ok {
			if err := fn(f, directive, store); err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *ObjectNode) ParseObjectDirectives(store *NodeStore) error {
	for _, directive := range o.Directives {
		if fn, ok := objectDirectiveMap[directive.Name]; ok {
			if err := fn(o, directive, store); err != nil {
				return err
			}
		}
	}
	return nil
}

package ast

type Argument struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Directives  []*Directive `json:"-"`
	Type        *TypeRef     `json:"type"`
	// Parse Value according to the type of Type
	DefaultValue any `json:"default_value"`
	Value        any `json:"value"`
}

func (a *Argument) Validate(store *NodeStore) error {
	if err := a.Type.Validate(store); err != nil {
		return err
	}
	if a.Value != nil {
		if err := a.Type.ValidateValue(a.Value); err != nil {
			return err
		}
	}
	if a.DefaultValue != nil {
		if err := a.Type.ValidateValue(a.DefaultValue); err != nil {
			return err
		}
	}
	err := ValidateDirectives(a.Name, a.Directives, store)
	if err != nil {
		return err
	}
	return nil
}

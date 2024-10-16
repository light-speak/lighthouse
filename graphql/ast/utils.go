package ast

func GetDirective(name string, directives []*Directive) []*Directive {
	var result []*Directive
	for _, directive := range directives {
		if directive.Name == name {
			result = append(result, directive)
		}
	}
	return result
}

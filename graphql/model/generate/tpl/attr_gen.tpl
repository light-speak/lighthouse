func init() {
{{- range $typeName, $fields := .Fields }}
{{- range $field := $fields }}
	excute.AddAttrResolver("{{ $field.Name | camelCase | ucFirst }}Attr", func(ctx *context.Context, data *sync.Map, resolve resolve.Resolve) (interface{}, error) {
		r := resolve.(*Resolver)
		res, err := r.{{ $field.Name | ucFirst }}AttrResolver(ctx, data)
		return res, err
	})
{{- end }}
{{- end }}
}
{{- range $typeName, $fields := .Fields }}
{{- range $field := $fields }}
func (r *Resolver) {{ $field.Name | ucFirst }}AttrResolver(ctx *context.Context, data *sync.Map
{{- range $arg := $field.Args -}}
, {{ $arg.Name | lcFirst }} {{ $arg.Type.GetGoType false | prefixModels }}
{{- end -}}
) (
	{{- if eq $field.Type.GetRealType.Kind "SCALAR" -}}
	{{- false | $field.Type.GetGoType -}}
	{{- else -}}
	{{- (false | $field.Type.GetGoType) | prefixModels -}}
	{{- end -}}
	, error) {
	{{ $field.Name | ucFirst | funcStart }}
	panic("not implemented")
	{{ $field.Name | ucFirst | funcEnd }}
}
{{- end }}
{{- end }}

{{ section }}
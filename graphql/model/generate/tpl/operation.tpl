{{- range $key, $field := .Fields }}
func {{ $field.Name | ucFirst }}Resolver(ctx *context.Context
	{{- range $index, $arg := $field.Args -}}, 
	{{- $arg.Name -}}
	{{- " " -}} 
	{{- if eq $arg.Type.GetRealType.Kind "SCALAR" -}}
	{{- false | $arg.Type.GetGoType -}}
	{{- else -}}
	{{- (false | $arg.Type.GetGoType) | prefixModels -}}
	{{- end -}}
	{{- end -}} ) (
	{{- if eq $field.Type.GetRealType.Kind "SCALAR" -}}
	{{- false | $field.Type.GetGoType -}}
	{{- else -}}
	{{- (false | $field.Type.GetGoType) | prefixModels -}}
	{{- end -}}
	, error) {
	{{ $field.Name | ucFirst | funcStart }}
	panic("not implement")
	{{ $field.Name | ucFirst | funcEnd }}
}
{{- end }}
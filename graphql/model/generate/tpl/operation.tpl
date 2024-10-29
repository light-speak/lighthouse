{{- range $key, $field := .Fields }}
func {{ $field.Name | ucFirst }}Resolver(ctx context.Context{{ range $index, $arg := $field.Args }}, {{ $arg.Name }} {{ if eq $arg.Type.GetRealType.Kind "SCALAR" }}{{ else }}models.{{ end }}{{ false | $arg.Type.GetGoType }}) ({{ if eq $field.Type.GetRealType.Kind "SCALAR" }}{{ else }}*models.{{ end }}{{ false | $field.Type.GetGoType }}{{ else }}{{ false | $field.Type.GetGoType }}{{ end }}, error) {
	{{ $field.Name | ucFirst | funcStart }}
	panic("not implement")
	{{ .Name | ucFirst | funcEnd }}
}
{{- end }}

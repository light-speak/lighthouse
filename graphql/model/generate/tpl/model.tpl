{{ range .Nodes }}
{{ . | genModel }}
{{- $name := .Name -}}
{{- range .Interfaces }}
func (*{{ $name | ucFirst }}) Is{{ .Name | ucFirst }}() bool { return true }
{{- range .Fields }}
func (this *{{ $name | ucFirst }}) Get{{ .Name | ucFirst }}() {{ false | .Type.GetGoType }} { return this.{{ .Name | ucFirst }} }
{{- end }}
{{- end }}
func (*{{ $name | ucFirst }}) TableName() string { return "{{ if ne .Table "" }}{{ .Table }}{{ else }}{{ .Name | pluralize | lcFirst }}{{ end }}" }
func (*{{ $name | ucFirst }}) TypeName() string { return "{{ .Name | lcFirst }}" }
var {{ $name | ucFirst }}EnumFields = map[string]func(interface{}) interface{} {
  {{- range .Fields }}
  {{- if .Type.IsEnum }}
  "{{ .Name | lcFirst }}": func(value interface{}) interface{} {
    switch v := value.(type) {
    case int64:
      return {{ .Type.GetRealType.Name }}(v)
    case int8:
      return {{ .Type.GetRealType.Name }}(v)
    default:
      return v
    }
  },
  {{- end }}
  {{- end }}
}
{{ end }}

func Migrate() error {
	return model.GetDB().AutoMigrate(
    {{- range $index, $node := .Nodes }}
    &{{ $node.Name }}{},
    {{- end }}
  )
}
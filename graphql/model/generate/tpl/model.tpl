{{ range .Nodes }}
{{ . | genModel }}
{{- $name := .Name -}}
{{- range .Interfaces }}
func (*{{ $name | ucFirst }}) Is{{ .Name | ucFirst }}() bool { return true }
{{- range .Fields }}
func (this *{{ $name | ucFirst }}) Get{{ .Name | ucFirst }}() {{ false | .Type.GetGoType }} { return this.{{ .Name | ucFirst }} }
{{- end }}
{{- end }}
func (*{{ $name | ucFirst }}) GetProvide() map[string]*ast.Relation { return map[string]*ast.Relation{
  {{- range .Fields }}
  {{- if ne .Name "__typename" }}"{{ .Name }}": {{ if .Relation }}{{ buildRelation . }}{{ else }}{}{{ end }},{{ end -}}
  {{- end -}}
}}
{{ end }}

func Migrate() error {
	return model.GetDB().AutoMigrate(
    {{- range $index, $node := .Nodes }}
    &{{ $node.Name }}{},
    {{- end }}
  )
}
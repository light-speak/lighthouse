{{ range .Nodes }}
{{ . | genModel }}
{{- $name := .GetName -}}
{{- range .ImplementTypes }}
func (*{{ $name | ucFirst }}) Is{{ .GetName | ucFirst }}() bool { return true }
{{- range .Fields }}
func (this *{{ $name | ucFirst }}) Get{{ .GetName | ucFirst }}() {{ .Type.GoType }} { return this.{{ .GetName | ucFirst }} }
{{- end }}
{{- end }}
{{ end }}

func AutoMigrate() {
	model.GetDB().AutoMigrate(
    {{- range $index, $node := .Nodes }}
    &{{ $node.GetName }}{},
    {{- end }}
  )
}
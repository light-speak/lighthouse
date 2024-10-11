{{ range .Nodes }}
{{ . | genModel }}
{{ end }}

func init() {
	model.GetDB().AutoMigrate(
    {{- range $index, $node := .Nodes }}
    &{{ $node.GetName }}{},
    {{- end }}
  )
}
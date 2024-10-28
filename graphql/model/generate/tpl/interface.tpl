{{ range .Nodes }}
type {{ .GetName }} interface {
  Is{{ .GetName }}()
  {{- range .GetFields }}
  Get{{ .Name | ucFirst }}() {{ false | .Type.GetGoType }}
  {{- end }}
}
{{ end }}
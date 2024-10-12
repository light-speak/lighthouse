{{ range .Nodes }}
type {{ .GetName }} interface {
  Is{{ .GetName }}()
  {{- range .GetFields }}
  Get{{ .GetName | ucFirst }}() {{ .Type.GoType }}
  {{- end }}
}
{{ end }}
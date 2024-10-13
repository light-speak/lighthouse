{{ range .Nodes }}
type {{ .GetName }} struce {
  {{- range .GetFields }}
  {{ .GetName | ucFirst }} {{ .Type.GoType }} `json:"{{ .GetName | ucFirst }}"`
  {{- end }}
}
{{ end }}
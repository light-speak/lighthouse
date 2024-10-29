{{ range .Nodes }}
type {{ .GetName }} struct {
  {{- range .Fields }}
  {{ .Name | ucFirst }} {{ false | .Type.GetGoType }} `json:"{{ .Name | ucFirst }}"`
  {{- end }}
}
{{ end }}
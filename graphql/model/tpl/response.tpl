{{ range .Nodes }}
type {{ .GetName }} struct {
{{ .GetFields | fields }}
}
{{ end }}
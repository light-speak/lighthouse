{{ range .Nodes }}
type {{ .GetName }} struct {
  {{ range .GetFields -}}
  {{ .GetName | ucFirst }} {{ .Type.GoType }} 
  {{ end }}
}
{{ end }}
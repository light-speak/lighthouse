{{ range .Nodes }}
{{ . | genModel }}
{{ end }}
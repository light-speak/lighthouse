func init () {
{{- range .Nodes }}
{{- if eq .Name "Query" }}
{{- range .Fields }}
{{- if not (isInternalType .Name) }}
{{- if eq (len .Directives) 0 }}
{{- $args := .Args }}
  excute.AddResolver("{{ .Name | ucFirst }}Resolver", func(ctx *context.Context, args map[string]any) (interface{}, error) {
  {{- range $index, $arg := $args }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].({{ if eq $arg.Type.GetRealType.Kind "SCALAR" }}{{ else }}models.{{ end }}{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a {{ if eq $arg.Type.GetRealType.Kind "SCALAR" }}{{ else }}models.{{ end }}{{ false | $arg.Type.GetGoType }}")
    }
  {{- end }}
    return {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
  })
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

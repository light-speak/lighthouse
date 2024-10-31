func init() {
{{- range .Nodes }}
{{- if eq .Name "Query" }}
{{- range .Fields }}
{{- if not (isInternalType .Name) }}
{{- if eq (len .Directives) 0 }}
{{- $args := .Args }}
  excute.AddResolver("{{ .Name }}", func(ctx *context.Context, args map[string]any) (interface{}, error) {
  {{- range $index, $arg := $args }}

    {{- if eq $arg.Type.GetRealType.Kind "SCALAR" }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].({{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a {{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }

    {{- else if eq $arg.Type.GetRealType.Kind "ENUM" }}
    {{ $arg.Name | lcFirst }}Value, ok := args["{{ $index }}"].(int8)
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a int8, got %T", args["{{ $index }}"])
    }
    {{ $arg.Name | lcFirst }} := models.{{ false | $arg.Type.GetGoType }}({{ $arg.Name | lcFirst }}Value)

    {{- else if eq $arg.Type.GetRealType.Kind "OBJECT" }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].(models.{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }

    {{- else }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].(models.{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- end }}
  {{- end }}
    return {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
  })
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

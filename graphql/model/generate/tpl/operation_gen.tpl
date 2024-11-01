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
    {{ $arg.Name | lcFirst }}, ok := models.{{ false | $arg.Type.GetGoType }}Map[args["{{ $arg.Name | lcFirst }}"].(string)]
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- else if eq $arg.Type.GetRealType.Kind "OBJECT" }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].(models.{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- else if eq $arg.Type.GetRealType.Kind "INPUT_OBJECT" }}
    {{ $arg.Name | lcFirst }}Ptr, err := models.MapTo{{ false | $arg.Type.GetGoType }}(args["{{ $index }}"].(map[string]interface{}))
    if err != nil {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{ $arg.Name | lcFirst }} := *{{ $arg.Name | lcFirst }}Ptr
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

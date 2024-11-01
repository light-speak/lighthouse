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
    pv, e := graphql.Parser.NodeStore.Scalars["{{ $arg.Type.GetRealType.Name }}"].ScalarType.ParseValue(args["{{ $index }}"], nil)
    if e != nil {
      return nil, e
    }
    {{ $arg.Name | lcFirst }}, ok := pv.({{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a {{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }

    {{- else if eq $arg.Type.GetRealType.Kind "ENUM" }}
    {{ $arg.Name | lcFirst }}Value, ok := models.{{ false | $arg.Type.GetGoType }}Map[args["{{ $arg.Name | lcFirst }}"].(string)]
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{ $arg.Name | lcFirst }} := &{{ $arg.Name | lcFirst }}Value
    {{- else if eq $arg.Type.GetRealType.Kind "OBJECT" }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].(models.{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- else if eq $arg.Type.GetRealType.Kind "INPUT_OBJECT" }}
    {{ $arg.Name | lcFirst }}, err := models.MapTo{{ false | $arg.Type.GetGoType }}(args["{{ $index }}"].(map[string]interface{}))
    if err != nil {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- else }}
    {{ $arg.Name | lcFirst }}, ok := args["{{ $index }}"].(models.{{ false | $arg.Type.GetGoType }})
    if !ok {
      return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a models.{{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
    }
    {{- end }}
  {{- end }}
    {{- if .Type.IsList }}
    {{- if .Type.IsObject }}
    list, err := {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if list == nil {
      return nil, err
    }
    res := []map[string]interface{}{}
    for _, item := range list {
      itemMap, err := model.StructToMap(item)
      if err != nil {
        return nil, err
      }
      res = append(res, itemMap)
    }
    return res, nil
    {{- else }}
    res, err := {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if res == nil {
      return nil, err
    }
    return res, nil
    {{- end }}
    {{- else if .Type.IsObject }}
    res, err := {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if res == nil {
      return nil, err
    }
    return model.StructToMap(res)
    {{- else }}
    res, err := {{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    return res, err
    {{- end }}
  })
{{- end }}
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

func init() {
{{- range .Nodes }}
{{- range .Fields }}
{{- if not (isInternalType .Name) }}
{{- if not (isQuickDirectives .Directives) }}
  {{- $args := .Args }}
  excute.AddResolver("{{ .Name }}", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
  {{- range $index, $arg := $args }}
    {{- if eq $arg.Type.GetRealType.Kind "SCALAR" }}
    var {{ $arg.Name | lcFirst }} {{ if $arg.Type.IsNullable }}*{{ end }}{{ false | $arg.Type.GetGoType }}
    if args["{{ $index }}"] != nil {
      p{{ $arg.Name | lcFirst }}, e := graphql.Parser.NodeStore.Scalars["{{ $arg.Type.GetRealType.Name }}"].ScalarType.ParseValue(args["{{ $index }}"], nil)
      if e != nil {
        return nil, e
      }
      var ok bool
      {{- if $arg.Type.IsNullable }}
      tmp, ok := p{{ $arg.Name | lcFirst }}.({{ false | $arg.Type.GetGoType }})
      if !ok {
        return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a {{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
      }
      {{ $arg.Name | lcFirst }} = &tmp
      {{- else }}
      {{ $arg.Name | lcFirst }}, ok = p{{ $arg.Name | lcFirst }}.({{ false | $arg.Type.GetGoType }})
      if !ok {
        return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a {{ false | $arg.Type.GetGoType }}, got %T", args["{{ $index }}"])
      }
      {{- end }}
    {{- else if eq $arg.Type.GetRealType.Kind "ENUM" }}
    var {{ $arg.Name | lcFirst }} {{ if $arg.Type.IsNullable }}*{{ end }}{{ false | $arg.Type.GetGoType | prefixModels }}
    if args["{{ $index }}"] != nil {
      {{- if $arg.Type.IsList }}
      var enumValues []{{ false | $arg.Type.GetGoType | prefixModels }}
      for _, v := range args["{{ $index }}"].([]interface{}) {
        enumValue, ok := models.{{ $arg.Type.GetRealType.Name }}Map[v.(string)]
        if !ok {
          return nil, fmt.Errorf("argument: '{{ $arg.Name }}' contains invalid enum value, got %v", v)
        }
        enumValues = append(enumValues, enumValue)
      }
      {{ $arg.Name | lcFirst }} = enumValues
      {{- else }}
      enumValue, ok := models.{{ $arg.Type.GetRealType.Name }}Map[args["{{ $index }}"].(string)]
      if !ok {
        return nil, fmt.Errorf("argument: '{{ $arg.Name }}' is not a valid models.{{ $arg.Type.GetRealType.Name }}, got %v", args["{{ $index }}"])
      }
      {{ $arg.Name | lcFirst }} = &enumValue
      {{- end }}
    {{- else if eq $arg.Type.GetRealType.Kind "INPUT_OBJECT" }}
    var {{ $arg.Name | lcFirst }} {{ if $arg.Type.IsNullable }}*{{ end }}{{ false | $arg.Type.GetGoType | prefixModels }}
    if args["{{ $index }}"] != nil {
      {{- if $arg.Type.IsList }}
      var inputList {{ false | $arg.Type.GetGoType | prefixModels }}
      for _, v := range args["{{ $index }}"].([]interface{}) {
        input, err := models.MapTo{{ $arg.Type.GetRealType.Name }}(v.(map[string]interface{}))
        if err != nil {
          return nil, fmt.Errorf("argument: '{{ $arg.Name }}' contains invalid input, error: %v", err)
        }
        inputList = append(inputList, input)
      }
      {{ $arg.Name | lcFirst }} = inputList
      {{- else }}
      var err error
      {{- if $arg.Type.IsNullable }}
      tmp, err := models.MapTo{{ $arg.Type.GetRealType.Name }}(args["{{ $index }}"].(map[string]interface{}))
      if err != nil {
        return nil, fmt.Errorf("argument: '{{ $arg.Name }}' can not convert to models.{{ $arg.Type.GetRealType.Name }}, error: %v", err)
      }
      {{ $arg.Name | lcFirst }} = &tmp
      {{- else }}
      {{ $arg.Name | lcFirst }}, err = models.MapTo{{ $arg.Type.GetRealType.Name }}(args["{{ $index }}"].(map[string]interface{}))
      if err != nil {
        return nil, fmt.Errorf("argument: '{{ $arg.Name }}' can not convert to models.{{ $arg.Type.GetRealType.Name }}, error: %v", err)
      }
      {{- end }}
      {{- end }}
    {{- else }}
    {{- end }}
    }
  {{- end }}
    {{- if .Type.IsList }}
    {{- if .Type.IsObject }}
    list, err := r.{{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if err != nil {
      return nil, err
    }
    if list == nil {
      return nil, nil
    }
    res := []*sync.Map{}
    for _, item := range list {
      itemMap, err := model.StructToMap(item)
      if err != nil {
        return nil, err
      }
      res = append(res, itemMap)
    }
    return res, nil
    {{- else }}
    res, err := r.{{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if err != nil {
      return nil, err
    }
    return res, nil
    {{- end }}
    {{- else if .Type.IsObject }}
    res, err := r.{{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    if err != nil {
      return nil, err
    }
    if res == nil {
      return nil, nil
    }
    {{- if .Type.GetRealType.TypeNode.IsModel }}
    return model.StructToMap(res)
    {{- else }}
    return model.TypeToMap(res)
    {{- end }}
    {{- else }}
    res, err := r.{{ .Name | ucFirst }}Resolver(ctx{{ range $index, $arg := $args }}, {{ $arg.Name | lcFirst }}{{ end }})
    return res, err
    {{- end }}
  })
{{- end }}
{{- end }}
{{- end }}
{{- end }}
}

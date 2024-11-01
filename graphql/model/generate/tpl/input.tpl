{{ range .Nodes }}
type {{ .GetName }} struct {
  {{- range .Fields }}
  {{ .Name | ucFirst }} {{ false | .Type.GetGoType }} `json:"{{ .Name | ucFirst }}"`
  {{- end }}
}

func MapTo{{ .GetName }}(data map[string]interface{}) (*{{ .GetName }}, error) {
  result := &{{ .GetName }}{}
  {{- range .Fields }}
  {{ if eq .Type.GetRealType.Kind "SCALAR" }}
  {{ .Name | lcFirst }}, ok := data["{{ .Name }}"].({{ false | .Type.GetGoType }})
  if !ok {
    return nil, fmt.Errorf("invalid value for field '{{ .Name }}', got %T", data["{{ .Name }}"])
  }
  result.{{ .Name | ucFirst }} = {{ .Name | lcFirst }}
  {{- else if eq .Type.GetRealType.Kind "ENUM" }}
  {{ .Name | lcFirst }}, ok := {{ false | .Type.GetGoType }}Map[data["{{ .Name }}"].(string)]
  if !ok {
    return nil, fmt.Errorf("invalid value for field '{{ .Name }}', got %T", data["{{ .Name }}"])
  }
  result.{{ .Name | ucFirst }} = {{ .Name | lcFirst }}
  {{- else if eq .Type.GetRealType.Kind "OBJECT_INPUT" }}
  inputPtr, err := MapTo{{ .Type.GetGoType }}(data["{{ .Name }}"].(map[string]interface{}))
  if err != nil {
    return nil, fmt.Errorf("invalid value for field '{{ .Name }}', got %T", data["{{ .Name }}"])
  }
  result.{{ .Name | ucFirst }} = *inputPtr
  {{- else }}
  {{ .Name | lcFirst }}, ok := data["{{ .Name }}"].({{ false | .Type.GetGoType }})
  if !ok {
    return nil, fmt.Errorf("invalid value for field '{{ .Name }}', got %T", data["{{ .Name }}"])
  }
  result.{{ .Name | ucFirst }} = {{ .Name | lcFirst }}
  {{- end }}
  {{- end }}
  return result, nil
}
{{- end }}

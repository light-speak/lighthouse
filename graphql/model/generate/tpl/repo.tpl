{{ range .Nodes }}
{{- $name := .Name -}}
func Provide__{{ $name | ucFirst }}() map[string]*ast.Relation { return map[string]*ast.Relation{
  {{- range .Fields }}
  {{- if ne .Name "__typename" }}"{{ .Name }}": {{ if .Relation }}{{ buildRelation . }}{{ else }}{}{{ end }},{{ end -}}
  {{- end -}}
}}
func Query__{{ $name | ucFirst }}(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.{{ $name | ucFirst }}{}).Scopes(scopes...)
}
func Fields__{{ $name | ucFirst }}({{ $name | lcFirst }} *models.{{ $name | ucFirst }}, key string) (interface{}, error) {
  switch key {
  {{- range .Fields }}
    {{- if ne .Name "__typename" }}
    case "{{ .Name }}": 
      return {{ $name | lcFirst }}.{{ .Name | camelCase | ucFirst }}, nil
    {{- end -}}
  {{- end }}
  }
  return nil, nil
} 
func First__{{ $name | ucFirst }}(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  selectColumns, selectRelations := model.GetSelectInfo(columns, Provide__{{ $name | ucFirst }}())
  {{ $name | lcFirst }} := &models.{{ $name | ucFirst }}{}
  err := Query__{{ $name | ucFirst }}().Scopes(scopes...).Select(selectColumns).First({{ $name | lcFirst }}).Error
  if err != nil {
    return nil, err
  }
  res, err := model.StructToMap({{ $name | lcFirst }})
  if err != nil {
    return nil, err
  }
  for _, relation := range selectRelations {
    fieldValue, err := Fields__{{ $name | ucFirst }}({{ $name | lcFirst }}, relation.Relation.ForeignKey)
    if err != nil {
      return nil, err
    }
    res, err = model.FetchRelation(res, relation, fieldValue)
    if err != nil {
      return nil, err
    }
  }
  return res, nil
}
{{ end }}

func init() {
  {{- range .Nodes }}
  model.AddQuickFirst("{{ .Name | ucFirst }}", First__{{ .Name | ucFirst }})
  {{- end }}
}

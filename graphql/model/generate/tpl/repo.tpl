{{ range .Nodes }}
{{- $name := .Name -}}
func Provide__{{ $name | ucFirst }}() map[string]*ast.Relation { return map[string]*ast.Relation{
  {{- range .Fields }}
  {{- if ne .Name "__typename" }}"{{ .Name }}": {{ if .Relation }}{{ buildRelation . }}{{ else }}{}{{ end }},{{ end -}}
  {{- end -}}
}}
func Load__{{ $name | ucFirst }}(ctx *context.Context, key int64, field string) (map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "{{ if ne .Table "" }}{{ .Table }}{{ else }}{{ $name | pluralize | lcFirst }}{{ end }}", field).Load(key)
}
func LoadList__{{ $name | ucFirst }}(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "{{ if ne .Table "" }}{{ .Table }}{{ else }}{{ $name | pluralize | lcFirst }}{{ end }}", field).LoadList(key)
}
func Query__{{ $name | ucFirst }}(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.{{ $name | ucFirst }}{}).Scopes(scopes...)
}
func First__{{ $name | ucFirst }}(ctx *context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  var err error
  selectColumns := model.GetSelectInfo(columns, Provide__{{ $name | ucFirst }}())
  if data == nil {
    data = make(map[string]interface{})
    err = Query__{{ $name | ucFirst }}().Scopes(scopes...).Select(selectColumns).First(data).Error
    if err != nil {
      return nil, err
    }
  }
  return data, nil
}
func List__{{ $name | ucFirst }}(ctx *context.Context, columns map[string]interface{},datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
  var err error
  selectColumns := model.GetSelectInfo(columns, Provide__{{ $name | ucFirst }}())
  if datas == nil {
    datas = make([]map[string]interface{}, 0)
    err = Query__{{ $name | ucFirst }}().Scopes(scopes...).Select(selectColumns).Find(&datas).Error
    if err != nil {
      return nil, err
    }
  }
  return datas, nil
}
func Count__{{ $name | ucFirst }}(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  var count int64
  err := Query__{{ $name | ucFirst }}().Scopes(scopes...).Count(&count).Error
  return count, err
}
{{ end }}

func init() {
  {{- range .Nodes }}
  model.AddQuickFirst("{{ .Name | ucFirst }}", First__{{ .Name | ucFirst }})
  model.AddQuickList("{{ .Name | ucFirst }}", List__{{ .Name | ucFirst }})
  model.AddQuickLoad("{{ .Name | ucFirst }}", Load__{{ .Name | ucFirst }})
  model.AddQuickLoadList("{{ .Name | ucFirst }}", LoadList__{{ .Name | ucFirst }})
  model.AddQuickCount("{{ .Name | ucFirst }}", Count__{{ .Name | ucFirst }})
  {{- end }}
}

{{ range $node, $fields := .Fields }}
func GetSearchableModelMapping() map[string]model.SearchModel {
	return map[string]model.SearchModel{
		"{{ $node.GetName }}": &{{ $node.GetName }}{},
	}
}
{{ end }}

{{ range $node, $fields := .Fields }}
{{- $t := $node.GetName -}}
func ({{$t | lcFirst}} *{{ $t }}) FieldMapping() map[string]interface{} {
	return map[string]interface{}{
		{{- range $field := $fields }}
		"{{ $field.Field.Name | lcFirst }}": map[string]string{
			"type": "{{ $field.Type | lc }}",
			"analyzer": "{{ $field.IndexAnalyzer | lc }}",
			"search_analyzer": "{{ $field.SearchAnalyzer | lc }}",
		},
		{{- end }}
	}
}

func ({{$t | lcFirst}} *{{ $t }}) SearchId() int64 {
	return {{$t | lcFirst}}.Id
}

func ({{$t | lcFirst}} *{{ $t }}) IndexName() string {
	return {{ $t | lcFirst }}.TableName()
}

func ({{$t | lcFirst}} *{{ $t }}) GetSearchData(mapData ...map[string]interface{}) map[string]interface{} {
	if len(mapData) > 0 {
		obj := mapData[0]
		data := map[string]interface{}{
			{{- range $field := $fields }}
			"{{ $field.Field.Name | lcFirst }}": obj["{{ $field.Field.Name | snakeCase }}"],
			{{- end }}
		}
		return data
	}else{
		data := map[string]interface{}{
			{{- range $field := $fields }}
			"{{ $field.Field.Name | lcFirst }}": {{$t | lcFirst}}.{{ $field.Field.Name | camelCase | ucFirst}},
			{{- end }}
		}
		return data
	}
}

{{ end }}

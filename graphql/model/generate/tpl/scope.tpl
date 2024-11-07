
{{- range .Scopes }}
func {{ $.Name | camelCase | ucFirst }}{{ . | camelCase | ucFirst }}(ctx *context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
    {{ ( . | camelCase | ucFirst) | funcStart}}
		return db
    {{ ( . | camelCase | ucFirst) | funcEnd}}
	}
}
{{- end }}


func init() {
{{- range .Scopes }}
	model.AddScopes("{{ $.Name | camelCase | ucFirst }}{{ . | camelCase | ucFirst }}", {{ $.Name | camelCase | ucFirst }}{{ . | camelCase | ucFirst }})
{{- end }}
}

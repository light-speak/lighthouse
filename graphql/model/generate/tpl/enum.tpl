{{- range $key, $node := .Nodes }}
{{- if not (isInternalType $node.Name) }}
type {{ $key }} int8

const (
  {{- $index := 0}}
  {{- range $vKey, $enumValue := $node.EnumValues }}
  {{- if $enumValue.Value }}
  {{ $vKey }} {{ $key }} = {{ $enumValue.Value }}
  {{- else }}
  {{ $vKey }}{{ if eq $index 0 }} {{ $key }} = iota{{ end }}
  {{- $index = add $index 1 }}
  {{- end }}
  {{- end }}
)

var {{ $key }}Map = map[string]{{ $key }}{
  {{- range $vKey, $enumValue := $node.EnumValues }}
  "{{ $vKey }}": {{ $vKey }},
  {{- end }}
}
{{- end }}
{{- end }}

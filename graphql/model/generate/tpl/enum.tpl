{{- range $key, $node := .Nodes }}
{{- if not (isInternalType $node.Name) }}
type {{ $key }} int8

const (
  {{- $index := 0}}
  {{- range $vKey, $enumValue := $node.EnumValues }}
  // {{ $enumValue.Description }}
  {{- if $enumValue.Value }}
  {{ $key }}{{ $vKey }} {{ $key }} = {{ $enumValue.Value }}
  {{- else }}
  {{ $key }}{{ $vKey }}{{ if eq $index 0 }} = iota{{ end }}
  {{- $index = add $index 1 }}
  {{- end }}
  {{- end }}
)

func (e {{ $key }}) ToString() string {
  switch e {
  {{- range $vKey, $enumValue := $node.EnumValues }}
  case {{ $key }}{{ $vKey }}:
    return "{{ $vKey }}"
  {{- end }}
  default:
    return "unknown"
  }
}

var {{ $key }}Map = map[string]{{ $key }}{
  {{- range $vKey, $enumValue := $node.EnumValues }}
  "{{ $vKey }}": {{ $key }}{{ $vKey }},
  {{- end }}
}
{{- end }}
{{- end }}

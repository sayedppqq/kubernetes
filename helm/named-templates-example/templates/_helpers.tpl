{{/* Generate basic labels */}}
{{- define "mychart.labels" }}
labels:
  generator: helm
  date: {{ now | htmlDate }}
  chart: {{ .Chart.Name }}
  version: {{ .Chart.Version }}
{{- end }}

{{/* Testing */}}
{{- define "mychart.test" }}
  labels:
    generator: helmmmmmmmmmmm
    date: {{ now | htmlDate }}
{{- end }}
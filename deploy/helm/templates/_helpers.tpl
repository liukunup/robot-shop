{{/*
Expand the name of the chart.
*/}}
{{- define "robotshop.name" -}}
{{- default .Chart.Name .Values.global.projectName | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "robotshop.fullname" -}}
{{- $name := default .Chart.Name .Values.global.projectName -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "robotshop.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "robotshop.labels" -}}
helm.sh/chart: {{ include "robotshop.chart" . }}
{{ include "robotshop.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "robotshop.selectorLabels" -}}
app.kubernetes.io/name: {{ include "robotshop.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "robotshop.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "robotshop.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end -}}
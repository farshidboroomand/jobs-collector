{{/*
Expand the name of the chart.
*/}}
{{- define "jobs-collector.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "jobs-collector.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "jobs-collector.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "jobs-collector.labels" -}}
helm.sh/chart: {{ include "jobs-collector.chart" . }}
{{ include "jobs-collector.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels.
*/}}
{{- define "jobs-collector.selectorLabels" -}}
app.kubernetes.io/name: {{ include "jobs-collector.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Service account name.
*/}}
{{- define "jobs-collector.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "jobs-collector.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Namespace for rendered resources.
*/}}
{{- define "jobs-collector.namespace" -}}
{{- default .Release.Namespace .Values.namespace.name }}
{{- end }}

{{/*
Application image.
*/}}
{{- define "jobs-collector.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag -}}
{{- printf "%s:%s" .Values.image.repository $tag }}
{{- end }}

{{/*
MySQL image.
*/}}
{{- define "jobs-collector.mysqlImage" -}}
{{- printf "%s:%s" .Values.mysql.image.repository .Values.mysql.image.tag }}
{{- end }}

{{/*
MySQL service name.
*/}}
{{- define "jobs-collector.mysql.fullname" -}}
{{- printf "%s-mysql" (include "jobs-collector.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Authentication secret name.
*/}}
{{- define "jobs-collector.secretName" -}}
{{- if .Values.mysql.auth.existingSecret }}
{{- .Values.mysql.auth.existingSecret }}
{{- else }}
{{- printf "%s-auth" (include "jobs-collector.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

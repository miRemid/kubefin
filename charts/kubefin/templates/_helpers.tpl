{{/*
Expand the name of the chart.
*/}}
{{- define "kubefin.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubefin.fullname" -}}
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
{{- define "kubefin.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Mimir labels
*/}}
{{- define "kubefin-mimir.labels" -}}
helm.sh/chart: {{ include "kubefin.chart" . }}
{{ include "kubefin-mimir.selectorLabels" . }}
app.kubernetes.io/version: {{ .Values.mimir.image.tag | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Mimir selector labels
*/}}
{{- define "kubefin-mimir.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubefin.name" . }}-mimir
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "kubefin-cost-analyzer.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "kubefin-dashboard.name" -}}
{{- printf "%s" "kubefin-dashboard"  }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "kubefin-cost-analyzer.fullname" -}}
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

{{- define "kubefin-dashboard.fullname" -}}
{{- if eq .Release.Name .Chart.Name }}
{{- include "kubefin-dashboard.name" . }}
{{- else }}
{{- printf "%s-%s" .Release.Name (include "kubefin-dashboard.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "kubefin-cost-analyzer.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Kubefin-cost-analyzer labels
*/}}
{{- define "kubefin-cost-analyzer.labels" -}}
helm.sh/chart: {{ include "kubefin-cost-analyzer.chart" . }}
{{ include "kubefin-cost-analyzer.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Kubefin-dashboard labels
*/}}
{{- define "kubefin-dashboard.labels" -}}
helm.sh/chart: {{ include "kubefin-cost-analyzer.chart" . }}
{{ include "kubefin-dashboard.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Kubefin-cost-analyzer selector labels
*/}}
{{- define "kubefin-cost-analyzer.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubefin-cost-analyzer.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Kubefin-dashboard selector labels
*/}}
{{- define "kubefin-dashboard.selectorLabels" -}}
app.kubernetes.io/name: {{ include "kubefin-dashboard.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

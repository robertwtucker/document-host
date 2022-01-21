{{/*
Expand the name of the docuhost
*/}}
{{- define "docuhost.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "docuhost.fullname" -}}
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
{{- define "docuhost.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "docuhost.labels" -}}
helm.sh/chart: {{ include "docuhost.chart" . }}
{{ include "docuhost.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "docuhost.selectorLabels" -}}
app.kubernetes.io/name: {{ include "docuhost.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "docuhost.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "docuhost.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Return true if a secret object should be created
*/}}
{{- define "docuhost.createSecret" }}
{{- if not (include "docuhost.useExistingSecret" .) }}
    {{- true }}
{{- end }}
{{- end }}

{{/*
Get the password secret.
*/}}
{{- define "docuhost.secretName" }}
{{- if .Values.existingSecret }}
    {{- printf "%s" (tpl .Values.existingSecret $) }}
{{- else }}
    {{- printf "%s" (include "docuhost.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Return true if we should use an existingSecret.
*/}}
{{- define "docuhost.useExistingSecret" }}
{{- if .Values.existingSecret }}
    {{- true }}
{{- end }}
{{- end }}

{{/*
Return the configuration configmap name
*/}}
{{- define "docuhost.configMapName" }}
{{- if .Values.existingConfigMap }}
    {{- printf "%s" (tpl .Values.existingConfigMap $) }}
{{- else }}
    {{- printf "%s-configuration" (include "docuhost.fullname" .)  }}
{{- end }}
{{- end }}

{{/*
Return true if a configmap object should be created
*/}}
{{- define "docuhost.createConfigMap" }}
{{- if empty .Values.existingConfigmap }}
    {{- true }}
{{- end }}
{{- end }}

{{/*
Return the DocuHost database environment settings
*/}}
{{- define "docuhost.env.database" -}}
- name: DB_PREFIX
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-prefix
- name: DB_USER
  valueFrom:
    secretKeyRef:
      name: {{ include "docuhost.secretName" . }}
      key: db-username
- name: DB_PASS
  valueFrom:
    secretKeyRef:
      name: {{ include "docuhost.secretName" . }}
      key: db-password
- name: DB_HOST
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-host
- name: DB_PORT
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-port
- name: DB_NAME
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-name
- name: DB_COLLECTION
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-collection
- name: DB_TIMEOUT
  valueFrom:
    configMapKeyRef:
      name: {{ include "docuhost.configMapName" . }}
      key: db-timeout
{{- end }}

{{/*
Create labels for the app.
*/}}
{{- define "smart-home.appLabels" -}}
app: {{ .name  }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "smart-home.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "smart-home.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

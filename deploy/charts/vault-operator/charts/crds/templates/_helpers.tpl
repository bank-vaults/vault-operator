{{/* Shortened name suffixed with upgrade-crd */}}
{{- define "vault-operator.crd.upgradeJob.name" -}}
{{- print (include "vault-operator.fullname" .) "-upgrade" -}}
{{- end -}}

{{- define "vault-operator.crd.upgradeJob.labels" -}}
app.kubernetes.io/name: {{ include "vault-operator.name" . }}
app.kubernetes.io/component: crds-upgrade
{{- end -}}

{{/* Create the name of crd.upgradeJob service account to use */}}
{{- define "vault-operator.crd.upgradeJob.serviceAccountName" -}}
{{- if .Values.upgradeJob.serviceAccount.create -}}
    {{ default (include "vault-operator.crd.upgradeJob.name" .) .Values.upgradeJob.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.upgradeJob.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Expand the name of the chart.
*/}}
{{- define "deshchain-testnet.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "deshchain-testnet.fullname" -}}
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
{{- define "deshchain-testnet.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "deshchain-testnet.labels" -}}
helm.sh/chart: {{ include "deshchain-testnet.chart" . }}
{{ include "deshchain-testnet.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "deshchain-testnet.selectorLabels" -}}
app.kubernetes.io/name: {{ include "deshchain-testnet.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "deshchain-testnet.serviceAccountName" -}}
{{- if .Values.security.serviceAccount.create }}
{{- default (include "deshchain-testnet.fullname" .) .Values.security.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.security.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the image registry
*/}}
{{- define "deshchain-testnet.imageRegistry" -}}
{{- if .Values.global.imageRegistry -}}
{{- .Values.global.imageRegistry -}}
{{- else -}}
{{- .Values.deshchain.image.registry -}}
{{- end -}}
{{- end -}}

{{/*
Get the image repository
*/}}
{{- define "deshchain-testnet.imageRepository" -}}
{{- if .Values.deshchain.image.repository -}}
{{- .Values.deshchain.image.repository -}}
{{- else -}}
{{- .Values.global.imageRepository | default "deshchain/node" -}}
{{- end -}}
{{- end -}}

{{/*
Get the image tag
*/}}
{{- define "deshchain-testnet.imageTag" -}}
{{- if .Values.deshchain.image.tag -}}
{{- .Values.deshchain.image.tag -}}
{{- else -}}
{{- .Chart.AppVersion -}}
{{- end -}}
{{- end -}}

{{/*
Get the image pull policy
*/}}
{{- define "deshchain-testnet.imagePullPolicy" -}}
{{- if .Values.deshchain.image.pullPolicy -}}
{{- .Values.deshchain.image.pullPolicy -}}
{{- else -}}
{{- "IfNotPresent" -}}
{{- end -}}
{{- end -}}

{{/*
Get the storage class
*/}}
{{- define "deshchain-testnet.storageClass" -}}
{{- if .Values.global.storageClass -}}
{{- .Values.global.storageClass -}}
{{- else -}}
{{- "standard" -}}
{{- end -}}
{{- end -}}

{{/*
Create image pull secrets
*/}}
{{- define "deshchain-testnet.imagePullSecrets" -}}
{{- if .Values.global.imagePullSecrets }}
imagePullSecrets:
{{- range .Values.global.imagePullSecrets }}
  - name: {{ . }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Validator labels
*/}}
{{- define "deshchain-testnet.validatorLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: validator
{{- end }}

{{/*
Validator selector labels
*/}}
{{- define "deshchain-testnet.validatorSelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: validator
{{- end }}

{{/*
Sentry labels
*/}}
{{- define "deshchain-testnet.sentryLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: sentry
{{- end }}

{{/*
Sentry selector labels
*/}}
{{- define "deshchain-testnet.sentrySelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: sentry
{{- end }}

{{/*
Frontend labels
*/}}
{{- define "deshchain-testnet.frontendLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "deshchain-testnet.frontendSelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Explorer labels
*/}}
{{- define "deshchain-testnet.explorerLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: explorer
{{- end }}

{{/*
Explorer selector labels
*/}}
{{- define "deshchain-testnet.explorerSelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: explorer
{{- end }}

{{/*
Faucet labels
*/}}
{{- define "deshchain-testnet.faucetLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: faucet
{{- end }}

{{/*
Faucet selector labels
*/}}
{{- define "deshchain-testnet.faucetSelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: faucet
{{- end }}

{{/*
IPFS labels
*/}}
{{- define "deshchain-testnet.ipfsLabels" -}}
{{ include "deshchain-testnet.labels" . }}
app.kubernetes.io/component: ipfs
{{- end }}

{{/*
IPFS selector labels
*/}}
{{- define "deshchain-testnet.ipfsSelectorLabels" -}}
{{ include "deshchain-testnet.selectorLabels" . }}
app.kubernetes.io/component: ipfs
{{- end }}

{{/*
Generate certificates secret name
*/}}
{{- define "deshchain-testnet.certificatesSecretName" -}}
{{- printf "%s-certificates" (include "deshchain-testnet.fullname" .) -}}
{{- end -}}

{{/*
Generate validator keys secret name
*/}}
{{- define "deshchain-testnet.validatorKeysSecretName" -}}
{{- printf "%s-validator-keys" (include "deshchain-testnet.fullname" .) -}}
{{- end -}}

{{/*
Generate faucet keys secret name
*/}}
{{- define "deshchain-testnet.faucetKeysSecretName" -}}
{{- printf "%s-faucet-keys" (include "deshchain-testnet.fullname" .) -}}
{{- end -}}

{{/*
Generate genesis config name
*/}}
{{- define "deshchain-testnet.genesisConfigName" -}}
{{- printf "%s-genesis" (include "deshchain-testnet.fullname" .) -}}
{{- end -}}

{{/*
Generate node config name
*/}}
{{- define "deshchain-testnet.nodeConfigName" -}}
{{- printf "%s-config" (include "deshchain-testnet.fullname" .) -}}
{{- end -}}

{{/*
Common environment variables for DeshChain nodes
*/}}
{{- define "deshchain-testnet.commonEnv" -}}
- name: CHAIN_ID
  value: {{ .Values.deshchain.chain.id | quote }}
- name: MONIKER
  valueFrom:
    fieldRef:
      fieldPath: metadata.name
{{- end }}

{{/*
Persistent peers configuration
*/}}
{{- define "deshchain-testnet.persistentPeers" -}}
{{- $peers := list -}}
{{- range $i := until (int .Values.validators.replicas) -}}
{{- $peers = append $peers (printf "%s-validator-%d.%s-validator:26656" (include "deshchain-testnet.fullname" $) $i (include "deshchain-testnet.fullname" $)) -}}
{{- end -}}
{{- join "," $peers -}}
{{- end -}}
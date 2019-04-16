{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "nginx-ingress.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create labels
*/}}
{{- define "nginx-ingress.labels" -}}
app.kubernetes.io/name: {{ include "nginx-ingress.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "kubernetes-ingress.wallarmTarantoolPort" -}}3313{{- end -}}
{{- define "kubernetes-ingress.wallarmTarantoolName" -}}{{ .Values.controller.name }}-wallarm-tarantool{{- end -}}
{{- define "kubernetes-ingress.wallarmSecret" -}}{{ .Values.controller.name }}-secret{{- end -}}

{{- define "kubernetes-ingress.wallarmInitContainer" -}}
- name: addnode
  image: "{{ .Values.controller.image.repository }}:{{ .Values.controller.image.tag }}"
  imagePullPolicy: "{{ .Values.controller.image.pullPolicy }}"
  command:
  - sh
  - -c
  - /usr/share/wallarm-common/synccloud --one-time && chmod 0644 /etc/wallarm/*
  env:
  - name: WALLARM_API_HOST
    value: {{ .Values.controller.wallarm.apiHost | default "api.wallarm.com" }}
  - name: WALLARM_API_TOKEN
    valueFrom:
      secretKeyRef:
        key: token
        name: {{ template "kubernetes-ingress.wallarmSecret" . }}
  - name: WALLARM_SYNCNODE_OWNER
    value: nginx
  - name: WALLARM_SYNCNODE_GROUP
    value: nginx
  volumeMounts:
  - mountPath: /etc/wallarm
    name: wallarm
  securityContext:
    runAsUser: 0
{{- end -}}

{{- define "kubernetes-ingress.wallarmSyncnodeContainer" -}}
- name: synccloud
  image: "{{ .Values.controller.image.repository }}:{{ .Values.controller.image.tag }}"
  imagePullPolicy: "{{ .Values.controller.image.pullPolicy }}"
  command:
  - sh
  - -c
  - /usr/share/wallarm-common/synccloud
  env:
  - name: WALLARM_API_HOST
    value: {{ .Values.controller.wallarm.apiHost | default "api.wallarm.com" }}
  - name: WALLARM_API_TOKEN
    valueFrom:
      secretKeyRef:
        key: token
        name: {{ template "kubernetes-ingress.wallarmSecret" . }}
  - name: WALLARM_SYNCNODE_OWNER
    value: nginx
  - name: WALLARM_SYNCNODE_GROUP
    value: nginx
  volumeMounts:
  - mountPath: /etc/wallarm
    name: wallarm
  securityContext:
    runAsUser: 0
{{- end -}}

{{- define "kubernetes-ingress.wallarmCollectdContainer" -}}
- name: collectd
  image: "{{ .Values.controller.image.repository }}:{{ .Values.controller.image.tag }}"
  imagePullPolicy: "{{ .Values.controller.image.pullPolicy }}"
  command: ["/usr/sbin/collectd", "-f"]
  volumeMounts:
    - name: wallarm
      mountPath: /etc/wallarm
    - name: collectd-config
      mountPath: /etc/collectd
{{- end -}}

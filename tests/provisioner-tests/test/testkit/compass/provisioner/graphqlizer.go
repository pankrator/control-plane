package provisioner

import (
	"bytes"
	"encoding/json"
	"reflect"
	"text/template"

	"github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/pkg/errors"
)

// Graphqlizer is responsible for converting Go objects to input arguments in graphql format
type graphqlizer struct{}

func (g *graphqlizer) ProvisionRuntimeInputToGraphQL(in gqlschema.ProvisionRuntimeInput) (string, error) {
	return g.genericToGraphQL(in, `{
		{{- if .RuntimeInput }}
		runtimeInput: {{ RuntimeInputToGraphQL .RuntimeInput }}
		{{- end }}
		{{- if .ClusterConfig }}
		clusterConfig: {{ ClusterConfigToGraphQL .ClusterConfig }}
		{{- end }}
		{{- if .KymaConfig }}
		kymaConfig: {{ KymaConfigToGraphQL .KymaConfig }}
		{{- end }}
	}`)
}

func (g *graphqlizer) RuntimeInputToGraphQL(in gqlschema.RuntimeInput) (string, error) {
	return g.genericToGraphQL(in, `{
		name: "{{ .Name }}"	
		{{- if .Description }}
		description: {{ .Description }}
		{{- end }}
		{{- if .Labels }}
		labels: {{ .Labels }}
		{{- end }}
	}`)
}

func (g *graphqlizer) UpgradeShootInputToGraphQL(in gqlschema.UpgradeShootInput) (string, error) {
	return g.genericToGraphQL(in, `{
		gardenerConfig: {{ GardenerUpgradeInputToGraphQL .GardenerConfig }}
	}`)
}

func (g *graphqlizer) UpgradeRuntimeInputToGraphQL(in gqlschema.UpgradeRuntimeInput) (string, error) {
	return g.genericToGraphQL(in, `{
		kymaConfig: {{ KymaConfigToGraphQL .KymaConfig }}
	}`)
}

func (g *graphqlizer) ClusterConfigToGraphQL(in gqlschema.ClusterConfigInput) (string, error) {
	return g.genericToGraphQL(in, `{
		{{- if .GardenerConfig }}
		gardenerConfig: {{ GardenerConfigInputToGraphQL .GardenerConfig }}
		{{- end }}
	}`)
}

func (g *graphqlizer) GardenerConfigInputToGraphQL(in gqlschema.GardenerConfigInput) (string, error) {

	return g.genericToGraphQL(in, `{
		kubernetesVersion: "{{ .KubernetesVersion }}"
		volumeSizeGB: {{ .VolumeSizeGb }}
		machineType: "{{ .MachineType }}"
		region: "{{ .Region }}"
		provider: "{{ .Provider }}"
		{{- if .Purpose }}
		purpose: "{{ .Purpose }}"
		{{- end }}
		diskType: "{{ .DiskType }}"
		{{- if .Seed }}
		seed: "{{ .Seed }}"
		{{- end }}
		targetSecret: "{{ .TargetSecret }}"
		workerCidr: "{{ .WorkerCidr }}"
        autoScalerMin: {{ .AutoScalerMin }}
        autoScalerMax: {{ .AutoScalerMax }}
        maxSurge: {{ .MaxSurge }}
		maxUnavailable: {{ .MaxUnavailable }}
		{{- if .EnableKubernetesVersionAutoUpdate }}
		enableKubernetesVersionAutoUpdate: {{ .EnableKubernetesVersionAutoUpdate }}
		{{- end }}
		{{- if .EnableMachineImageVersionAutoUpdate }}
		enableMachineImageVersionAutoUpdate: {{ .EnableMachineImageVersionAutoUpdate }}
		{{- end }}
		{{- if .AllowPrivilegedContainers }}
		allowPrivilegedContainers: {{ .AllowPrivilegedContainers }}
		{{- end }}
		providerSpecificConfig: {{ ProviderSpecificInputToGraphQL .ProviderSpecificConfig }}
	}`)
}

func (g *graphqlizer) GardenerUpgradeInputToGraphQL(in gqlschema.GardenerUpgradeInput) (string, error) {

	return g.genericToGraphQL(in, `{
		{{- if .KubernetesVersion }}
		kubernetesVersion: "{{ .KubernetesVersion }}"
        {{- end }}
		{{- if .VolumeSizeGb }}
		volumeSizeGB: {{ .VolumeSizeGb }}
        {{- end }}
        {{- if .MachineType }}
		machineType: "{{ .MachineType }}"
		{{- end }}
		{{- if .Purpose }}
		purpose: "{{ .Purpose }}"
		{{- end }}
		{{- if .DiskType }}
		diskType: "{{ .DiskType }}"
		{{- end }}
		{{- if .AutoScalerMin }}	
        autoScalerMin: {{ .AutoScalerMin }}
		{{- end }}
		{{- if .AutoScalerMax }}
        autoScalerMax: {{ .AutoScalerMax }}
		{{- end }}
		{{- if .MaxSurge }}
        maxSurge: {{ .MaxSurge }}
		{{- end }}
		{{- if .MaxUnavailable }}
		maxUnavailable: {{ .MaxUnavailable }}
		{{- end }}
		{{- if .EnableKubernetesVersionAutoUpdate }}
		enableKubernetesVersionAutoUpdate: {{ .EnableKubernetesVersionAutoUpdate }}
		{{- end }}
		{{- if .EnableMachineImageVersionAutoUpdate }}
		enableMachineImageVersionAutoUpdate: {{ .EnableMachineImageVersionAutoUpdate }}
		{{- end }}
		{{- if .ProviderSpecificConfig }}
		providerSpecificConfig: {{ ProviderSpecificInputToGraphQL .ProviderSpecificConfig }}
        {{- end }}
	}`)
}

func (g *graphqlizer) ProviderSpecificInputToGraphQL(in *gqlschema.ProviderSpecificInput) (string, error) {
	return g.genericToGraphQL(in, `{
		{{- if .AzureConfig }}
		azureConfig: {{ AzureProviderConfigInputToGraphQL .AzureConfig }}
		{{- end }}
		{{- if .GcpConfig }}
		gcpConfig: {{ GcpProviderConfigInputToGraphQL .GcpConfig }}
		{{- end }}
	}`)
}

func (g *graphqlizer) AzureProviderConfigInputToGraphQL(in *gqlschema.AzureProviderConfigInput) (string, error) {
	return g.genericToGraphQL(in, `{
		vnetCidr: "{{.VnetCidr}}",
		{{- if .Zones }}
		zones: {{.Zones | marshal }},
		{{- end }}
	}`)
}

func (g *graphqlizer) GcpProviderConfigInputToGraphQL(in *gqlschema.GCPProviderConfigInput) (string, error) {
	return g.genericToGraphQL(in, `{
		zones: "{{ .Zones | marshal }}"
	}`)
}

func (g *graphqlizer) UpgradeClusterConfigToGraphQL(in gqlschema.UpgradeRuntimeInput) (string, error) {
	return g.genericToGraphQL(in, `{
		{{- if .Version }}
		version: "{{.Version}}"
		{{- end }}
	}`)
}

func (g *graphqlizer) KymaConfigToGraphQL(in gqlschema.KymaConfigInput) (string, error) {
	return g.genericToGraphQL(in, `{
		version: "{{.Version}}"
		{{- if .Components }}
		components: [
			{{- range $i, $e := .Components }}
			{{- if $i}}, {{- end}} {{ ComponentConfigurationInputToGQL $e }}
			{{- end }}]
		{{- end }}
		{{- if .Configuration }}
		configuration: [
			{{- range $i, $e := .Configuration }}
			{{- if $i}}, {{- end}} {{ ConfigEntryInputToGQL $e }}
			{{- end }}]
		{{- end }}
	}`)
}

func (g *graphqlizer) ComponentConfigurationInputToGQL(in gqlschema.ComponentConfigurationInput) (string, error) {
	return g.genericToGraphQL(in, `{
		component: "{{.Component}}"
		namespace: "{{.Namespace}}"
		{{- if .Configuration }}
		configuration: [
			{{- range $i, $e := .Configuration }}
			{{- if $i}}, {{- end}} {{ ConfigEntryInputToGQL $e }}
			{{- end }}]
		{{- end }}
	}`)
}

func (g *graphqlizer) ConfigEntryInputToGQL(in gqlschema.ConfigEntryInput) (string, error) {
	return g.genericToGraphQL(in, `{
		key: "{{.Key}}"
		value: "{{.Value}}"
		{{- if .Secret }}
		secret: {{.Secret}}
		{{- end }}
	}`)
}

func (g *graphqlizer) marshal(obj interface{}) string {
	var out string

	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Map:
		s, err := g.genericToGraphQL(obj, `{ {{- range $k, $v := . }}{{ $k }}:{{ marshal $v }},{{ end -}} }`)
		if err != nil {
			logrus.Warnf("failed to marshal labels: %s", err.Error())
			return ""
		}
		out = s
	case reflect.Slice, reflect.Array:
		s, err := g.genericToGraphQL(obj, `[{{ range $i, $e := . }}{{ if $i }},{{ end }}{{ marshal $e }}{{ end }}]`)
		if err != nil {
			logrus.Warnf("failed to marshal labels: %s", err.Error())
			return ""
		}
		out = s
	default:
		marshalled, err := json.Marshal(obj)
		if err != nil {
			logrus.Warnf("failed to marshal labels: %s", err.Error())
			return ""
		}
		out = string(marshalled)
	}

	return out
}

func (g *graphqlizer) genericToGraphQL(obj interface{}, tmpl string) (string, error) {
	fm := sprig.TxtFuncMap()
	fm["marshal"] = g.marshal
	fm["RuntimeInputToGraphQL"] = g.RuntimeInputToGraphQL
	fm["ComponentConfigurationInputToGQL"] = g.ComponentConfigurationInputToGQL
	fm["ConfigEntryInputToGQL"] = g.ConfigEntryInputToGQL
	fm["ClusterConfigToGraphQL"] = g.ClusterConfigToGraphQL
	fm["KymaConfigToGraphQL"] = g.KymaConfigToGraphQL
	fm["UpgradeClusterConfigToGraphQL"] = g.UpgradeClusterConfigToGraphQL
	fm["GardenerConfigInputToGraphQL"] = g.GardenerConfigInputToGraphQL
	fm["ProviderSpecificInputToGraphQL"] = g.ProviderSpecificInputToGraphQL
	fm["AzureProviderConfigInputToGraphQL"] = g.AzureProviderConfigInputToGraphQL
	fm["GcpProviderConfigInputToGraphQL"] = g.GcpProviderConfigInputToGraphQL
	fm["GardenerUpgradeInputToGraphQL"] = g.GardenerUpgradeInputToGraphQL

	t, err := template.New("tmpl").Funcs(fm).Parse(tmpl)
	if err != nil {
		return "", errors.Wrapf(err, "while parsing template")
	}

	var b bytes.Buffer

	if err := t.Execute(&b, obj); err != nil {
		return "", errors.Wrap(err, "while executing template")
	}
	return b.String(), nil
}

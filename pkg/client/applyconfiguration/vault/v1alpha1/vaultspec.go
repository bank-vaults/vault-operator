// Copyright © 2023 Bank-Vaults
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// VaultSpecApplyConfiguration represents an declarative configuration of the VaultSpec type for use
// with apply.
type VaultSpecApplyConfiguration struct {
	Size                       *int32                                            `json:"size,omitempty"`
	Image                      *string                                           `json:"image,omitempty"`
	BankVaultsImage            *string                                           `json:"bankVaultsImage,omitempty"`
	BankVaultsVolumeMounts     []v1.VolumeMount                                  `json:"bankVaultsVolumeMounts,omitempty"`
	StatsDDisabled             *bool                                             `json:"statsdDisabled,omitempty"`
	StatsDImage                *string                                           `json:"statsdImage,omitempty"`
	StatsdConfig               *string                                           `json:"statsdConfig,omitempty"`
	FluentDEnabled             *bool                                             `json:"fluentdEnabled,omitempty"`
	FluentDImage               *string                                           `json:"fluentdImage,omitempty"`
	FluentDConfLocation        *string                                           `json:"fluentdConfLocation,omitempty"`
	FluentDConfFile            *string                                           `json:"fluentdConfFile,omitempty"`
	FluentDConfig              *string                                           `json:"fluentdConfig,omitempty"`
	WatchedSecretsLabels       []map[string]string                               `json:"watchedSecretsLabels,omitempty"`
	WatchedSecretsAnnotations  []map[string]string                               `json:"watchedSecretsAnnotations,omitempty"`
	Annotations                map[string]string                                 `json:"annotations,omitempty"`
	VaultAnnotations           map[string]string                                 `json:"vaultAnnotations,omitempty"`
	VaultLabels                map[string]string                                 `json:"vaultLabels,omitempty"`
	VaultPodSpec               *EmbeddedPodSpecApplyConfiguration                `json:"vaultPodSpec,omitempty"`
	VaultContainerSpec         *v1.Container                                     `json:"vaultContainerSpec,omitempty"`
	VaultConfigurerAnnotations map[string]string                                 `json:"vaultConfigurerAnnotations,omitempty"`
	VaultConfigurerLabels      map[string]string                                 `json:"vaultConfigurerLabels,omitempty"`
	VaultConfigurerPodSpec     *EmbeddedPodSpecApplyConfiguration                `json:"vaultConfigurerPodSpec,omitempty"`
	Config                     *v1beta1.JSON                                     `json:"config,omitempty"`
	ExternalConfig             *v1beta1.JSON                                     `json:"externalConfig,omitempty"`
	UnsealConfig               *UnsealConfigApplyConfiguration                   `json:"unsealConfig,omitempty"`
	CredentialsConfig          *CredentialsConfigApplyConfiguration              `json:"credentialsConfig,omitempty"`
	EnvsConfig                 []v1.EnvVar                                       `json:"envsConfig,omitempty"`
	SecurityContext            *v1.PodSecurityContext                            `json:"securityContext,omitempty"`
	ServiceType                *string                                           `json:"serviceType,omitempty"`
	LoadBalancerIP             *string                                           `json:"loadBalancerIP,omitempty"`
	ServiceRegistrationEnabled *bool                                             `json:"serviceRegistrationEnabled,omitempty"`
	RaftLeaderAddress          *string                                           `json:"raftLeaderAddress,omitempty"`
	ServicePorts               map[string]int32                                  `json:"servicePorts,omitempty"`
	Affinity                   *v1.Affinity                                      `json:"affinity,omitempty"`
	PodAntiAffinity            *string                                           `json:"podAntiAffinity,omitempty"`
	NodeAffinity               *v1.NodeAffinity                                  `json:"nodeAffinity,omitempty"`
	NodeSelector               map[string]string                                 `json:"nodeSelector,omitempty"`
	Tolerations                []v1.Toleration                                   `json:"tolerations,omitempty"`
	ServiceAccount             *string                                           `json:"serviceAccount,omitempty"`
	Volumes                    []v1.Volume                                       `json:"volumes,omitempty"`
	VolumeMounts               []v1.VolumeMount                                  `json:"volumeMounts,omitempty"`
	VolumeClaimTemplates       []EmbeddedPersistentVolumeClaimApplyConfiguration `json:"volumeClaimTemplates,omitempty"`
	VaultEnvsConfig            []v1.EnvVar                                       `json:"vaultEnvsConfig,omitempty"`
	SidecarEnvsConfig          []v1.EnvVar                                       `json:"sidecarEnvsConfig,omitempty"`
	Resources                  *ResourcesApplyConfiguration                      `json:"resources,omitempty"`
	Ingress                    *IngressApplyConfiguration                        `json:"ingress,omitempty"`
	ServiceMonitorEnabled      *bool                                             `json:"serviceMonitorEnabled,omitempty"`
	ExistingTLSSecretName      *string                                           `json:"existingTlsSecretName,omitempty"`
	TLSExpiryThreshold         *string                                           `json:"tlsExpiryThreshold,omitempty"`
	TLSAdditionalHosts         []string                                          `json:"tlsAdditionalHosts,omitempty"`
	CANamespaces               []string                                          `json:"caNamespaces,omitempty"`
	IstioEnabled               *bool                                             `json:"istioEnabled,omitempty"`
	VeleroEnabled              *bool                                             `json:"veleroEnabled,omitempty"`
	VeleroFsfreezeImage        *string                                           `json:"veleroFsfreezeImage,omitempty"`
	VaultContainers            []v1.Container                                    `json:"vaultContainers,omitempty"`
	VaultInitContainers        []v1.Container                                    `json:"vaultInitContainers,omitempty"`
}

// VaultSpecApplyConfiguration constructs an declarative configuration of the VaultSpec type for use with
// apply.
func VaultSpec() *VaultSpecApplyConfiguration {
	return &VaultSpecApplyConfiguration{}
}

// WithSize sets the Size field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Size field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithSize(value int32) *VaultSpecApplyConfiguration {
	b.Size = &value
	return b
}

// WithImage sets the Image field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Image field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithImage(value string) *VaultSpecApplyConfiguration {
	b.Image = &value
	return b
}

// WithBankVaultsImage sets the BankVaultsImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the BankVaultsImage field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithBankVaultsImage(value string) *VaultSpecApplyConfiguration {
	b.BankVaultsImage = &value
	return b
}

// WithBankVaultsVolumeMounts adds the given value to the BankVaultsVolumeMounts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the BankVaultsVolumeMounts field.
func (b *VaultSpecApplyConfiguration) WithBankVaultsVolumeMounts(values ...v1.VolumeMount) *VaultSpecApplyConfiguration {
	for i := range values {
		b.BankVaultsVolumeMounts = append(b.BankVaultsVolumeMounts, values[i])
	}
	return b
}

// WithStatsDDisabled sets the StatsDDisabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the StatsDDisabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithStatsDDisabled(value bool) *VaultSpecApplyConfiguration {
	b.StatsDDisabled = &value
	return b
}

// WithStatsDImage sets the StatsDImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the StatsDImage field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithStatsDImage(value string) *VaultSpecApplyConfiguration {
	b.StatsDImage = &value
	return b
}

// WithStatsdConfig sets the StatsdConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the StatsdConfig field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithStatsdConfig(value string) *VaultSpecApplyConfiguration {
	b.StatsdConfig = &value
	return b
}

// WithFluentDEnabled sets the FluentDEnabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FluentDEnabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithFluentDEnabled(value bool) *VaultSpecApplyConfiguration {
	b.FluentDEnabled = &value
	return b
}

// WithFluentDImage sets the FluentDImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FluentDImage field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithFluentDImage(value string) *VaultSpecApplyConfiguration {
	b.FluentDImage = &value
	return b
}

// WithFluentDConfLocation sets the FluentDConfLocation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FluentDConfLocation field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithFluentDConfLocation(value string) *VaultSpecApplyConfiguration {
	b.FluentDConfLocation = &value
	return b
}

// WithFluentDConfFile sets the FluentDConfFile field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FluentDConfFile field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithFluentDConfFile(value string) *VaultSpecApplyConfiguration {
	b.FluentDConfFile = &value
	return b
}

// WithFluentDConfig sets the FluentDConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FluentDConfig field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithFluentDConfig(value string) *VaultSpecApplyConfiguration {
	b.FluentDConfig = &value
	return b
}

// WithWatchedSecretsLabels adds the given value to the WatchedSecretsLabels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the WatchedSecretsLabels field.
func (b *VaultSpecApplyConfiguration) WithWatchedSecretsLabels(values ...map[string]string) *VaultSpecApplyConfiguration {
	for i := range values {
		b.WatchedSecretsLabels = append(b.WatchedSecretsLabels, values[i])
	}
	return b
}

// WithWatchedSecretsAnnotations adds the given value to the WatchedSecretsAnnotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the WatchedSecretsAnnotations field.
func (b *VaultSpecApplyConfiguration) WithWatchedSecretsAnnotations(values ...map[string]string) *VaultSpecApplyConfiguration {
	for i := range values {
		b.WatchedSecretsAnnotations = append(b.WatchedSecretsAnnotations, values[i])
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *VaultSpecApplyConfiguration) WithAnnotations(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithVaultAnnotations puts the entries into the VaultAnnotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the VaultAnnotations field,
// overwriting an existing map entries in VaultAnnotations field with the same key.
func (b *VaultSpecApplyConfiguration) WithVaultAnnotations(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.VaultAnnotations == nil && len(entries) > 0 {
		b.VaultAnnotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.VaultAnnotations[k] = v
	}
	return b
}

// WithVaultLabels puts the entries into the VaultLabels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the VaultLabels field,
// overwriting an existing map entries in VaultLabels field with the same key.
func (b *VaultSpecApplyConfiguration) WithVaultLabels(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.VaultLabels == nil && len(entries) > 0 {
		b.VaultLabels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.VaultLabels[k] = v
	}
	return b
}

// WithVaultPodSpec sets the VaultPodSpec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VaultPodSpec field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithVaultPodSpec(value *EmbeddedPodSpecApplyConfiguration) *VaultSpecApplyConfiguration {
	b.VaultPodSpec = value
	return b
}

// WithVaultContainerSpec sets the VaultContainerSpec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VaultContainerSpec field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithVaultContainerSpec(value v1.Container) *VaultSpecApplyConfiguration {
	b.VaultContainerSpec = &value
	return b
}

// WithVaultConfigurerAnnotations puts the entries into the VaultConfigurerAnnotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the VaultConfigurerAnnotations field,
// overwriting an existing map entries in VaultConfigurerAnnotations field with the same key.
func (b *VaultSpecApplyConfiguration) WithVaultConfigurerAnnotations(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.VaultConfigurerAnnotations == nil && len(entries) > 0 {
		b.VaultConfigurerAnnotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.VaultConfigurerAnnotations[k] = v
	}
	return b
}

// WithVaultConfigurerLabels puts the entries into the VaultConfigurerLabels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the VaultConfigurerLabels field,
// overwriting an existing map entries in VaultConfigurerLabels field with the same key.
func (b *VaultSpecApplyConfiguration) WithVaultConfigurerLabels(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.VaultConfigurerLabels == nil && len(entries) > 0 {
		b.VaultConfigurerLabels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.VaultConfigurerLabels[k] = v
	}
	return b
}

// WithVaultConfigurerPodSpec sets the VaultConfigurerPodSpec field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VaultConfigurerPodSpec field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithVaultConfigurerPodSpec(value *EmbeddedPodSpecApplyConfiguration) *VaultSpecApplyConfiguration {
	b.VaultConfigurerPodSpec = value
	return b
}

// WithConfig sets the Config field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Config field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithConfig(value v1beta1.JSON) *VaultSpecApplyConfiguration {
	b.Config = &value
	return b
}

// WithExternalConfig sets the ExternalConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExternalConfig field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithExternalConfig(value v1beta1.JSON) *VaultSpecApplyConfiguration {
	b.ExternalConfig = &value
	return b
}

// WithUnsealConfig sets the UnsealConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the UnsealConfig field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithUnsealConfig(value *UnsealConfigApplyConfiguration) *VaultSpecApplyConfiguration {
	b.UnsealConfig = value
	return b
}

// WithCredentialsConfig sets the CredentialsConfig field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the CredentialsConfig field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithCredentialsConfig(value *CredentialsConfigApplyConfiguration) *VaultSpecApplyConfiguration {
	b.CredentialsConfig = value
	return b
}

// WithEnvsConfig adds the given value to the EnvsConfig field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the EnvsConfig field.
func (b *VaultSpecApplyConfiguration) WithEnvsConfig(values ...v1.EnvVar) *VaultSpecApplyConfiguration {
	for i := range values {
		b.EnvsConfig = append(b.EnvsConfig, values[i])
	}
	return b
}

// WithSecurityContext sets the SecurityContext field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SecurityContext field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithSecurityContext(value v1.PodSecurityContext) *VaultSpecApplyConfiguration {
	b.SecurityContext = &value
	return b
}

// WithServiceType sets the ServiceType field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceType field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithServiceType(value string) *VaultSpecApplyConfiguration {
	b.ServiceType = &value
	return b
}

// WithLoadBalancerIP sets the LoadBalancerIP field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LoadBalancerIP field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithLoadBalancerIP(value string) *VaultSpecApplyConfiguration {
	b.LoadBalancerIP = &value
	return b
}

// WithServiceRegistrationEnabled sets the ServiceRegistrationEnabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceRegistrationEnabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithServiceRegistrationEnabled(value bool) *VaultSpecApplyConfiguration {
	b.ServiceRegistrationEnabled = &value
	return b
}

// WithRaftLeaderAddress sets the RaftLeaderAddress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the RaftLeaderAddress field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithRaftLeaderAddress(value string) *VaultSpecApplyConfiguration {
	b.RaftLeaderAddress = &value
	return b
}

// WithServicePorts puts the entries into the ServicePorts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the ServicePorts field,
// overwriting an existing map entries in ServicePorts field with the same key.
func (b *VaultSpecApplyConfiguration) WithServicePorts(entries map[string]int32) *VaultSpecApplyConfiguration {
	if b.ServicePorts == nil && len(entries) > 0 {
		b.ServicePorts = make(map[string]int32, len(entries))
	}
	for k, v := range entries {
		b.ServicePorts[k] = v
	}
	return b
}

// WithAffinity sets the Affinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Affinity field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithAffinity(value v1.Affinity) *VaultSpecApplyConfiguration {
	b.Affinity = &value
	return b
}

// WithPodAntiAffinity sets the PodAntiAffinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PodAntiAffinity field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithPodAntiAffinity(value string) *VaultSpecApplyConfiguration {
	b.PodAntiAffinity = &value
	return b
}

// WithNodeAffinity sets the NodeAffinity field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeAffinity field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithNodeAffinity(value v1.NodeAffinity) *VaultSpecApplyConfiguration {
	b.NodeAffinity = &value
	return b
}

// WithNodeSelector puts the entries into the NodeSelector field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the NodeSelector field,
// overwriting an existing map entries in NodeSelector field with the same key.
func (b *VaultSpecApplyConfiguration) WithNodeSelector(entries map[string]string) *VaultSpecApplyConfiguration {
	if b.NodeSelector == nil && len(entries) > 0 {
		b.NodeSelector = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.NodeSelector[k] = v
	}
	return b
}

// WithTolerations adds the given value to the Tolerations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Tolerations field.
func (b *VaultSpecApplyConfiguration) WithTolerations(values ...v1.Toleration) *VaultSpecApplyConfiguration {
	for i := range values {
		b.Tolerations = append(b.Tolerations, values[i])
	}
	return b
}

// WithServiceAccount sets the ServiceAccount field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceAccount field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithServiceAccount(value string) *VaultSpecApplyConfiguration {
	b.ServiceAccount = &value
	return b
}

// WithVolumes adds the given value to the Volumes field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Volumes field.
func (b *VaultSpecApplyConfiguration) WithVolumes(values ...v1.Volume) *VaultSpecApplyConfiguration {
	for i := range values {
		b.Volumes = append(b.Volumes, values[i])
	}
	return b
}

// WithVolumeMounts adds the given value to the VolumeMounts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VolumeMounts field.
func (b *VaultSpecApplyConfiguration) WithVolumeMounts(values ...v1.VolumeMount) *VaultSpecApplyConfiguration {
	for i := range values {
		b.VolumeMounts = append(b.VolumeMounts, values[i])
	}
	return b
}

// WithVolumeClaimTemplates adds the given value to the VolumeClaimTemplates field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VolumeClaimTemplates field.
func (b *VaultSpecApplyConfiguration) WithVolumeClaimTemplates(values ...*EmbeddedPersistentVolumeClaimApplyConfiguration) *VaultSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithVolumeClaimTemplates")
		}
		b.VolumeClaimTemplates = append(b.VolumeClaimTemplates, *values[i])
	}
	return b
}

// WithVaultEnvsConfig adds the given value to the VaultEnvsConfig field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VaultEnvsConfig field.
func (b *VaultSpecApplyConfiguration) WithVaultEnvsConfig(values ...v1.EnvVar) *VaultSpecApplyConfiguration {
	for i := range values {
		b.VaultEnvsConfig = append(b.VaultEnvsConfig, values[i])
	}
	return b
}

// WithSidecarEnvsConfig adds the given value to the SidecarEnvsConfig field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the SidecarEnvsConfig field.
func (b *VaultSpecApplyConfiguration) WithSidecarEnvsConfig(values ...v1.EnvVar) *VaultSpecApplyConfiguration {
	for i := range values {
		b.SidecarEnvsConfig = append(b.SidecarEnvsConfig, values[i])
	}
	return b
}

// WithResources sets the Resources field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Resources field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithResources(value *ResourcesApplyConfiguration) *VaultSpecApplyConfiguration {
	b.Resources = value
	return b
}

// WithIngress sets the Ingress field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Ingress field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithIngress(value *IngressApplyConfiguration) *VaultSpecApplyConfiguration {
	b.Ingress = value
	return b
}

// WithServiceMonitorEnabled sets the ServiceMonitorEnabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ServiceMonitorEnabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithServiceMonitorEnabled(value bool) *VaultSpecApplyConfiguration {
	b.ServiceMonitorEnabled = &value
	return b
}

// WithExistingTLSSecretName sets the ExistingTLSSecretName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ExistingTLSSecretName field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithExistingTLSSecretName(value string) *VaultSpecApplyConfiguration {
	b.ExistingTLSSecretName = &value
	return b
}

// WithTLSExpiryThreshold sets the TLSExpiryThreshold field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the TLSExpiryThreshold field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithTLSExpiryThreshold(value string) *VaultSpecApplyConfiguration {
	b.TLSExpiryThreshold = &value
	return b
}

// WithTLSAdditionalHosts adds the given value to the TLSAdditionalHosts field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the TLSAdditionalHosts field.
func (b *VaultSpecApplyConfiguration) WithTLSAdditionalHosts(values ...string) *VaultSpecApplyConfiguration {
	for i := range values {
		b.TLSAdditionalHosts = append(b.TLSAdditionalHosts, values[i])
	}
	return b
}

// WithCANamespaces adds the given value to the CANamespaces field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the CANamespaces field.
func (b *VaultSpecApplyConfiguration) WithCANamespaces(values ...string) *VaultSpecApplyConfiguration {
	for i := range values {
		b.CANamespaces = append(b.CANamespaces, values[i])
	}
	return b
}

// WithIstioEnabled sets the IstioEnabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IstioEnabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithIstioEnabled(value bool) *VaultSpecApplyConfiguration {
	b.IstioEnabled = &value
	return b
}

// WithVeleroEnabled sets the VeleroEnabled field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VeleroEnabled field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithVeleroEnabled(value bool) *VaultSpecApplyConfiguration {
	b.VeleroEnabled = &value
	return b
}

// WithVeleroFsfreezeImage sets the VeleroFsfreezeImage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VeleroFsfreezeImage field is set to the value of the last call.
func (b *VaultSpecApplyConfiguration) WithVeleroFsfreezeImage(value string) *VaultSpecApplyConfiguration {
	b.VeleroFsfreezeImage = &value
	return b
}

// WithVaultContainers adds the given value to the VaultContainers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VaultContainers field.
func (b *VaultSpecApplyConfiguration) WithVaultContainers(values ...v1.Container) *VaultSpecApplyConfiguration {
	for i := range values {
		b.VaultContainers = append(b.VaultContainers, values[i])
	}
	return b
}

// WithVaultInitContainers adds the given value to the VaultInitContainers field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the VaultInitContainers field.
func (b *VaultSpecApplyConfiguration) WithVaultInitContainers(values ...v1.Container) *VaultSpecApplyConfiguration {
	for i := range values {
		b.VaultInitContainers = append(b.VaultInitContainers, values[i])
	}
	return b
}
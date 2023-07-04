// Copyright © 2019 Banzai Cloud
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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Vault is the Schema for the vaults API
// Deprecated. Use v1alpha2.Vault
type Vault struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VaultSpec   `json:"spec,omitempty"`
	Status VaultStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VaultList contains a list of Vault
// Deprecated. Use v1alpha2.VaultList
type VaultList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Vault `json:"items"`
}

// VaultSpec defines the desired state of Vault
// Deprecated. Use v1alpha2.VaultSpec
type VaultSpec struct {

	// Size defines the number of Vault instances in the cluster (>= 1 means HA)
	// default: 1
	Size int32 `json:"size,omitempty"`

	// Image specifies the Vault image to use for the Vault instances
	// default: library/vault:latest
	Image string `json:"image,omitempty"`

	// BankVaultsImage specifies the Bank Vaults image to use for Vault unsealing and configuration
	// default: ghcr.io/banzaicloud/bank-vaults:latest
	BankVaultsImage string `json:"bankVaultsImage,omitempty"`

	// BankVaultsVolumeMounts define some extra Kubernetes Volume mounts for the Bank Vaults Sidecar container.
	// default:
	BankVaultsVolumeMounts []v1.VolumeMount `json:"bankVaultsVolumeMounts,omitempty"`

	// StatsDDisabled specifies if StatsD based metrics should be disabled
	// default: false
	StatsDDisabled bool `json:"statsdDisabled,omitempty"`

	// StatsDImage specifices the StatsD image to use for Vault metrics exportation
	// default: prom/statsd-exporter:latest
	StatsDImage string `json:"statsdImage,omitempty"`

	// StatsdConfig specifices the StatsD mapping configuration
	// default:
	StatsdConfig string `json:"statsdConfig,omitempty"`

	// FluentDEnabled specifies if FluentD based log exportation should be enabled
	// default: false
	FluentDEnabled bool `json:"fluentdEnabled,omitempty"`

	// FluentDImage specifices the FluentD image to use for Vault log exportation
	// default: fluent/fluentd:edge
	FluentDImage string `json:"fluentdImage,omitempty"`

	// FluentDConfLocation is the location of the fluent.conf file
	// default: "/fluentd/etc"
	FluentDConfLocation string `json:"fluentdConfLocation,omitempty"`

	// FluentDConfFile specifices the FluentD configuration file name to use for Vault log exportation
	// default:
	FluentDConfFile string `json:"fluentdConfFile,omitempty"`

	// FluentDConfig specifices the FluentD configuration to use for Vault log exportation
	// default:
	FluentDConfig string `json:"fluentdConfig,omitempty"`

	// WatchedSecretsLabels specifices a set of Kubernetes label selectors which select Secrets to watch.
	// If these Secrets change the Vault cluster gets restarted. For example a Secret that Cert-Manager is
	// managing a public Certificate for Vault using let's Encrypt.
	// default:
	WatchedSecretsLabels []map[string]string `json:"watchedSecretsLabels,omitempty"`

	// WatchedSecretsAnnotations specifices a set of Kubernetes annotations selectors which select Secrets to watch.
	// If these Secrets change the Vault cluster gets restarted. For example a Secret that Cert-Manager is
	// managing a public Certificate for Vault using let's Encrypt.
	// default:
	WatchedSecretsAnnotations []map[string]string `json:"watchedSecretsAnnotations,omitempty"`

	// Annotations define a set of common Kubernetes annotations that will be added to all operator managed resources.
	// default:
	Annotations map[string]string `json:"annotations,omitempty"`

	// VaultAnnotations define a set of Kubernetes annotations that will be added to all Vault Pods.
	// default:
	VaultAnnotations map[string]string `json:"vaultAnnotations,omitempty"`

	// VaultLabels define a set of Kubernetes labels that will be added to all Vault Pods.
	// default:
	VaultLabels map[string]string `json:"vaultLabels,omitempty"`

	// VaultPodSpec is a Kubernetes Pod specification snippet (`spec:` block) that will be merged into the operator generated
	// Vault Pod specification.
	// default:
	VaultPodSpec *EmbeddedPodSpec `json:"vaultPodSpec,omitempty"`

	// VaultContainerSpec is a Kubernetes Container specification snippet that will be merged into the operator generated
	// Vault Container specification.
	// default:
	VaultContainerSpec v1.Container `json:"vaultContainerSpec,omitempty"`

	// VaultConfigurerAnnotations define a set of Kubernetes annotations that will be added to the Vault Configurer Pod.
	// default:
	VaultConfigurerAnnotations map[string]string `json:"vaultConfigurerAnnotations,omitempty"`

	// VaultConfigurerLabels define a set of Kubernetes labels that will be added to all Vault Configurer Pod.
	// default:
	VaultConfigurerLabels map[string]string `json:"vaultConfigurerLabels,omitempty"`

	// VaultConfigurerPodSpec is a Kubernetes Pod specification snippet (`spec:` block) that will be merged into
	// the operator generated Vault Configurer Pod specification.
	// default:
	VaultConfigurerPodSpec *EmbeddedPodSpec `json:"vaultConfigurerPodSpec,omitempty"`

	// Config is the Vault Server configuration. See https://www.vaultproject.io/docs/configuration/ for more details.
	// default:
	Config extv1beta1.JSON `json:"config"`

	// ExternalConfig is higher level configuration block which instructs the Bank Vaults Configurer to configure Vault
	// through its API, thus allows setting up:
	// - Secret Engines
	// - Auth Methods
	// - Audit Devices
	// - Plugin Backends
	// - Policies
	// - Startup Secrets (Bank Vaults feature)
	// A documented example: https://github.com/bank-vaults/vault-operator/blob/main/vault-config.yml
	// default:
	ExternalConfig extv1beta1.JSON `json:"externalConfig,omitempty"`

	// UnsealConfig defines where the Vault cluster's unseal keys and root token should be stored after initialization.
	// See the type's documentation for more details. Only one method may be specified.
	// default: Kubernetes Secret based unsealing
	UnsealConfig UnsealConfig `json:"unsealConfig,omitempty"`

	// CredentialsConfig defines a external Secret for Vault and how it should be mounted to the Vault Pod
	// for example accessing Cloud resources.
	// default:
	CredentialsConfig CredentialsConfig `json:"credentialsConfig,omitempty"`

	// EnvsConfig is a list of Kubernetes environment variable definitions that will be passed to all Bank-Vaults pods.
	// default:
	EnvsConfig []v1.EnvVar `json:"envsConfig,omitempty"`

	// SecurityContext is a Kubernetes PodSecurityContext that will be applied to all Pods created by the operator.
	// default:
	SecurityContext v1.PodSecurityContext `json:"securityContext,omitempty"`

	// ServiceType is a Kubernetes Service type of the Vault Service.
	// default: ClusterIP
	ServiceType string `json:"serviceType,omitempty"`

	// LoadBalancerIP is an optional setting for allocating a specific address for the entry service object
	// of type LoadBalancer
	// default: ""
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`

	// serviceRegistrationEnabled enables the injection of the service_registration Vault stanza.
	// This requires elaborated RBAC privileges for updating Pod labels for the Vault Pod.
	// default: false
	ServiceRegistrationEnabled bool `json:"serviceRegistrationEnabled,omitempty"`

	// RaftLeaderAddress defines the leader address of the raft cluster in multi-cluster deployments.
	// (In single cluster (namespace) deployments it is automatically detected).
	// "self" is a special value which means that this instance should be the bootstrap leader instance.
	// default: ""
	RaftLeaderAddress string `json:"raftLeaderAddress,omitempty"`

	// ServicePorts is an extra map of ports that should be exposed by the Vault Service.
	// default:
	ServicePorts map[string]int32 `json:"servicePorts,omitempty"`

	// Affinity is a group of affinity scheduling rules applied to all Vault Pods.
	// default:
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// PodAntiAffinity is the TopologyKey in the Vault Pod's PodAntiAffinity.
	// No PodAntiAffinity is used if empty.
	// Deprecated. Use Affinity.
	// default:
	PodAntiAffinity string `json:"podAntiAffinity,omitempty"`

	// NodeAffinity is Kubernetees NodeAffinity definition that should be applied to all Vault Pods.
	// Deprecated. Use Affinity.
	// default:
	NodeAffinity v1.NodeAffinity `json:"nodeAffinity,omitempty"`

	// NodeSelector is Kubernetees NodeSelector definition that should be applied to all Vault Pods.
	// default:
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations is Kubernetes Tolerations definition that should be applied to all Vault Pods.
	// default:
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// ServiceAccount is Kubernetes ServiceAccount in which the Vault Pods should be running in.
	// default: default
	ServiceAccount string `json:"serviceAccount,omitempty"`

	// Volumes define some extra Kubernetes Volumes for the Vault Pods.
	// default:
	Volumes []v1.Volume `json:"volumes,omitempty"`

	// VolumeMounts define some extra Kubernetes Volume mounts for the Vault Pods.
	// default:
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`

	// VolumeClaimTemplates define some extra Kubernetes PersistentVolumeClaim templates for the Vault Statefulset.
	// default:
	VolumeClaimTemplates []EmbeddedPersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`

	// VaultEnvsConfig is a list of Kubernetes environment variable definitions that will be passed to the Vault container.
	// default:
	VaultEnvsConfig []v1.EnvVar `json:"vaultEnvsConfig,omitempty"`

	// SidecarEnvsConfig is a list of Kubernetes environment variable definitions that will be passed to Vault sidecar containers.
	// default:
	SidecarEnvsConfig []v1.EnvVar `json:"sidecarEnvsConfig,omitempty"`

	// Resources defines the resource limits for all the resources created by the operator.
	// See the type for more details.
	// default:
	Resources *Resources `json:"resources,omitempty"`

	// Ingress, if it is specified the operator will create an Ingress resource for the Vault Service and
	// will annotate it with the correct Ingress annotations specific to the TLS settings in the configuration.
	// See the type for more details.
	// default:
	Ingress *Ingress `json:"ingress,omitempty"`

	// ServiceMonitorEnabled enables the creation of Prometheus Operator specific ServiceMonitor for Vault.
	// default: false
	ServiceMonitorEnabled bool `json:"serviceMonitorEnabled,omitempty"`

	// ExistingTLSSecretName is name of the secret that contains a TLS server certificate and key and the corresponding CA certificate.
	// Required secret format kubernetes.io/tls type secret keys + ca.crt key
	// If it is set, generating certificate will be disabled
	// default: ""
	ExistingTLSSecretName string `json:"existingTlsSecretName,omitempty"`

	// TLSExpiryThreshold is the Vault TLS certificate expiration threshold in Go's Duration format.
	// default: 168h
	TLSExpiryThreshold string `json:"tlsExpiryThreshold,omitempty"`

	// TLSAdditionalHosts is a list of additional hostnames or IP addresses to add to the SAN on the automatically generated TLS certificate.
	// default:
	TLSAdditionalHosts []string `json:"tlsAdditionalHosts,omitempty"`

	// CANamespaces define a list of namespaces where the generated CA certificate for Vault should be distributed,
	// use ["*"] for all namespaces.
	// default:
	CANamespaces []string `json:"caNamespaces,omitempty"`

	// IstioEnabled describes if the cluster has a Istio running and enabled.
	// default: false
	IstioEnabled bool `json:"istioEnabled,omitempty"`

	// VeleroEnabled describes if the cluster has a Velero running and enabled.
	// default: false
	VeleroEnabled bool `json:"veleroEnabled,omitempty"`

	// VeleroFsfreezeImage specifices the Velero Fsrfeeze image to use in Velero backup hooks
	// default: velero/fsfreeze-pause:latest
	VeleroFsfreezeImage string `json:"veleroFsfreezeImage,omitempty"`

	// VaultContainers add extra containers
	VaultContainers []v1.Container `json:"vaultContainers,omitempty"`

	// VaultInitContainers add extra initContainers
	VaultInitContainers []v1.Container `json:"vaultInitContainers,omitempty"`
}

// VaultStatus defines the observed state of Vault
// Deprecated. Use v1alpha2.VaultStatus
type VaultStatus struct {
	// Important: Run "make generate-code" to regenerate code after modifying this file
	Nodes      []string                `json:"nodes"`
	Leader     string                  `json:"leader"`
	Conditions []v1.ComponentCondition `json:"conditions,omitempty"`
}

// UnsealConfig represents the UnsealConfig field of a VaultSpec Kubernetes object
type UnsealConfig struct {
	Options    UnsealOptions          `json:"options,omitempty"`
	Kubernetes KubernetesUnsealConfig `json:"kubernetes,omitempty"`
	Google     *GoogleUnsealConfig    `json:"google,omitempty"`
	Alibaba    *AlibabaUnsealConfig   `json:"alibaba,omitempty"`
	Azure      *AzureUnsealConfig     `json:"azure,omitempty"`
	AWS        *AWSUnsealConfig       `json:"aws,omitempty"`
	Vault      *VaultUnsealConfig     `json:"vault,omitempty"`
	HSM        *HSMUnsealConfig       `json:"hsm,omitempty"`
}

// UnsealOptions represents the common options to all unsealing backends
type UnsealOptions struct {
	PreFlightChecks *bool `json:"preFlightChecks,omitempty"`
	StoreRootToken  *bool `json:"storeRootToken,omitempty"`
	SecretThreshold *uint `json:"secretThreshold,omitempty"`
	SecretShares    *uint `json:"secretShares,omitempty"`
}

// KubernetesUnsealConfig holds the parameters for Kubernetes based unsealing
type KubernetesUnsealConfig struct {
	SecretNamespace string `json:"secretNamespace,omitempty"`
	SecretName      string `json:"secretName,omitempty"`
}

// GoogleUnsealConfig holds the parameters for Google KMS based unsealing
type GoogleUnsealConfig struct {
	KMSKeyRing    string `json:"kmsKeyRing"`
	KMSCryptoKey  string `json:"kmsCryptoKey"`
	KMSLocation   string `json:"kmsLocation"`
	KMSProject    string `json:"kmsProject"`
	StorageBucket string `json:"storageBucket"`
}

// AlibabaUnsealConfig holds the parameters for Alibaba Cloud KMS based unsealing
//
//	--alibaba-kms-region eu-central-1 --alibaba-kms-key-id 9d8063eb-f9dc-421b-be80-15d195c9f148 --alibaba-oss-endpoint oss-eu-central-1.aliyuncs.com --alibaba-oss-bucket bank-vaults
type AlibabaUnsealConfig struct {
	KMSRegion   string `json:"kmsRegion"`
	KMSKeyID    string `json:"kmsKeyId"`
	OSSEndpoint string `json:"ossEndpoint"`
	OSSBucket   string `json:"ossBucket"`
	OSSPrefix   string `json:"ossPrefix"`
}

// AzureUnsealConfig holds the parameters for Azure Key Vault based unsealing
type AzureUnsealConfig struct {
	KeyVaultName string `json:"keyVaultName"`
}

// AWSUnsealConfig holds the parameters for AWS KMS based unsealing
type AWSUnsealConfig struct {
	KMSKeyID  string `json:"kmsKeyId"`
	KMSRegion string `json:"kmsRegion,omitempty"`
	S3Bucket  string `json:"s3Bucket"`
	S3Prefix  string `json:"s3Prefix"`
	S3Region  string `json:"s3Region,omitempty"`
	S3SSE     string `json:"s3SSE,omitempty"`
}

// VaultUnsealConfig holds the parameters for remote Vault based unsealing
type VaultUnsealConfig struct {
	Address        string `json:"address"`
	UnsealKeysPath string `json:"unsealKeysPath"`
	Role           string `json:"role,omitempty"`
	AuthPath       string `json:"authPath,omitempty"`
	TokenPath      string `json:"tokenPath,omitempty"`
	Token          string `json:"token,omitempty"`
}

// HSMUnsealConfig holds the parameters for remote HSM based unsealing
type HSMUnsealConfig struct {
	Daemon     bool   `json:"daemon,omitempty"`
	ModulePath string `json:"modulePath"`
	SlotID     uint   `json:"slotId,omitempty"`
	TokenLabel string `json:"tokenLabel,omitempty"`
	Pin        string `json:"pin"`
	KeyLabel   string `json:"keyLabel"`
}

// CredentialsConfig configuration for a credentials file provided as a secret
type CredentialsConfig struct {
	Env        string `json:"env"`
	Path       string `json:"path"`
	SecretName string `json:"secretName"`
}

// Resources holds different container's ResourceRequirements
type Resources struct {
	Vault              *v1.ResourceRequirements `json:"vault,omitempty"`
	BankVaults         *v1.ResourceRequirements `json:"bankVaults,omitempty"`
	HSMDaemon          *v1.ResourceRequirements `json:"hsmDaemon,omitempty"`
	PrometheusExporter *v1.ResourceRequirements `json:"prometheusExporter,omitempty"`
	FluentD            *v1.ResourceRequirements `json:"fluentd,omitempty"`
}

// Ingress specification for the Vault cluster
type Ingress struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Spec        netv1.IngressSpec `json:"spec,omitempty"`
}
